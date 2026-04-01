package tools

import (
	"context"

	"github.com/liao-eli/cc-cli/claude-code-go/internal/types"
)

type ToolResult struct {
	Content     interface{}
	IsError     bool
	NewMessages []*types.Message
}

type PermissionResult struct {
	Behavior string
	Reason   string
}

type ToolContext struct {
	WorkingDir  string
	AbortSignal context.Context
}

type Tool interface {
	Name() string
	Description() string
	InputSchema() map[string]interface{}

	Execute(ctx context.Context, input map[string]interface{}, tc *ToolContext) (*ToolResult, error)

	IsEnabled() bool
	IsReadOnly(input map[string]interface{}) bool
	IsConcurrencySafe(input map[string]interface{}) bool

	UserFacingName(input map[string]interface{}) string
}

func ToToolParam(t Tool) map[string]interface{} {
	return map[string]interface{}{
		"name":         t.Name(),
		"description":  t.Description(),
		"input_schema": t.InputSchema(),
	}
}
