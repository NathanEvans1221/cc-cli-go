package query

import (
	"time"

	"github.com/user-name/cc-cli-go/internal/permission"
	"github.com/user-name/cc-cli-go/internal/tools"
	"github.com/user-name/cc-cli-go/internal/types"
)

type QueryParams struct {
	Messages          []*types.Message
	SystemPrompt      []string
	Tools             []tools.Tool
	Model             string
	MaxTokens         int
	MaxTurns          int
	ToolTimeout       time.Duration
	PermissionChecker *permission.Checker
}

type QueryResult struct {
	Reason string
	Error  error
}

type StreamEvent struct {
	Type    string
	Message *types.Message
	Content *types.ContentBlock
	Delta   string
	Usage   *types.Usage

	PermissionRequest *PermissionRequestEvent
}

type PermissionRequestEvent struct {
	ToolName string
	Input    map[string]interface{}
	Decision *permission.Decision
	Index    int
}
