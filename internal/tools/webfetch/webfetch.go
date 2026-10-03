package webfetch

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/user-name/cc-cli-go/internal/tools"
)

type FetchTool struct {
	Client *http.Client
}

func New() *FetchTool { return &FetchTool{} }

func (t *FetchTool) Name() string        { return "WebFetch" }
func (t *FetchTool) Description() string { return "Fetch a URL and return the response body." }
func (t *FetchTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{"type": "object", "properties": map[string]interface{}{"url": map[string]interface{}{"type": "string"}}}
}
func (t *FetchTool) IsEnabled() bool                               { return true }
func (t *FetchTool) IsReadOnly(map[string]interface{}) bool        { return true }
func (t *FetchTool) IsConcurrencySafe(map[string]interface{}) bool { return true }
func (t *FetchTool) UserFacingName(map[string]interface{}) string  { return "WebFetch" }

func (t *FetchTool) Execute(ctx context.Context, input map[string]interface{}, _ *tools.ToolContext) (*tools.ToolResult, error) {
	raw, _ := input["url"].(string)
	if raw == "" {
		return &tools.ToolResult{Content: "url is required", IsError: true}, nil
	}
	client := t.Client
	if client == nil {
		client = http.DefaultClient
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, raw, nil)
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
	if resp.StatusCode >= 400 {
		return &tools.ToolResult{Content: fmt.Sprintf("HTTP %d: %s", resp.StatusCode, strings.TrimSpace(string(body))), IsError: true}, nil
	}
	return &tools.ToolResult{Content: string(body)}, nil
}
