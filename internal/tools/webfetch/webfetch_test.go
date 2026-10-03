package webfetch

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFetchTool_Execute(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("page-body"))
	}))
	defer srv.Close()

	tool := New()
	if tool.Name() != "WebFetch" || !tool.IsReadOnly(nil) {
		t.Fatal("metadata")
	}
	tool.Client = srv.Client()
	got, err := tool.Execute(context.Background(), map[string]interface{}{"url": srv.URL}, nil)
	if err != nil || got.IsError || !strings.Contains(got.Content.(string), "page-body") {
		t.Fatalf("%+v %v", got, err)
	}
	missing, err := tool.Execute(context.Background(), map[string]interface{}{}, nil)
	if err != nil || !missing.IsError {
		t.Fatalf("empty url %+v", missing)
	}
}
