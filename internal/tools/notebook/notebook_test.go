package notebook

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNotebookTool_Execute(t *testing.T) {
	tool := New()
	if tool.Name() != "NotebookEdit" {
		t.Fatal(tool.Name())
	}
	path := filepath.Join(t.TempDir(), "n.ipynb")
	if err := os.WriteFile(path, []byte(`{"nbformat":4,"cells":[{"cell_type":"code","source":["print(1)"]}]}`), 0644); err != nil {
		t.Fatal(err)
	}
	got, err := tool.Execute(context.Background(), map[string]interface{}{"path": path, "source": "print(2)\n"}, nil)
	if err != nil || got.IsError {
		t.Fatalf("%+v %v", got, err)
	}
	saved, err := os.ReadFile(path)
	if err != nil || !strings.Contains(string(saved), "print(2)") {
		t.Fatalf("%s", saved)
	}
}
