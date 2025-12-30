package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type InputModel struct {
	textInput textinput.Model
	quitting  bool
	submitted bool
}

type InputResult struct {
	Value     string
	Submitted bool
}

func NewInputModel(defaultValue string) InputModel {
	ti := textinput.New()
	ti.Placeholder = defaultValue
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 40
	ti.PromptStyle = lipgloss.NewStyle().Foreground(Pink)
	ti.TextStyle = lipgloss.NewStyle().Foreground(White)

	return InputModel{
		textInput: ti,
	}
}

func (m InputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m InputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if m.textInput.Value() == "" {
				m.textInput.SetValue(m.textInput.Placeholder)
			}
			m.submitted = true
			m.quitting = true
			return m, tea.Quit
		case "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m InputModel) View() string {
	if m.quitting {
		return ""
	}
	title := Title.Render("Enter folder name to search:")
	input := m.textInput.View()
	hint := Hint.Render("enter to search • esc to quit")
	return fmt.Sprintf("%s\n%s\n%s", title, input, hint)
}

func RunInput(defaultValue string) (InputResult, error) {
	model := NewInputModel(defaultValue)

	p := tea.NewProgram(model)
	finalModel, err := p.Run()
	if err != nil {
		return InputResult{}, fmt.Errorf("error running input: %w", err)
	}

	m := finalModel.(InputModel)
	return InputResult{
		Value:     m.textInput.Value(),
		Submitted: m.submitted,
	}, nil
}
