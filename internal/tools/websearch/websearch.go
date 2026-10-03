package websearch

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/user-name/cc-cli-go/internal/tools"
)

type SearchTool struct {
	Client  *http.Client
	BaseURL string
}

func New() *SearchTool { return &SearchTool{} }

func (t *SearchTool) Name() string        { return "WebSearch" }
func (t *SearchTool) Description() string { return "Search the web through an injected HTTP endpoint." }
func (t *SearchTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{"type": "object", "properties": map[string]interface{}{"query": map[string]interface{}{"type": "string"}}}
}
func (t *SearchTool) IsEnabled() bool                               { return true }
func (t *SearchTool) IsReadOnly(map[string]interface{}) bool        { return true }
func (t *SearchTool) IsConcurrencySafe(map[string]interface{}) bool { return true }
func (t *SearchTool) UserFacingName(map[string]interface{}) string  { return "WebSearch" }

type searchResponse struct {
	Results []struct {
		Title string `json:"title"`
		URL   string `json:"url"`
	} `json:"results"`
}

func (t *SearchTool) Execute(ctx context.Context, input map[string]interface{}, _ *tools.ToolContext) (*tools.ToolResult, error) {
	query, _ := input["query"].(string)
	if query == "" {
		return &tools.ToolResult{Content: "query is required", IsError: true}, nil
	}
	if t.BaseURL == "" {
		return &tools.ToolResult{Content: "search endpoint is not configured", IsError: true}, nil
	}
	client := t.Client
	if client == nil {
		client = http.DefaultClient
	}
	endpoint := t.BaseURL + "?q=" + url.QueryEscape(query)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return &tools.ToolResult{Content: err.Error(), IsError: true}, nil
	}
	resp, err := client.Do(req)
	if err != nil {
		return &tools.ToolResult{Content: err.Error(), IsError: true}, nil
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return &tools.ToolResult{Content: err.Error(), IsError: true}, nil
	}
	var parsed searchResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return &tools.ToolResult{Content: err.Error(), IsError: true}, nil
	}
	lines := make([]string, 0, len(parsed.Results))
	for _, item := range parsed.Results {
		lines = append(lines, fmt.Sprintf("%s %s", item.Title, item.URL))
	}
	return &tools.ToolResult{Content: stringsJoin(lines)}, nil
}

func stringsJoin(lines []string) string {
	out := ""
	for i, line := range lines {
		if i > 0 {
			out += "\n"
		}
		out += line
	}
	return out
}
