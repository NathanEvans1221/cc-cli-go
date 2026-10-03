package types

import "testing"

func TestMessageAndContentConstructors(t *testing.T) {
	user := NewUserMessage("hello")
	if user.Role != "user" || user.UUID == "" || user.Content[0].Text != "hello" {
		t.Fatalf("%+v", user)
	}
	assistant := NewAssistantMessage()
	if assistant.Role != "assistant" || assistant.Type != MessageTypeAssistant || assistant.UUID == user.UUID {
		t.Fatalf("%+v", assistant)
	}

	text := NewTextBlock("body")
	use := NewToolUseBlock("toolu", "Read", map[string]string{"path": "a"})
	result := NewToolResultBlock("toolu", "ok", false)
	if text.Type != "text" || use.Name != "Read" || result.ToolUseID != "toolu" || result.IsError {
		t.Fatalf("%+v %+v %+v", text, use, result)
	}

	usage := &Usage{InputTokens: 3, OutputTokens: 4}
	if usage.Total() != 7 {
		t.Fatalf("total %d", usage.Total())
	}
}
