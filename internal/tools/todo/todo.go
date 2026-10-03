package todo

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/user-name/cc-cli-go/internal/tools"
)

type TodoTool struct{}

func New() *TodoTool { return &TodoTool{} }

func (t *TodoTool) Name() string        { return "TodoWrite" }
func (t *TodoTool) Description() string { return "Write a todo list to a file or return it." }
func (t *TodoTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{"type": "object"}
}
func (t *TodoTool) IsEnabled() bool                               { return true }
func (t *TodoTool) IsReadOnly(map[string]interface{}) bool        { return false }
func (t *TodoTool) IsConcurrencySafe(map[string]interface{}) bool { return false }
func (t *TodoTool) UserFacingName(map[string]interface{}) string  { return "TodoWrite" }

func (t *TodoTool) Execute(_ context.Context, input map[string]interface{}, _ *tools.ToolContext) (*tools.ToolResult, error) {
	var items []string
	switch raw := input["items"].(type) {
	case []string:
		items = raw
	case []interface{}:
		for _, item := range raw {
			items = append(items, fmt.Sprint(item))
		}
	case string:
		items = []string{raw}
	}
	if len(items) == 0 {
		return &tools.ToolResult{Content: "items are required", IsError: true}, nil
	}
	body := strings.Join(items, "\n")
	if path, _ := input["path"].(string); path != "" {
		if err := os.WriteFile(path, []byte(body+"\n"), 0644); err != nil {
			return &tools.ToolResult{Content: err.Error(), IsError: true}, nil
		}
	}
	return &tools.ToolResult{Content: body}, nil
}
