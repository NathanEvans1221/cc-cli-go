package tui

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/liao-eli/cc-cli/claude-code-go/internal/query"
	"github.com/liao-eli/cc-cli/claude-code-go/internal/types"
)

type Model struct {
	input    textinput.Model
	viewport viewport.Model
	spinner  spinner.Model

	messages []*types.Message
	loading  bool
	ready    bool

	QueryEngine *query.Engine
	eventChan   <-chan query.StreamEvent
	resultChan  <-chan query.QueryResult
}

func InitialModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Type your message..."
	ti.Focus()

	vp := viewport.New(80, 20)

	s := spinner.New()
	s.Spinner = spinner.Dot

	return Model{
		input:    ti,
		viewport: vp,
		spinner:  s,
		messages: []*types.Message{},
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		spinner.Tick,
	)
}
