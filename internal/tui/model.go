package tui

import (
	"context"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/user-name/cc-cli-go/internal/config"
	envctx "github.com/user-name/cc-cli-go/internal/context"
	"github.com/user-name/cc-cli-go/internal/permission"
	"github.com/user-name/cc-cli-go/internal/query"
	"github.com/user-name/cc-cli-go/internal/render"
	"github.com/user-name/cc-cli-go/internal/session"
	"github.com/user-name/cc-cli-go/internal/types"
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

	language string
	theme    string
	notice   string
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
		language:    "en",
		theme:       "12",
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

	m := Model{
		input:       input,
		viewport:    vp,
		spinner:     s,
		messages:    []*types.Message{},
		ctx:         ctx,
		cancel:      cancel,
		contextInfo: contextInfo,
		session:     session.NewSession(contextInfo.WorkingDir),
		permChecker: checker,
		language:    settings.UI.Language,
		theme:       settings.UI.Theme,
	}
	if m.language == "" {
		m.language = "en"
	}
	if m.theme == "" {
		m.theme = "12"
	}
	return m
}

func (m *Model) ApplySettings(settings *config.Settings) {
	if settings == nil {
		return
	}
	if settings.UI.Language != "" {
		m.language = settings.UI.Language
	}
	if settings.UI.Theme != "" {
		m.theme = settings.UI.Theme
	}
}

func (m Model) UserLabel() string {
	if m.language == "zh-TW" {
		return "你: "
	}
	return "You: "
}

func (m Model) ProgressLabel() string {
	if m.language == "zh-TW" {
		return "思考中..."
	}
	return "Thinking..."
}

func (m Model) Theme() string {
	if m.theme == "" {
		return "12"
	}
	return m.theme
}

func (m Model) Complete(prefix string, extra []string) []string {
	names := []string{"Bash", "Read", "Edit", "Write", "WebFetch", "WebSearch", "TodoWrite", "NotebookEdit"}
	return render.Complete(prefix, append(names, extra...))
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
