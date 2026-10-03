package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user-name/cc-cli-go/internal/tui"
)

func TestPrepareInteractive_RunPathDoesNotStartTUI(t *testing.T) {
	home := t.TempDir()
	t.Setenv("USERPROFILE", home)
	t.Setenv("HOME", home)
	t.Setenv("ANTHROPIC_API_KEY", "test-key")
	continueFlag = false
	resumeFlag = ""
	t.Cleanup(func() {
		continueFlag = false
		resumeFlag = ""
	})

	called := false
	prev := runTUI
	runTUI = func(model tui.Model) error {
		called = true
		if model.QueryEngine == nil {
			t.Fatal("query engine was not wired")
		}
		return nil
	}
	t.Cleanup(func() { runTUI = prev })

	if err := runInteractive(nil, nil); err != nil {
		t.Fatal(err)
	}
	if !called {
		t.Fatal("run did not reach the TUI hook")
	}
}

func TestPrepareInteractive_BadResumeAndInvalidSettings(t *testing.T) {
	home := t.TempDir()
	t.Setenv("USERPROFILE", home)
	t.Setenv("HOME", home)
	t.Setenv("ANTHROPIC_API_KEY", "test-key")
	if err := os.MkdirAll(filepath.Join(home, ".claude"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(home, ".claude", "settings.json"), []byte(`{"permission":{"mode":"nope","rules":[{"tool_name":"","pattern":"","behavior":"no"}]}}`), 0644); err != nil {
		t.Fatal(err)
	}
	resumeFlag = "missing-session-id"
	continueFlag = false
	t.Cleanup(func() { resumeFlag = "" })

	model, err := prepareInteractive()
	if err != nil {
		t.Fatal(err)
	}
	if model.QueryEngine == nil {
		t.Fatal("expected a model after a missing resume id")
	}
}
