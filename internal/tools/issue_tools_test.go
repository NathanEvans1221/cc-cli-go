package tools_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/user-name/cc-cli-go/internal/tools"
	"github.com/user-name/cc-cli-go/internal/tools/agent"
	"github.com/user-name/cc-cli-go/internal/tools/notebook"
	"github.com/user-name/cc-cli-go/internal/tools/todo"
	"github.com/user-name/cc-cli-go/internal/tools/webfetch"
	"github.com/user-name/cc-cli-go/internal/tools/websearch"
)

func TestFetchSearchTodoNotebookAndAgent(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "search") || r.URL.Query().Get("q") != "" {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"results": []map[string]string{{"title": "Hit", "url": "https://example.test/a"}},
			})
			return
		}
		_, _ = w.Write([]byte("page-body"))
	}))
	defer srv.Close()

	fetch := webfetch.New()
	fetch.Client = srv.Client()
	got, err := fetch.Execute(context.Background(), map[string]interface{}{"url": srv.URL}, nil)
	if err != nil || got.IsError || !strings.Contains(got.Content.(string), "page-body") {
		t.Fatalf("fetch %+v %v", got, err)
	}

	search := websearch.New()
	search.Client = srv.Client()
	search.BaseURL = srv.URL
	got, err = search.Execute(context.Background(), map[string]interface{}{"query": "cc"}, nil)
	if err != nil || got.IsError || !strings.Contains(got.Content.(string), "Hit") {
		t.Fatalf("search %+v %v", got, err)
	}

	path := filepath.Join(t.TempDir(), "todos.txt")
	todoTool := todo.New()
	got, err = todoTool.Execute(context.Background(), map[string]interface{}{"path": path, "items": []interface{}{"write tests"}}, nil)
	if err != nil || got.IsError {
		t.Fatalf("todo %+v %v", got, err)
	}
	body, err := os.ReadFile(path)
	if err != nil || !strings.Contains(string(body), "write tests") {
		t.Fatalf("file %s %v", body, err)
	}

	nb := filepath.Join(t.TempDir(), "n.ipynb")
	if err := os.WriteFile(nb, []byte(`{"nbformat":4,"cells":[{"cell_type":"code","source":["print(1)"]}]}`), 0644); err != nil {
		t.Fatal(err)
	}
	got, err = notebook.New().Execute(context.Background(), map[string]interface{}{"path": nb, "source": "print(2)\n"}, nil)
	if err != nil || got.IsError {
		t.Fatalf("notebook %+v %v", got, err)
	}
	saved, err := os.ReadFile(nb)
	if err != nil || !strings.Contains(string(saved), "print(2)") {
		t.Fatalf("saved %s", saved)
	}

	agentTool := agent.New(runnerFunc(func(_ context.Context, prompt string) (string, error) {
		return "did:" + prompt, nil
	}))
	got, err = agentTool.Execute(context.Background(), map[string]interface{}{"prompt": "ping"}, &tools.ToolContext{})
	if err != nil || got.Content != "did:ping" {
		t.Fatalf("agent %+v %v", got, err)
	}
}

type runnerFunc func(context.Context, string) (string, error)

func (f runnerFunc) Run(ctx context.Context, prompt string) (string, error) {
	return f(ctx, prompt)
}
