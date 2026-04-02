// Package compact provides context compaction functionality for long conversations.
// When conversations exceed a token threshold, it automatically summarizes older messages
// to maintain context while reducing token usage.
package compact

import (
	"context"
	"fmt"
	"strings"

	"github.com/liao-eli/cc-cli-go/internal/types"
)

// Compactor manages conversation compaction to stay within token limits.
type Compactor struct {
	// maxTokens is the maximum allowed tokens before compaction
	maxTokens int
	// threshold is the percentage (0.0-1.0) at which compaction triggers
	threshold float64
	// summaryStyle determines how summaries are generated
	summaryStyle string
}

// CompactionResult contains the results of a compaction operation.
type CompactionResult struct {
	// OriginalTokens is the token count before compaction
	OriginalTokens int
	// CompactedTokens is the token count after compaction
	CompactedTokens int
	// Summary is the generated summary of removed messages
	Summary string
	// MessagesRemoved is the count of messages that were compacted
	MessagesRemoved int
}

// Option is a functional option for configuring the Compactor.
type Option func(*Compactor)

// WithMaxTokens sets the maximum token limit.
func WithMaxTokens(max int) Option {
	return func(c *Compactor) {
		c.maxTokens = max
	}
}

// WithThreshold sets the compaction trigger threshold.
func WithThreshold(threshold float64) Option {
	return func(c *Compactor) {
		c.threshold = threshold
	}
}

// WithSummaryStyle sets the summary generation style.
func WithSummaryStyle(style string) Option {
	return func(c *Compactor) {
		c.summaryStyle = style
	}
}

// NewCompactor creates a new Compactor with the given options.
func NewCompactor(opts ...Option) *Compactor {
	c := &Compactor{
		maxTokens:    100000,
		threshold:    0.8,
		summaryStyle: "concise",
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// ShouldCompact determines if compaction is needed based on token count.
func (c *Compactor) ShouldCompact(messages []*types.Message) bool {
	totalTokens := c.estimateTokens(messages)
	return totalTokens > int(float64(c.maxTokens)*c.threshold)
}

// Compact performs compaction on the message history.
// It keeps recent messages and generates a summary of older ones.
func (c *Compactor) Compact(messages []*types.Message) (*CompactionResult, error) {
	if len(messages) <= 2 {
		return &CompactionResult{
			OriginalTokens:  c.estimateTokens(messages),
			CompactedTokens: c.estimateTokens(messages),
			MessagesRemoved: 0,
		}, nil
	}

	// Keep recent messages (more for longer conversations)
	keepRecent := 2
	if len(messages) > 10 {
		keepRecent = 3
	}

	recentMessages := messages[len(messages)-keepRecent:]
	olderMessages := messages[:len(messages)-keepRecent]

	// Generate summary of older messages
	summary := c.generateSummary(olderMessages)

	result := &CompactionResult{
		OriginalTokens:  c.estimateTokens(messages),
		CompactedTokens: c.estimateTokens(recentMessages) + len(summary)/4,
		Summary:         summary,
		MessagesRemoved: len(olderMessages),
	}

	return result, nil
}

// generateSummary creates a concise summary of the message history.
func (c *Compactor) generateSummary(messages []*types.Message) string {
	var sb strings.Builder

	sb.WriteString("## Conversation Summary\n\n")

	// Count messages and extract information
	userCount := 0
	assistantCount := 0
	toolsUsed := make(map[string]int)
	topics := []string{}

	for _, msg := range messages {
		if msg.Role == "user" {
			userCount++
			if len(msg.Content) > 0 && msg.Content[0].Type == "text" {
				text := msg.Content[0].Text
				if len(text) > 0 {
					topic := c.extractTopic(text)
					if topic != "" {
						topics = append(topics, topic)
					}
				}
			}
		} else if msg.Role == "assistant" {
			assistantCount++
			for _, block := range msg.Content {
				if block.Type == "tool_use" {
					toolsUsed[block.Name]++
				}
			}
		}
	}

	// Format summary
	sb.WriteString(fmt.Sprintf("- %d user messages, %d assistant responses\n", userCount, assistantCount))

	if len(topics) > 0 {
		sb.WriteString("- Topics discussed: ")
		if len(topics) > 5 {
			sb.WriteString(strings.Join(topics[:5], ", "))
			sb.WriteString(fmt.Sprintf(" (and %d more)", len(topics)-5))
		} else {
			sb.WriteString(strings.Join(topics, ", "))
		}
		sb.WriteString("\n")
	}

	if len(toolsUsed) > 0 {
		sb.WriteString("- Tools used: ")
		tools := make([]string, 0, len(toolsUsed))
		for tool, count := range toolsUsed {
			tools = append(tools, fmt.Sprintf("%s (%d)", tool, count))
		}
		sb.WriteString(strings.Join(tools, ", "))
		sb.WriteString("\n")
	}

	return sb.String()
}

// extractTopic extracts key topics from user messages.
func (c *Compactor) extractTopic(text string) string {
	words := strings.Fields(text)
	if len(words) == 0 {
		return ""
	}

	keywords := []string{}
	stopWords := map[string]bool{
		"the": true, "a": true, "an": true, "is": true, "are": true,
		"was": true, "were": true, "be": true, "been": true, "being": true,
		"have": true, "has": true, "had": true, "do": true, "does": true,
		"did": true, "will": true, "would": true, "could": true, "should": true,
		"can": true, "may": true, "might": true, "must": true, "shall": true,
		"to": true, "of": true, "in": true, "for": true, "on": true,
		"with": true, "at": true, "by": true, "from": true, "as": true,
		"i": true, "you": true, "we": true, "they": true, "it": true,
		"this": true, "that": true, "these": true, "those": true,
	}

	for _, word := range words {
		lower := strings.ToLower(word)
		if !stopWords[lower] && len(word) > 2 {
			keywords = append(keywords, word)
			if len(keywords) >= 3 {
				break
			}
		}
	}

	if len(keywords) == 0 {
		return ""
	}

	return strings.Join(keywords, " ")
}

// estimateTokens provides a rough token count for messages.
// Uses simple heuristic: ~4 characters per token.
func (c *Compactor) estimateTokens(messages []*types.Message) int {
	total := 0

	for _, msg := range messages {
		for _, block := range msg.Content {
			if block.Type == "text" {
				total += len(block.Text) / 4
			} else if block.Type == "tool_use" {
				total += 50 // Base overhead for tool use
				total += len(block.Name)
			} else if block.Type == "tool_result" {
				if str, ok := block.Content.(string); ok {
					total += len(str) / 4
				}
			}
		}
	}

	return total
}

// ApplyCompaction creates a new message list with the summary applied.
func (c *Compactor) ApplyCompaction(messages []*types.Message, result *CompactionResult) []*types.Message {
	if result.MessagesRemoved == 0 {
		return messages
	}

	// Create summary message
	summaryMsg := types.NewUserMessage(result.Summary)
	summaryMsg.Role = "system"

	// Keep recent messages
	keepRecent := len(messages) - result.MessagesRemoved
	recentMessages := messages[keepRecent:]

	// Build new message list
	newMessages := []*types.Message{summaryMsg}
	newMessages = append(newMessages, recentMessages...)

	return newMessages
}

// ManualCompact performs compaction on demand.
func ManualCompact(ctx context.Context, messages []*types.Message) (*CompactionResult, error) {
	compactor := NewCompactor()
	return compactor.Compact(messages)
}
