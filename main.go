package main

import (
	"fmt"
	"log"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	width    int
	height   int
	textarea textarea.Model
}

func initialModel() model {
	ta := textarea.New()
	ta.Placeholder = "Placeholder..."
	ta.Focus()
	ta.MaxHeight = 3
	ta.ShowLineNumbers = false
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.SetPromptFunc(2, func(lineIdx int) string {
		if lineIdx == 0 {
			return "> "
		}
		return "  "
	})
	ta.SetHeight(1)
	ta.SetWidth(30)

	ta.FocusedStyle.Base = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		Background(lipgloss.NoColor{}).
		Align(lipgloss.Center)

	ta.BlurredStyle.Base = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		Background(lipgloss.NoColor{}).
		Align(lipgloss.Center)

	return model{
		textarea: ta,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}

	m.textarea, cmd = m.textarea.Update(msg)
	return m, cmd
}

func (m model) View() string {
	layout := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center)

	textArea := m.textarea.View()
	textAreaHeight := lipgloss.Height(textArea)
	textAreaWidth := lipgloss.Width(textArea)

	lineInfo := m.textarea.LineInfo()
	lineCount := m.textarea.LineCount()

	debugInfo := lipgloss.NewStyle().
		Width(m.width).
		Foreground(lipgloss.Color("240")).
		Align(lipgloss.Center).
		Render(fmt.Sprintf("m.width: %d, m.height: %d, textarea.Width: %d, textarea.Height: %d, textAreaHeight: %d, textAreaWidth: %d, lineInfo.Height: %d, lineCount: %d", m.width, m.height, m.textarea.Width(), m.textarea.Height(), textAreaHeight, textAreaWidth, lineInfo.Height, lineCount))

	inputContent := lipgloss.NewStyle().
		Width(m.width).
		Align(lipgloss.Center).
		Render("Input value: \"" + m.textarea.Value() + "\"")

	footer := lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center).Render(textArea)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		debugInfo,
		inputContent,
		footer,
	)

	return layout.Render(content)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
