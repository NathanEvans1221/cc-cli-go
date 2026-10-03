package tui

import (
	"context"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/user-name/cc-cli-go/internal/api"
	"github.com/user-name/cc-cli-go/internal/config"
	apperr "github.com/user-name/cc-cli-go/internal/errors"
	"github.com/user-name/cc-cli-go/internal/permission"
	"github.com/user-name/cc-cli-go/internal/query"
	"github.com/user-name/cc-cli-go/internal/tools"
	"github.com/user-name/cc-cli-go/internal/types"
)

type stopStreamer struct{}

func (stopStreamer) Stream(context.Context, *api.Request) (<-chan api.StreamEvent, error) {
	ch := make(chan api.StreamEvent, 1)
	ch <- api.StreamEvent{Type: "message_stop"}
	close(ch)
	return ch, nil
}

func TestModel_SubmitShowsProgressAndLanguage(t *testing.T) {
	t.Setenv("USERPROFILE", t.TempDir())
	t.Setenv("HOME", t.TempDir())
	settings := config.DefaultSettings()
	settings.UI.Language = "zh-TW"
	settings.UI.Theme = "10"
	m := InitialModelWithSettings(settings)
	if m.UserLabel() != "你: " || m.Theme() != "10" || m.ProgressLabel() != "思考中..." {
		t.Fatalf("label %s theme %s progress %s", m.UserLabel(), m.Theme(), m.ProgressLabel())
	}
	m.QueryEngine = query.NewEngineWithStreamer(stopStreamer{}, tools.NewRegistry())
	m.input.SetValue("hello # Hi")
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(Model)
	if !m.loading || cmd == nil {
		t.Fatal("submit did not start a query")
	}
	sized, _ := m.Update(tea.WindowSizeMsg{Width: 60, Height: 20})
	m = sized.(Model)
	if !strings.Contains(m.View(), "思考中...") {
		t.Fatalf("view %q", m.View())
	}
	done := make(chan query.QueryResult, 1)
	done <- query.QueryResult{Reason: "completed"}
	m.resultChan = done
	m.eventChan = nil
	updated, _ = m.Update(m.waitForEvents()())
	m = updated.(Model)
	if m.loading {
		t.Fatal("query result should clear loading")
	}
	matches := m.Complete("Web", []string{"docs/api.md"})
	if len(matches) < 2 {
		t.Fatalf("matches %v", matches)
	}
}

func TestModel_ErrorNoticeAndCancel(t *testing.T) {
	t.Setenv("USERPROFILE", t.TempDir())
	t.Setenv("HOME", t.TempDir())
	m := InitialModel()
	m.loading = true
	updated, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	m = updated.(Model)
	if m.loading || cmd != nil {
		t.Fatal("ctrl+c while loading should cancel without quitting")
	}
	updated, _ = m.Update(tea.WindowSizeMsg{Width: 40, Height: 10})
	m = updated.(Model)
	err := apperr.New(apperr.ErrorTypeAPI, "AUTH_ERROR", "登入失敗").WithSuggestion("檢查 API key")
	updated, _ = m.Update(QueryResultMsg(query.QueryResult{Reason: "error", Error: err}))
	m = updated.(Model)
	if !strings.Contains(m.View(), "檢查 API key") {
		t.Fatalf("view %q", m.View())
	}
	assistant := types.NewAssistantMessage()
	assistant.Content = []types.ContentBlock{{Type: "text", Text: ""}}
	m.messages = []*types.Message{assistant}
	updated, cmd = m.Update(StreamEventMsg(query.StreamEvent{Type: "content_block_delta", Delta: "token"}))
	m = updated.(Model)
	if cmd == nil || m.messages[0].Content[0].Text != "token" {
		t.Fatalf("delta %+v cmd %v", m.messages[0].Content[0], cmd != nil)
	}
}

func TestPermissionDialogViewAndApplySettings(t *testing.T) {
	t.Setenv("USERPROFILE", t.TempDir())
	t.Setenv("HOME", t.TempDir())
	m := InitialModel()
	m.ApplySettings(&config.Settings{UI: config.UISettings{Language: "en", Theme: "14"}})
	if m.Theme() != "14" || m.UserLabel() != "You: " {
		t.Fatalf("theme %s label %s", m.Theme(), m.UserLabel())
	}
	m.ApplySettings(nil)
	m.input.SetWidth(40)
	dialog := NewPermissionDialog("Bash", map[string]interface{}{"command": strings.Repeat("a", 120)}, &permission.Decision{Behavior: permission.BehaviorAsk, Reason: "dangerous command"})
	updated, _ := dialog.Update(tea.KeyMsg{Type: tea.KeyLeft})
	if view := updated.View(); !strings.Contains(view, "Bash") || !strings.Contains(view, "dangerous") {
		t.Fatalf("dialog view %s", view)
	}
	m.permDialog = &updated
	m.ready = true
	if !strings.Contains(m.View(), "Permission Request") {
		t.Fatal("model view did not show the dialog")
	}
	updated, _ = updated.Update(tea.KeyMsg{Type: tea.KeyEscape})
	if ok, action := updated.GetDecision(); ok || action != "Allow" {
		t.Fatalf("escape decision %v %s", ok, action)
	}
	events := make(chan query.StreamEvent)
	close(events)
	m.eventChan = events
	m.resultChan = nil
	if m.waitForEvents()() != nil {
		t.Fatal("closed event channel should yield nil")
	}
}

func TestPermissionAlwaysAllow(t *testing.T) {
	t.Setenv("USERPROFILE", t.TempDir())
	t.Setenv("HOME", t.TempDir())
	m := InitialModel()
	updated, _ := m.Update(StreamEventMsg(query.StreamEvent{
		Type: "permission_request",
		PermissionRequest: &query.PermissionRequestEvent{
			ToolName: "Bash",
			Input:    map[string]interface{}{"command": "ls"},
			Decision: nil,
		},
	}))
	m = updated.(Model)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRight})
	m = updated.(Model)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyRight})
	m = updated.(Model)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updated.(Model)
	if m.permDialog != nil {
		t.Fatal("dialog still open")
	}
}
