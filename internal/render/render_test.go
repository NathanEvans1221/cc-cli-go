package render

import (
	"strings"
	"testing"
)

func TestMarkdownHighlightAndComplete(t *testing.T) {
	out := Markdown("# Title\nuse `go` here")
	if !strings.Contains(out, "TITLE") || !strings.Contains(out, "[go]") {
		t.Fatalf("markdown %q", out)
	}
	if !strings.Contains(HighlightGo("func main() { return }"), "«func»") {
		t.Fatal("highlight")
	}
	got := Complete("Web", []string{"WebFetch", "Read", "WebSearch"})
	if len(got) != 2 {
		t.Fatalf("complete %v", got)
	}
}
