package cli

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestExecute_Version(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"--version"})
	t.Cleanup(func() {
		rootCmd.SetArgs(nil)
		rootCmd.SetOut(nil)
		rootCmd.SetErr(nil)
	})

	if err := Execute(); err != nil {
		t.Fatalf("Execute: %v", err)
	}
	if !strings.Contains(buf.String(), Version) {
		t.Fatalf("version output %q", buf.String())
	}
}

func TestRunInteractive_RequiresAPIKey(t *testing.T) {
	t.Setenv("ANTHROPIC_API_KEY", "")
	continueFlag = false
	resumeFlag = ""

	err := runInteractive(nil, nil)
	if err == nil {
		t.Fatal("expected missing API key error")
	}
	if !strings.Contains(err.Error(), "ANTHROPIC_API_KEY") {
		t.Fatalf("error = %v", err)
	}
	if os.Getenv("ANTHROPIC_API_KEY") != "" {
		t.Fatal("test left an API key set")
	}
}
