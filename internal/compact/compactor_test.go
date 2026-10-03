package compact

import (
	"context"
	"strings"
	"testing"

	"github.com/user-name/cc-cli-go/internal/types"
)

func msgs(n int) []*types.Message {
	out := make([]*types.Message, n)
	for i := 0; i < n; i++ {
		role := "user"
		if i%2 == 1 {
			role = "assistant"
		}
		text := "the user is asking about compaction tokens and history"
		if i > 6 {
			text = "more topics about files tools and sessions"
		}
		m := types.NewUserMessage(text)
		m.Role = role
		if role == "assistant" {
			m.Content = append(m.Content, types.NewToolUseBlock("id", "Read", map[string]string{}))
			m.Content = append(m.Content, types.NewToolResultBlock("id", strings.Repeat("x", 40), false))
		}
		out[i] = m
	}
	return out
}

func TestCompactor_ShortHistoryIsUnchanged(t *testing.T) {
	c := NewCompactor(WithMaxTokens(1000000), WithThreshold(0.99), WithSummaryStyle("brief"))
	messages := msgs(2)
	if c.ShouldCompact(messages) {
		t.Fatal("two short messages should stay under the limit")
	}
	result, err := c.Compact(messages)
	if err != nil {
		t.Fatal(err)
	}
	if result.MessagesRemoved != 0 {
		t.Fatalf("removed %d", result.MessagesRemoved)
	}
	if got := c.ApplyCompaction(messages, result); len(got) != len(messages) {
		t.Fatalf("applied length %d", len(got))
	}
}

func TestCompactor_SummarizesOlderMessages(t *testing.T) {
	c := NewCompactor(WithMaxTokens(20), WithThreshold(0.1))
	messages := msgs(12)
	if !c.ShouldCompact(messages) {
		t.Fatal("expected compaction")
	}
	result, err := c.Compact(messages)
	if err != nil {
		t.Fatal(err)
	}
	if result.MessagesRemoved == 0 || result.Summary == "" {
		t.Fatalf("result %+v", result)
	}
	if !strings.Contains(result.Summary, "Tools used") {
		t.Fatalf("summary %s", result.Summary)
	}
	if !strings.Contains(result.Summary, "and ") && !strings.Contains(result.Summary, "Topics") {
		t.Fatalf("summary %s", result.Summary)
	}
	applied := c.ApplyCompaction(messages, result)
	if applied[0].Role != "system" {
		t.Fatalf("first role %s", applied[0].Role)
	}
	if len(applied) != 1+len(messages)-result.MessagesRemoved {
		t.Fatalf("applied %d removed %d total %d", len(applied), result.MessagesRemoved, len(messages))
	}
}

func TestManualCompact(t *testing.T) {
	result, err := ManualCompact(context.Background(), msgs(4))
	if err != nil || result.MessagesRemoved == 0 {
		t.Fatalf("result %+v err %v", result, err)
	}
}

func TestExtractTopic_Empty(t *testing.T) {
	c := NewCompactor()
	if c.extractTopic("a to") != "" {
		t.Fatal("expected no topic")
	}
}
