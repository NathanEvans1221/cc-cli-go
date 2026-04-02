package glob

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/liao-eli/cc-cli-go/internal/tools"
)

type GlobTool struct{}

func New() *GlobTool {
	return &GlobTool{}
}

func (t *GlobTool) Name() string {
	return "Glob"
}

func (t *GlobTool) Description() string {
	return "Fast file pattern matching tool that works with any codebase size. Supports glob patterns like '**/*.js' or 'src/**/*.ts'."
}

func (t *GlobTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"pattern": map[string]interface{}{
				"type":        "string",
				"description": "The glob pattern to match files against (e.g. '**/*.js')",
			},
			"path": map[string]interface{}{
				"type":        "string",
				"description": "The directory to search in. Defaults to current directory if not specified.",
			},
		},
		"required": []string{"pattern"},
	}
}

func (t *GlobTool) Execute(ctx context.Context, input map[string]interface{}, tc *tools.ToolContext) (*tools.ToolResult, error) {
	pattern, _ := input["pattern"].(string)
	path, _ := input["path"].(string)

	if pattern == "" {
		return &tools.ToolResult{
			Content: "Error: pattern is required",
			IsError: true,
		}, nil
	}

	searchPath := path
	if searchPath == "" {
		searchPath = tc.WorkingDir
	}

	if _, err := os.Stat(searchPath); os.IsNotExist(err) {
		return &tools.ToolResult{
			Content: fmt.Sprintf("Error: directory not found: %s", searchPath),
			IsError: true,
		}, nil
	}

	fullPattern := filepath.Join(searchPath, pattern)

	matches, err := filepath.Glob(fullPattern)
	if err != nil {
		return &tools.ToolResult{
			Content: fmt.Sprintf("Error: invalid glob pattern: %v", err),
			IsError: true,
		}, nil
	}

	if len(matches) == 0 {
		return &tools.ToolResult{
			Content: fmt.Sprintf("No files found matching pattern: %s", pattern),
		}, nil
	}

	sort.Strings(matches)

	resultLines := make([]string, len(matches))
	for i, match := range matches {
		relPath, err := filepath.Rel(searchPath, match)
		if err != nil {
			relPath = match
		}
		resultLines[i] = relPath
	}

	return &tools.ToolResult{
		Content: fmt.Sprintf("Found %d files:\n%s", len(matches), formatFileList(resultLines)),
	}, nil
}

func formatFileList(files []string) string {
	if len(files) <= 20 {
		return joinFiles(files)
	}

	shown := files[:20]
	remaining := len(files) - 20

	return fmt.Sprintf("%s\n... and %d more files", joinFiles(shown), remaining)
}

func joinFiles(files []string) string {
	result := ""
	for i, file := range files {
		if i > 0 {
			result += "\n"
		}
		result += fmt.Sprintf("  %s", file)
	}
	return result
}

func (t *GlobTool) IsEnabled() bool {
	return true
}

func (t *GlobTool) IsReadOnly(input map[string]interface{}) bool {
	return true
}

func (t *GlobTool) IsConcurrencySafe(input map[string]interface{}) bool {
	return true
}

func (t *GlobTool) UserFacingName(input map[string]interface{}) string {
	if pattern, ok := input["pattern"].(string); ok {
		return pattern
	}
	return "Glob"
}
