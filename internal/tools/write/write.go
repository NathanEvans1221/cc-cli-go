package write

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/liao-eli/cc-cli-go/internal/tools"
)

type WriteTool struct{}

func New() *WriteTool {
	return &WriteTool{}
}

func (t *WriteTool) Name() string {
	return "Write"
}

func (t *WriteTool) Description() string {
	return "Write content to a new file. The file must not already exist."
}

func (t *WriteTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"file_path": map[string]interface{}{
				"type":        "string",
				"description": "The absolute path to the file to write (must not already exist)",
			},
			"content": map[string]interface{}{
				"type":        "string",
				"description": "The content to write to the file",
			},
		},
		"required": []string{"file_path", "content"},
	}
}

func (t *WriteTool) Execute(ctx context.Context, input map[string]interface{}, tc *tools.ToolContext) (*tools.ToolResult, error) {
	filePath, _ := input["file_path"].(string)
	content, _ := input["content"].(string)

	if filePath == "" {
		return &tools.ToolResult{
			Content: "Error: file_path is required",
			IsError: true,
		}, nil
	}

	if _, err := os.Stat(filePath); err == nil {
		return &tools.ToolResult{
			Content: fmt.Sprintf("Error: file already exists: %s. Use Edit tool to modify existing files.", filePath),
			IsError: true,
		}, nil
	}

	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return &tools.ToolResult{
			Content: fmt.Sprintf("Error creating directory: %v", err),
			IsError: true,
		}, nil
	}

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return &tools.ToolResult{
			Content: fmt.Sprintf("Error writing file: %v", err),
			IsError: true,
		}, nil
	}

	return &tools.ToolResult{
		Content: fmt.Sprintf("Successfully wrote to %s", filePath),
	}, nil
}

func (t *WriteTool) IsEnabled() bool {
	return true
}

func (t *WriteTool) IsReadOnly(input map[string]interface{}) bool {
	return false
}

func (t *WriteTool) IsConcurrencySafe(input map[string]interface{}) bool {
	return false
}

func (t *WriteTool) UserFacingName(input map[string]interface{}) string {
	if path, ok := input["file_path"].(string); ok {
		return path
	}
	return "Write"
}
