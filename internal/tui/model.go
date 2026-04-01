package tui

import (
	"context"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/liao-eli/cc-cli-go/internal/query"
	"github.com/liao-eli/cc-cli-go/internal/types"
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

	ctx    context.Context
	cancel context.CancelFunc
}

func InitialModel() Model {
	ti := textinput.New()
	ti.Placeholder = "Type your message..."
	ti.Focus()

	vp := viewport.New(80, 20)

	s := spinner.New()
	s.Spinner = spinner.Dot

	ctx, cancel := context.WithCancel(context.Background())

	return Model{
		input:    ti,
		viewport: vp,
		spinner:  s,
		messages: []*types.Message{},
		ctx:      ctx,
		cancel:   cancel,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		spinner.Tick,
	)
}
