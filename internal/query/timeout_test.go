package query

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/user-name/cc-cli-go/internal/api"
	"github.com/user-name/cc-cli-go/internal/tools"
	"github.com/user-name/cc-cli-go/internal/types"
)

type blockingTool struct{ calls int }

func (t *blockingTool) Name() string                                  { return "Block" }
func (t *blockingTool) Description() string                           { return "blocks" }
func (t *blockingTool) InputSchema() map[string]interface{}           { return map[string]interface{}{} }
func (t *blockingTool) IsEnabled() bool                               { return true }
func (t *blockingTool) IsReadOnly(map[string]interface{}) bool        { return true }
func (t *blockingTool) IsConcurrencySafe(map[string]interface{}) bool { return true }
func (t *blockingTool) UserFacingName(map[string]interface{}) string  { return "Block" }
func (t *blockingTool) Execute(ctx context.Context, _ map[string]interface{}, _ *tools.ToolContext) (*tools.ToolResult, error) {
	t.calls++
	<-ctx.Done()
	return &tools.ToolResult{Content: ctx.Err().Error(), IsError: true}, nil
}

func TestQuery_ToolTimeoutDoesNotReportSuccess(t *testing.T) {
	tool := &blockingTool{}
	reg := tools.NewRegistry()
	reg.Register(tool)
	script := &scriptStreamer{turns: [][]api.StreamEvent{
		toolUseTurn("toolu_t", "Block", map[string]interface{}{"q": "wait"}),
		endTurn(),
	}}
	eng := NewEngineWithStreamer(script, reg)
	res := runQuery(t, eng, QueryParams{
		Messages:    []*types.Message{types.NewUserMessage("wait")},
		Tools:       []tools.Tool{tool},
		Model:       "test",
		MaxTokens:   8,
		MaxTurns:    3,
		ToolTimeout: 20 * time.Millisecond,
	})
	if res.Error != nil {
		t.Fatal(res.Error)
	}
	if tool.calls != 1 {
		t.Fatalf("calls %d", tool.calls)
	}
	body := script.bodies()[1]
	if !strings.Contains(body, "context deadline exceeded") {
		t.Fatalf("follow-up %s", body)
	}
	if !strings.Contains(body, `"is_error":true`) {
		t.Fatalf("timeout was not an error result: %s", body)
	}
}
