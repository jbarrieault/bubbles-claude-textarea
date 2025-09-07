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
	ta.ShowLineNumbers = false
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.SetPromptFunc(2, func(lineIdx int) string {
		if lineIdx == 0 {
			return "> "
		}
		return "  "
	})
	ta.SetHeight(1)
	// ta.MaxHeight = 3

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

// This program attempts to use a textarea to make an input that grows its height dynamically based on the content.
// The dynamic height should take into account both soft and hard line breaks.

// Latest Implementation Idea (doesn't appear to work):
// - render the text area
// - measure its height, subtract the border height
// - set the height of the text area to the measured height (or clamp to max height)
// - render the text area again

// I think the flaw in my logic has to do with m.textarea.View() i'm using to measure is constrainted
// its current (...more like previous) height.

func (m model) View() string {
	layout := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		Align(lipgloss.Center)

	textArea := m.textarea.View()
	textAreaMeasuredHeight := lipgloss.Height(textArea)

	initialTextAreaHeight := m.textarea.Height()

	m.textarea.SetHeight(textAreaMeasuredHeight - 2)
	textArea = m.textarea.View()

	newTextAreaHeight := m.textarea.Height()
	newTextAreaMeasuredHeight := lipgloss.Height(textArea)

	lineInfo := m.textarea.LineInfo()
	lineCount := m.textarea.LineCount()

	debugInfo := lipgloss.NewStyle().
		Width(m.width).
		Foreground(lipgloss.Color("240")).
		Align(lipgloss.Left).
		Render(fmt.Sprintf(
			"initialTextAreaHeight: %d,\n"+
				"newTextAreaHeight: %d,\n"+
				"textAreaMeasuredHeight: %d,\n"+
				"newTextAreaMeasuredHeight: %d,\n"+
				"lineInfo.Height: %d,\n"+
				"lineCount: %d",
			initialTextAreaHeight,
			newTextAreaHeight,
			textAreaMeasuredHeight,
			newTextAreaMeasuredHeight,
			lineInfo.Height,
			lineCount,
		))

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
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
