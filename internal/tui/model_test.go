package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/user-name/cc-cli-go/internal/config"
	"github.com/user-name/cc-cli-go/internal/permission"
	"github.com/user-name/cc-cli-go/internal/query"
	"github.com/user-name/cc-cli-go/internal/session"
	"github.com/user-name/cc-cli-go/internal/types"
)

func TestModel_StartupViewAndKeys(t *testing.T) {
	m := InitialModel()
	if m.Init() == nil || m.View() != "Loading..." {
		t.Fatalf("view %q", m.View())
	}

	updated, _ := m.Update(tea.WindowSizeMsg{Width: 40, Height: 12})
	m = updated.(Model)
	if !m.ready || m.View() == "Loading..." {
		t.Fatal("expected ready view")
	}

	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEscape})
	m = updated.(Model)
	if cmd != nil {
		t.Fatal("escape should only clear input")
	}
	m.input.SetValue("draft")
	updated, cmd = m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	if cmd == nil {
		t.Fatal("ctrl+c should quit when idle")
	}
	_ = updated
}

func TestModel_SettingsSessionAndPermission(t *testing.T) {
	settings := config.DefaultSettings()
	sess := session.NewSession(t.TempDir())
	sess.Messages = []*types.Message{types.NewUserMessage("kept")}
	m := InitialModelWithSessionAndSettings(sess, settings)
	if len(m.messages) != 1 {
		t.Fatalf("messages %d", len(m.messages))
	}

	decision := &permission.Decision{Behavior: permission.BehaviorAsk, Reason: "ask"}
	updated, _ := m.Update(StreamEventMsg(query.StreamEvent{
		Type: "permission_request",
		PermissionRequest: &query.PermissionRequestEvent{
			ToolName: "Bash",
			Input:    map[string]interface{}{"command": "ls"},
			Decision: decision,
		},
	}))
	m = updated.(Model)
	if m.permDialog == nil || m.View() == "" {
		t.Fatal("permission dialog was not shown")
	}
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(Model)
	if m.permDialog != nil {
		t.Fatal("dialog should close after enter")
	}
}

func TestModel_CtrlDQuits(t *testing.T) {
	m := InitialModelWithSettings(config.DefaultSettings())
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlD})
	if cmd == nil {
		t.Fatal("ctrl+d should quit")
	}
}
