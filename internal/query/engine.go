package query

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/user-name/cc-cli-go/internal/api"
	"github.com/user-name/cc-cli-go/internal/permission"
	"github.com/user-name/cc-cli-go/internal/tools"
	"github.com/user-name/cc-cli-go/internal/types"
)

const defaultMaxTurns = 8

type streamer interface {
	Stream(ctx context.Context, req *api.Request) (<-chan api.StreamEvent, error)
}

type Engine struct {
	client  streamer
	toolReg *tools.Registry
}

func NewEngine(client *api.Client, toolReg *tools.Registry) *Engine {
	return NewEngineWithStreamer(client, toolReg)
}

func NewEngineWithStreamer(client streamer, toolReg *tools.Registry) *Engine {
	return &Engine{
		client:  client,
		toolReg: toolReg,
	}
}

func (e *Engine) Query(ctx context.Context, params QueryParams) (<-chan StreamEvent, <-chan QueryResult) {
	events := make(chan StreamEvent, 100)
	results := make(chan QueryResult, 1)

	go func() {
		defer close(events)
		defer close(results)

		e.runQuery(ctx, params, events, results)
	}()

	return events, results
}

func (e *Engine) runQuery(ctx context.Context, params QueryParams, events chan<- StreamEvent, results chan<- QueryResult) {
	messages := append([]*types.Message{}, params.Messages...)
	maxTurns := params.MaxTurns
	if maxTurns <= 0 {
		maxTurns = defaultMaxTurns
	}

	for turn := 0; turn < maxTurns; turn++ {
		req := api.NewRequest(params.Model, params.MaxTokens)
		req.SetSystem(params.SystemPrompt)

		for _, msg := range messages {
			req.AddMessage(msg)
		}

		for _, tool := range params.Tools {
			req.AddTool(api.ToolParam{
				Name:        tool.Name(),
				Description: tool.Description(),
				InputSchema: tool.InputSchema(),
			})
		}

		stream, err := e.client.Stream(ctx, req)
		if err != nil {
			results <- QueryResult{Reason: "error", Error: err}
			return
		}

		var currentMessage *types.Message
		var currentContent *types.ContentBlock
		var turnBlocks []types.ContentBlock
		var toolUses []types.ContentBlock
		stopped := false

		for event := range stream {
			switch event.Type {
			case "message_start":
				msg, err := api.ParseMessageStart(event.Message)
				if err == nil {
					currentMessage = msg
					events <- StreamEvent{Type: "message_start", Message: msg}
				}

			case "content_block_start":
				block, err := api.ParseContentBlock(event.ContentBlock)
				if err == nil {
					currentContent = block
					events <- StreamEvent{Type: "content_block_start", Content: block}
				}

			case "content_block_delta":
				deltaType, deltaText, err := api.ParseDelta(event.Delta)
				if err == nil && deltaText != "" {
					if currentContent != nil && deltaType == "text_delta" {
						currentContent.Text += deltaText
					}
					events <- StreamEvent{Type: "content_block_delta", Delta: deltaText}
				}

			case "content_block_stop":
				if currentContent != nil {
					turnBlocks = append(turnBlocks, *currentContent)
					if currentContent.Type == "tool_use" {
						toolUses = append(toolUses, *currentContent)
					}
					currentContent = nil
				}
				events <- StreamEvent{Type: "content_block_stop"}

			case "message_delta":
				var delta api.MessageDeltaEvent
				if err := json.Unmarshal(event.Delta, &delta); err == nil {
					if currentMessage != nil {
						currentMessage.StopReason = delta.Delta.StopReason
					}
				}

			case "message_stop":
				events <- StreamEvent{Type: "message_stop"}
				stopped = true
			}
		}

		if !stopped {
			results <- QueryResult{Reason: "error", Error: fmt.Errorf("stream ended before message_stop")}
			return
		}

		if len(toolUses) == 0 {
			results <- QueryResult{Reason: "completed"}
			return
		}

		toolResults := e.executeTools(ctx, toolUses, params, events)
		messages = append(messages, &types.Message{
			Role:    "assistant",
			Content: turnBlocks,
		}, &types.Message{
			Role:    "user",
			Content: toolResultBlocks(toolUses, toolResults),
		})
	}

	results <- QueryResult{Reason: "max_turns"}
}

func toolResultBlocks(toolUses []types.ContentBlock, toolResults []*tools.ToolResult) []types.ContentBlock {
	blocks := make([]types.ContentBlock, len(toolUses))
	for i, toolUse := range toolUses {
		content := "tool result missing"
		isError := true
		if i < len(toolResults) && toolResults[i] != nil {
			content = fmt.Sprint(toolResults[i].Content)
			isError = toolResults[i].IsError
		}
		blocks[i] = types.NewToolResultBlock(toolUse.ID, content, isError)
	}
	return blocks
}

func (e *Engine) executeTools(ctx context.Context, toolUses []types.ContentBlock, params QueryParams, events chan<- StreamEvent) []*tools.ToolResult {
	var wg sync.WaitGroup
	var mu sync.Mutex
	results := make([]*tools.ToolResult, len(toolUses))

	publish := func(idx int, result *tools.ToolResult) {
		mu.Lock()
		results[idx] = result
		mu.Unlock()
	}

	for i, tu := range toolUses {
		wg.Add(1)
		go func(idx int, toolUse types.ContentBlock) {
			defer wg.Done()

			tool := e.toolReg.Get(toolUse.Name)
			if tool == nil {
				publish(idx, &tools.ToolResult{
					Content: "tool not found",
					IsError: true,
				})
				return
			}

			input, _ := toolUse.Input.(map[string]interface{})
			execCtx := ctx
			if params.ToolTimeout > 0 {
				var cancel context.CancelFunc
				execCtx, cancel = context.WithTimeout(ctx, params.ToolTimeout)
				defer cancel()
			}

			if params.PermissionChecker != nil {
				decision := params.PermissionChecker.Check(toolUse.Name, input)
				if decision == nil || decision.Behavior == permission.BehaviorDeny {
					reason := "denied"
					if decision != nil {
						reason = decision.Reason
					}
					publish(idx, &tools.ToolResult{
						Content: fmt.Sprintf("Permission denied: %s", reason),
						IsError: true,
					})
					return
				}

				if decision.Behavior == permission.BehaviorAsk {
					events <- StreamEvent{
						Type: "permission_request",
						PermissionRequest: &PermissionRequestEvent{
							ToolName: toolUse.Name,
							Input:    input,
							Decision: decision,
							Index:    idx,
						},
					}
					publish(idx, &tools.ToolResult{
						Content: fmt.Sprintf("Permission required: %s", decision.Reason),
						IsError: true,
					})
					return
				}
			}

			result, err := tool.Execute(execCtx, input, &tools.ToolContext{AbortSignal: execCtx})
			if err != nil {
				result = &tools.ToolResult{
					Content: err.Error(),
					IsError: true,
				}
			}
			publish(idx, result)
		}(i, tu)
	}

	wg.Wait()
	return results
}
