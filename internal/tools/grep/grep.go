package grep

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/liao-eli/cc-cli-go/internal/tools"
)

type GrepTool struct{}

func New() *GrepTool {
	return &GrepTool{}
}

func (t *GrepTool) Name() string {
	return "Grep"
}

func (t *GrepTool) Description() string {
	return "Fast content search tool that searches file contents using regular expressions. Supports full regex syntax."
}

func (t *GrepTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"pattern": map[string]interface{}{
				"type":        "string",
				"description": "The regular expression pattern to search for in file contents",
			},
			"path": map[string]interface{}{
				"type":        "string",
				"description": "The directory to search in. Defaults to current directory if not specified.",
			},
			"include": map[string]interface{}{
				"type":        "string",
				"description": "File pattern to include in the search (e.g. '*.js', '*.{ts,tsx}')",
			},
		},
		"required": []string{"pattern"},
	}
}

func (t *GrepTool) Execute(ctx context.Context, input map[string]interface{}, tc *tools.ToolContext) (*tools.ToolResult, error) {
	patternStr, _ := input["pattern"].(string)
	path, _ := input["path"].(string)
	includePattern, _ := input["include"].(string)

	if patternStr == "" {
		return &tools.ToolResult{
			Content: "Error: pattern is required",
			IsError: true,
		}, nil
	}

	re, err := regexp.Compile(patternStr)
	if err != nil {
		return &tools.ToolResult{
			Content: fmt.Sprintf("Error: invalid regex pattern: %v", err),
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

	includeRe := compileIncludePattern(includePattern)

	var results []GrepMatch
	err = walkDirectory(searchPath, re, includeRe, &results)
	if err != nil {
		return &tools.ToolResult{
			Content: fmt.Sprintf("Error searching files: %v", err),
			IsError: true,
		}, nil
	}

	if len(results) == 0 {
		return &tools.ToolResult{
			Content: fmt.Sprintf("No matches found for pattern: %s", patternStr),
		}, nil
	}

	output := formatGrepResults(results, searchPath)
	return &tools.ToolResult{
		Content: fmt.Sprintf("Found %d matches in %d files:\n%s", countMatches(results), countFiles(results), output),
	}, nil
}

type GrepMatch struct {
	File  string
	Line  int
	Match string
}

func compileIncludePattern(pattern string) *regexp.Regexp {
	if pattern == "" {
		return nil
	}

	globPattern := fmt.Sprintf("^%s$", regexp.QuoteMeta(pattern))
	globPattern = regexp.MustCompile(`\\\*`).ReplaceAllString(globPattern, ".*")
	globPattern = regexp.MustCompile(`\\\{([^}]+)\\\}`).ReplaceAllStringFunc(globPattern, func(m string) string {
		inner := regexp.MustCompile(`\\\{([^}]+)\\\}`).FindStringSubmatch(m)[1]
		parts := regexp.MustCompile(`,`).Split(inner, -1)
		result := "("
		for i, part := range parts {
			if i > 0 {
				result += "|"
			}
			result += part
		}
		result += ")"
		return result
	})

	re, _ := regexp.Compile(globPattern)
	return re
}

func walkDirectory(root string, pattern *regexp.Regexp, include *regexp.Regexp, results *[]GrepMatch) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if info.IsDir() {
			return nil
		}

		if include != nil {
			filename := info.Name()
			if !include.MatchString(filename) {
				return nil
			}
		}

		file, err := os.Open(path)
		if err != nil {
			return nil
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		lineNum := 1
		for scanner.Scan() {
			line := scanner.Text()
			if pattern.MatchString(line) {
				matches := pattern.FindAllString(line, -1)
				for _, match := range matches {
					*results = append(*results, GrepMatch{
						File:  path,
						Line:  lineNum,
						Match: match,
					})
				}
			}
			lineNum++
		}

		return scanner.Err()
	})
}

func formatGrepResults(results []GrepMatch, root string) string {
	if len(results) <= 50 {
		return joinMatches(results, root)
	}

	shown := results[:50]
	remaining := len(results) - 50

	return fmt.Sprintf("%s\n... and %d more matches", joinMatches(shown, root), remaining)
}

func joinMatches(matches []GrepMatch, root string) string {
	result := ""
	for i, m := range matches {
		if i > 0 {
			result += "\n"
		}
		relPath, err := filepath.Rel(root, m.File)
		if err != nil {
			relPath = m.File
		}
		result += fmt.Sprintf("%s:%d: %s", relPath, m.Line, m.Match)
	}
	return result
}

func countMatches(results []GrepMatch) int {
	return len(results)
}

func countFiles(results []GrepMatch) int {
	files := make(map[string]bool)
	for _, m := range results {
		files[m.File] = true
	}
	return len(files)
}

func (t *GrepTool) IsEnabled() bool {
	return true
}

func (t *GrepTool) IsReadOnly(input map[string]interface{}) bool {
	return true
}

func (t *GrepTool) IsConcurrencySafe(input map[string]interface{}) bool {
	return true
}

func (t *GrepTool) UserFacingName(input map[string]interface{}) string {
	if pattern, ok := input["pattern"].(string); ok {
		return pattern
	}
	return "Grep"
}
