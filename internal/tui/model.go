package tui

import (
	"context"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/liao-eli/cc-cli-go/internal/config"
	envctx "github.com/liao-eli/cc-cli-go/internal/context"
	"github.com/liao-eli/cc-cli-go/internal/permission"
	"github.com/liao-eli/cc-cli-go/internal/query"
	"github.com/liao-eli/cc-cli-go/internal/session"
	"github.com/liao-eli/cc-cli-go/internal/types"
)

type Model struct {
	input    Input
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

	contextInfo *envctx.ContextInfo
	session     *session.Session
	permChecker *permission.Checker
	permDialog  *PermissionDialog
}

func InitialModel() Model {
	input := NewInput()

	vp := viewport.New(80, 20)

	s := spinner.New()
	s.Spinner = spinner.Dot

	ctx, cancel := context.WithCancel(context.Background())

	contextInfo, _ := envctx.BuildContext()

	return Model{
		input:       input,
		viewport:    vp,
		spinner:     s,
		messages:    []*types.Message{},
		ctx:         ctx,
		cancel:      cancel,
		contextInfo: contextInfo,
		session:     session.NewSession(contextInfo.WorkingDir),
		permChecker: permission.NewChecker(permission.ModeDefault),
	}
}

func InitialModelWithSettings(settings *config.Settings) Model {
	input := NewInput()

	vp := viewport.New(80, 20)

	s := spinner.New()
	s.Spinner = spinner.Dot

	ctx, cancel := context.WithCancel(context.Background())
	contextInfo, _ := envctx.BuildContext()

	checker := permission.NewChecker(settings.GetPermissionMode())
	checker.SetRules(settings.ToPermissionRules())

	return Model{
		input:       input,
		viewport:    vp,
		spinner:     s,
		messages:    []*types.Message{},
		ctx:         ctx,
		cancel:      cancel,
		contextInfo: contextInfo,
		session:     session.NewSession(contextInfo.WorkingDir),
		permChecker: checker,
	}
}

func InitialModelWithSessionAndSettings(sess *session.Session, settings *config.Settings) Model {
	m := InitialModelWithSettings(settings)
	m.session = sess
	m.messages = sess.Messages
	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		spinner.Tick,
	)
}
