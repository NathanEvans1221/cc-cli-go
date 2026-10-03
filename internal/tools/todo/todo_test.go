package todo

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTodoTool_Execute(t *testing.T) {
	tool := New()
	if tool.Name() != "TodoWrite" || tool.IsReadOnly(nil) {
		t.Fatal("metadata")
	}
	path := filepath.Join(t.TempDir(), "todos.txt")
	got, err := tool.Execute(context.Background(), map[string]interface{}{
		"path":  path,
		"items": []interface{}{"write tests"},
	}, nil)
	if err != nil || got.IsError {
		t.Fatalf("%+v %v", got, err)
	}
	body, err := os.ReadFile(path)
	if err != nil || !strings.Contains(string(body), "write tests") {
		t.Fatalf("%s %v", body, err)
	}
	empty, err := tool.Execute(context.Background(), map[string]interface{}{}, nil)
	if err != nil || !empty.IsError {
		t.Fatalf("empty %+v", empty)
	}
}
