package tui

import (
	"context"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/liao-eli/cc-cli/claude-code-go/internal/api"
	"github.com/liao-eli/cc-cli/claude-code-go/internal/query"
	"github.com/liao-eli/cc-cli/claude-code-go/internal/types"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyCtrlD:
			return m, tea.Quit
		case tea.KeyEnter:
			if m.input.Value() != "" {
				return m.submitInput()
			}
		}

	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 4
		if !m.ready {
			m.ready = true
		}

	case StreamEventMsg:
		event := query.StreamEvent(msg)
		if event.Type == "content_block_delta" && event.Delta != "" {
			if len(m.messages) > 0 {
				lastMsg := m.messages[len(m.messages)-1]
				if lastMsg.Type == types.MessageTypeAssistant {
					for i := range lastMsg.Content {
						if lastMsg.Content[i].Type == "text" {
							lastMsg.Content[i].Text += event.Delta
						}
					}
				}
			}
			m.viewport.SetContent(m.renderMessages())
		}
		return m, m.waitForEvents()

	case QueryResultMsg:
		m.loading = false
		return m, nil
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)

	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	m.spinner, cmd = m.spinner.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if !m.ready {
		return "Loading..."
	}

	var b strings.Builder

	b.WriteString(m.viewport.View())
	b.WriteString("\n\n")

	if m.loading {
		b.WriteString(m.spinner.View())
		b.WriteString(" Thinking...")
		b.WriteString("\n")
	}

	b.WriteString(m.input.View())

	return b.String()
}

func (m Model) submitInput() (tea.Model, tea.Cmd) {
	text := m.input.Value()
	m.input.SetValue("")

	m.messages = append(m.messages, types.NewUserMessage(text))

	assistantMsg := types.NewAssistantMessage()
	assistantMsg.Content = []types.ContentBlock{{Type: "text", Text: ""}}
	m.messages = append(m.messages, assistantMsg)

	m.viewport.SetContent(m.renderMessages())
	m.loading = true

	params := query.QueryParams{
		Messages:     m.messages[:len(m.messages)-1],
		SystemPrompt: []string{"You are a helpful coding assistant."},
		Model:        api.DefaultModel,
		MaxTokens:    api.DefaultMaxTokens,
	}

	m.eventChan, m.resultChan = m.QueryEngine.Query(context.Background(), params)

	return m, m.waitForEvents()
}

func (m Model) waitForEvents() tea.Cmd {
	return func() tea.Msg {
		select {
		case event, ok := <-m.eventChan:
			if !ok {
				return nil
			}
			return StreamEventMsg(event)
		case result, ok := <-m.resultChan:
			if !ok {
				return nil
			}
			return QueryResultMsg(result)
		}
	}
}

func (m Model) renderMessages() string {
	var b strings.Builder
	for _, msg := range m.messages {
		b.WriteString(m.renderMessage(msg))
		b.WriteString("\n")
	}
	return b.String()
}

func (m Model) renderMessage(msg *types.Message) string {
	var style lipgloss.Style
	var prefix string

	switch msg.Type {
	case types.MessageTypeUser:
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
		prefix = "You: "
	case types.MessageTypeAssistant:
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
		prefix = "Claude: "
	}

	var content string
	for _, block := range msg.Content {
		if block.Type == "text" {
			content += block.Text
		}
	}

	return style.Render(prefix + content)
}

type StreamEventMsg query.StreamEvent
type QueryResultMsg query.QueryResult
