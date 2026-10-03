package agent

import (
	"context"
	"testing"

	"github.com/user-name/cc-cli-go/internal/tools"
)

type runnerFunc func(context.Context, string) (string, error)

func (f runnerFunc) Run(ctx context.Context, prompt string) (string, error) {
	return f(ctx, prompt)
}

func TestAgentTool_Execute(t *testing.T) {
	tool := New(runnerFunc(func(_ context.Context, prompt string) (string, error) {
		return "did:" + prompt, nil
	}))
	if !tool.IsEnabled() || tool.Name() != "Agent" {
		t.Fatal("metadata")
	}
	got, err := tool.Execute(context.Background(), map[string]interface{}{"prompt": "ping"}, &tools.ToolContext{})
	if err != nil || got.Content != "did:ping" || got.IsError {
		t.Fatalf("%+v %v", got, err)
	}
}
