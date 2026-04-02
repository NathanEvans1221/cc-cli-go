package cli

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/liao-eli/cc-cli-go/internal/api"
	"github.com/liao-eli/cc-cli-go/internal/config"
	"github.com/liao-eli/cc-cli-go/internal/query"
	"github.com/liao-eli/cc-cli-go/internal/session"
	"github.com/liao-eli/cc-cli-go/internal/tools"
	"github.com/liao-eli/cc-cli-go/internal/tools/bash"
	"github.com/liao-eli/cc-cli-go/internal/tools/edit"
	"github.com/liao-eli/cc-cli-go/internal/tools/glob"
	"github.com/liao-eli/cc-cli-go/internal/tools/grep"
	"github.com/liao-eli/cc-cli-go/internal/tools/read"
	"github.com/liao-eli/cc-cli-go/internal/tools/write"
	"github.com/liao-eli/cc-cli-go/internal/tui"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start interactive session",
	RunE:  runInteractive,
}

var continueFlag bool
var resumeFlag string

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolVarP(&continueFlag, "continue", "c", false, "Continue last session")
	runCmd.Flags().StringVar(&resumeFlag, "resume", "", "Resume specific session by ID")
}

func runInteractive(cmd *cobra.Command, args []string) error {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("ANTHROPIC_API_KEY environment variable is required")
	}

	settings, err := config.Load()
	if err != nil {
		fmt.Printf("Warning: Could not load settings: %v\n", err)
		settings = config.DefaultSettings()
	}

	if errors := settings.Validate(); len(errors) > 0 {
		fmt.Printf("Warning: Settings validation errors:\n")
		for _, e := range errors {
			fmt.Printf("  - %s\n", e.Error())
		}
	}

	client := api.NewClient(apiKey)

	toolReg := tools.NewRegistry()
	toolReg.Register(bash.New())
	toolReg.Register(read.New())
	toolReg.Register(edit.New())
	toolReg.Register(write.New())
	toolReg.Register(glob.New())
	toolReg.Register(grep.New())

	engine := query.NewEngine(client, toolReg)

	var model tui.Model
	if continueFlag || resumeFlag != "" {
		var sess *session.Session
		var err error

		if resumeFlag != "" {
			sess, err = session.LoadSession(resumeFlag)
		} else {
			sess, err = session.GetLastSession()
		}

		if err != nil {
			fmt.Printf("Warning: Could not resume session: %v\n", err)
			model = tui.InitialModelWithSettings(settings)
		} else {
			model = tui.InitialModelWithSessionAndSettings(sess, settings)
		}
	} else {
		model = tui.InitialModelWithSettings(settings)
	}

	model.QueryEngine = engine

	p := tea.NewProgram(model)
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("run TUI: %w", err)
	}

	return nil
}
