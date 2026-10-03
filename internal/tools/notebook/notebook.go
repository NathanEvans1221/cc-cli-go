package notebook

import (
	"context"
	"encoding/json"
	"os"

	"github.com/user-name/cc-cli-go/internal/tools"
)

type NotebookTool struct{}

func New() *NotebookTool { return &NotebookTool{} }

func (t *NotebookTool) Name() string        { return "NotebookEdit" }
func (t *NotebookTool) Description() string { return "Replace the source of one Jupyter cell." }
func (t *NotebookTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{"type": "object"}
}
func (t *NotebookTool) IsEnabled() bool                               { return true }
func (t *NotebookTool) IsReadOnly(map[string]interface{}) bool        { return false }
func (t *NotebookTool) IsConcurrencySafe(map[string]interface{}) bool { return false }
func (t *NotebookTool) UserFacingName(map[string]interface{}) string  { return "NotebookEdit" }

func (t *NotebookTool) Execute(_ context.Context, input map[string]interface{}, _ *tools.ToolContext) (*tools.ToolResult, error) {
	path, _ := input["path"].(string)
	source, _ := input["source"].(string)
	if path == "" || source == "" {
		return &tools.ToolResult{Content: "path and source are required", IsError: true}, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return &tools.ToolResult{Content: err.Error(), IsError: true}, nil
	}
	var doc map[string]interface{}
	if err := json.Unmarshal(data, &doc); err != nil {
		return &tools.ToolResult{Content: err.Error(), IsError: true}, nil
	}
	cells, _ := doc["cells"].([]interface{})
	if len(cells) == 0 {
		return &tools.ToolResult{Content: "notebook has no cells", IsError: true}, nil
	}
	cell, _ := cells[0].(map[string]interface{})
	if cell == nil {
		return &tools.ToolResult{Content: "first cell is not an object", IsError: true}, nil
	}
	cell["source"] = []string{source}
	out, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return &tools.ToolResult{Content: err.Error(), IsError: true}, nil
	}
	if err := os.WriteFile(path, out, 0644); err != nil {
		return &tools.ToolResult{Content: err.Error(), IsError: true}, nil
	}
	return &tools.ToolResult{Content: source}, nil
}
