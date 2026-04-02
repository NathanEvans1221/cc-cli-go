package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type InputMode int

const (
	InputModeSingle InputMode = iota
	InputModeMulti
)

type Input struct {
	textarea textarea.Model
	mode     InputMode
	history  []string
	histIdx  int
}

func NewInput() Input {
	ta := textarea.New()
	ta.Placeholder = "Type your message... (Shift+Enter for new line)"
	ta.Focus()
	ta.SetWidth(80)
	ta.SetHeight(1)

	return Input{
		textarea: ta,
		mode:     InputModeSingle,
		history:  []string{},
		histIdx:  -1,
	}
}

func (i *Input) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			if i.mode == InputModeSingle || i.textarea.Line() == 0 {
				return i.navigateHistory(-1)
			}

		case tea.KeyDown:
			if i.mode == InputModeSingle || i.textarea.Line() >= len(strings.Split(i.textarea.Value(), "\n"))-1 {
				return i.navigateHistory(1)
			}
		}
	}

	i.textarea, cmd = i.textarea.Update(msg)
	return cmd
}

func (i *Input) navigateHistory(direction int) tea.Cmd {
	if len(i.history) == 0 {
		return nil
	}

	if direction < 0 {
		if i.histIdx == -1 {
			i.histIdx = len(i.history) - 1
		} else if i.histIdx > 0 {
			i.histIdx--
		}
	} else {
		if i.histIdx == -1 {
			return nil
		}
		i.histIdx++
		if i.histIdx >= len(i.history) {
			i.histIdx = -1
		}
	}

	if i.histIdx == -1 {
		i.textarea.SetValue("")
	} else {
		i.textarea.SetValue(i.history[i.histIdx])
	}

	return nil
}

func (i *Input) AddToHistory(input string) {
	if input == "" {
		return
	}

	if len(i.history) > 0 && i.history[len(i.history)-1] == input {
		return
	}

	i.history = append(i.history, input)

	if len(i.history) > 1000 {
		i.history = i.history[1:]
	}

	i.histIdx = -1
}

func (i *Input) Value() string {
	return i.textarea.Value()
}

func (i *Input) SetValue(value string) {
	i.textarea.SetValue(value)
}

func (i *Input) Clear() {
	i.textarea.SetValue("")
	i.mode = InputModeSingle
	i.textarea.SetHeight(1)
}

func (i *Input) Focus() {
	i.textarea.Focus()
}

func (i *Input) Blur() {
	i.textarea.Blur()
}

func (i *Input) SetWidth(width int) {
	i.textarea.SetWidth(width)
}

func (i *Input) View() string {
	return i.textarea.View()
}

func (i *Input) GetMode() InputMode {
	return i.mode
}

func (i *Input) AdjustHeight() {
	lines := len(strings.Split(i.textarea.Value(), "\n"))
	if lines < 1 {
		lines = 1
	}
	if lines > 5 {
		lines = 5
	}
	i.textarea.SetHeight(lines)
}

func (i *Input) HandlePaste(text string) {
	currentValue := i.textarea.Value()
	newValue := currentValue + text
	i.textarea.SetValue(newValue)

	i.mode = InputModeMulti
	i.AdjustHeight()
}
