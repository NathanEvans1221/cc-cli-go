package tui

import (
	"testing"

	"github.com/liao-eli/cc-cli-go/internal/testutil"
)

func TestNewInput(t *testing.T) {
	input := NewInput()

	if input.Value() != "" {
		t.Error("expected empty input")
	}

	if input.GetMode() != InputModeSingle {
		t.Error("expected single line mode")
	}
}

func TestInput_SetValue(t *testing.T) {
	input := NewInput()

	input.SetValue("test value")
	testutil.AssertEqual(t, "test value", input.Value())
}

func TestInput_Clear(t *testing.T) {
	input := NewInput()

	input.SetValue("test")
	input.Clear()

	testutil.AssertEqual(t, "", input.Value())
	testutil.AssertEqual(t, InputModeSingle, input.GetMode())
}

func TestInput_AddToHistory(t *testing.T) {
	input := NewInput()

	input.AddToHistory("command 1")
	input.AddToHistory("command 2")
	input.AddToHistory("command 3")

	if len(input.history) != 3 {
		t.Errorf("expected 3 history items, got %d", len(input.history))
	}
}

func TestInput_AddToHistory_NoDuplicates(t *testing.T) {
	input := NewInput()

	input.AddToHistory("command")
	input.AddToHistory("command")

	if len(input.history) != 1 {
		t.Errorf("expected 1 history item (no duplicates), got %d", len(input.history))
	}
}

func TestInput_AddToHistory_NoEmpty(t *testing.T) {
	input := NewInput()

	input.AddToHistory("")

	if len(input.history) != 0 {
		t.Errorf("expected 0 history items (no empty), got %d", len(input.history))
	}
}

func TestInput_HistoryLimit(t *testing.T) {
	input := NewInput()

	for i := 0; i < 1100; i++ {
		input.AddToHistory("command")
	}

	if len(input.history) > 1000 {
		t.Errorf("expected max 1000 history items, got %d", len(input.history))
	}
}

func TestInput_NavigateHistory_Up(t *testing.T) {
	input := NewInput()
	input.AddToHistory("command 1")
	input.AddToHistory("command 2")

	input.navigateHistory(-1)

	testutil.AssertEqual(t, "command 2", input.Value())
}

func TestInput_NavigateHistory_Down(t *testing.T) {
	input := NewInput()
	input.AddToHistory("command 1")
	input.AddToHistory("command 2")

	input.navigateHistory(-1)
	input.navigateHistory(-1)
	input.navigateHistory(1)

	testutil.AssertEqual(t, "command 2", input.Value())
}

func TestInput_NavigateHistory_Reset(t *testing.T) {
	input := NewInput()
	input.AddToHistory("command 1")
	input.AddToHistory("command 2")

	input.navigateHistory(-1)
	input.navigateHistory(1)

	testutil.AssertEqual(t, "", input.Value())
}

func TestInput_HandlePaste(t *testing.T) {
	input := NewInput()

	input.HandlePaste("pasted text")

	testutil.AssertEqual(t, "pasted text", input.Value())
	testutil.AssertEqual(t, InputModeMulti, input.GetMode())
}

func TestInput_HandlePaste_Append(t *testing.T) {
	input := NewInput()

	input.SetValue("existing ")
	input.HandlePaste("pasted text")

	testutil.AssertEqual(t, "existing pasted text", input.Value())
}

func TestInput_AdjustHeight(t *testing.T) {
	input := NewInput()

	input.SetValue("line1\nline2\nline3")
	input.AdjustHeight()

	if input.textarea.Height() != 3 {
		t.Errorf("expected height 3, got %d", input.textarea.Height())
	}
}

func TestInput_AdjustHeight_MaxLimit(t *testing.T) {
	input := NewInput()

	lines := ""
	for i := 0; i < 10; i++ {
		lines += "line\n"
	}

	input.SetValue(lines)
	input.AdjustHeight()

	if input.textarea.Height() > 5 {
		t.Errorf("expected max height 5, got %d", input.textarea.Height())
	}
}

func TestInput_Focus(t *testing.T) {
	input := NewInput()

	input.Focus()

	if !input.textarea.Focused() {
		t.Error("expected input to be focused")
	}
}

func TestInput_Blur(t *testing.T) {
	input := NewInput()

	input.Focus()
	input.Blur()

	if input.textarea.Focused() {
		t.Error("expected input to be blurred")
	}
}
