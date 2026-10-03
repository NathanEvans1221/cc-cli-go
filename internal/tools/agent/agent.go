package agent

import (
	"context"

	"github.com/user-name/cc-cli-go/internal/tools"
)

type Runner interface {
	Run(ctx context.Context, prompt string) (string, error)
}

type AgentTool struct {
	Runner Runner
}

func New(runner Runner) *AgentTool { return &AgentTool{Runner: runner} }

func (t *AgentTool) Name() string        { return "Agent" }
func (t *AgentTool) Description() string { return "Delegate one prompt to a runner." }
func (t *AgentTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{"type": "object"}
}
func (t *AgentTool) IsEnabled() bool                               { return t.Runner != nil }
func (t *AgentTool) IsReadOnly(map[string]interface{}) bool        { return true }
func (t *AgentTool) IsConcurrencySafe(map[string]interface{}) bool { return false }
func (t *AgentTool) UserFacingName(map[string]interface{}) string  { return "Agent" }

func (t *AgentTool) Execute(ctx context.Context, input map[string]interface{}, _ *tools.ToolContext) (*tools.ToolResult, error) {
	prompt, _ := input["prompt"].(string)
	if t.Runner == nil || prompt == "" {
		return &tools.ToolResult{Content: "prompt and runner are required", IsError: true}, nil
	}
	out, err := t.Runner.Run(ctx, prompt)
	if err != nil {
		return &tools.ToolResult{Content: err.Error(), IsError: true}, nil
	}
	return &tools.ToolResult{Content: out}, nil
}
