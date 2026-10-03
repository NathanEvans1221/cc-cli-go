package tools

import (
	"context"
	"testing"
)

type stubTool struct{ name string }

func (s stubTool) Name() string        { return s.name }
func (s stubTool) Description() string { return "stub " + s.name }
func (s stubTool) InputSchema() map[string]interface{} {
	return map[string]interface{}{"type": "object"}
}
func (s stubTool) Execute(context.Context, map[string]interface{}, *ToolContext) (*ToolResult, error) {
	return &ToolResult{Content: s.name}, nil
}
func (s stubTool) IsEnabled() bool                               { return true }
func (s stubTool) IsReadOnly(map[string]interface{}) bool        { return true }
func (s stubTool) IsConcurrencySafe(map[string]interface{}) bool { return true }
func (s stubTool) UserFacingName(map[string]interface{}) string  { return s.name }

func TestRegistryAndToolParam(t *testing.T) {
	reg := NewRegistry()
	if reg.Get("missing") != nil {
		t.Fatal("empty registry returned a tool")
	}
	reg.Register(stubTool{name: "Read"})
	reg.Register(stubTool{name: "Bash"})
	if reg.Get("Read").Name() != "Read" {
		t.Fatal("Get")
	}
	all := reg.All()
	if len(all) != 2 {
		t.Fatalf("All %d", len(all))
	}
	param := ToToolParam(all[0])
	if param["name"] == "" || param["input_schema"] == nil {
		t.Fatalf("param %#v", param)
	}
}
