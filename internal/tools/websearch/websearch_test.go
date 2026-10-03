package websearch

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSearchTool_Execute(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("q") == "" {
			t.Fatal("missing query")
		}
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"results": []map[string]string{{"title": "Hit", "url": "https://example.test/a"}},
		})
	}))
	defer srv.Close()

	tool := New()
	if tool.Name() != "WebSearch" {
		t.Fatal(tool.Name())
	}
	unconfigured, err := tool.Execute(context.Background(), map[string]interface{}{"query": "cc"}, nil)
	if err != nil || !unconfigured.IsError {
		t.Fatalf("unconfigured %+v", unconfigured)
	}
	tool.Client = srv.Client()
	tool.BaseURL = srv.URL
	got, err := tool.Execute(context.Background(), map[string]interface{}{"query": "cc"}, nil)
	if err != nil || got.IsError || !strings.Contains(got.Content.(string), "Hit") {
		t.Fatalf("%+v %v", got, err)
	}
}
