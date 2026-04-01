package cli

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/liao-eli/cc-cli/claude-code-go/internal/api"
	"github.com/liao-eli/cc-cli/claude-code-go/internal/query"
	"github.com/liao-eli/cc-cli/claude-code-go/internal/tools"
	"github.com/liao-eli/cc-cli/claude-code-go/internal/tools/bash"
	"github.com/liao-eli/cc-cli/claude-code-go/internal/tools/edit"
	"github.com/liao-eli/cc-cli/claude-code-go/internal/tools/read"
	"github.com/liao-eli/cc-cli/claude-code-go/internal/tui"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start interactive session",
	RunE:  runInteractive,
}

func init() {
	rootCmd.AddCommand(runCmd)
}

func runInteractive(cmd *cobra.Command, args []string) error {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("ANTHROPIC_API_KEY environment variable is required")
	}

	client := api.NewClient(apiKey)

	toolReg := tools.NewRegistry()
	toolReg.Register(bash.New())
	toolReg.Register(read.New())
	toolReg.Register(edit.New())

	engine := query.NewEngine(client, toolReg)

	model := tui.InitialModel()
	model.QueryEngine = engine

	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("run TUI: %w", err)
	}

	return nil
}
