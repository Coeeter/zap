package tui

import "github.com/charmbracelet/lipgloss"

var (
	Pink  = lipgloss.Color("212")
	Cyan  = lipgloss.Color("39")
	Gray  = lipgloss.Color("240")
	White = lipgloss.Color("252")
	Green = lipgloss.Color("42")
	Red   = lipgloss.Color("196")

	Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(Pink).
		MarginBottom(1)

	Selected = lipgloss.NewStyle().
			Foreground(Pink).
			Bold(true)

	Cursor = lipgloss.NewStyle().
		Foreground(Cyan)

	Dim = lipgloss.NewStyle().
		Foreground(Gray)

	Hint = lipgloss.NewStyle().
		Foreground(Gray).
		MarginTop(1)

	Dir = lipgloss.NewStyle().
		Foreground(Cyan).
		Bold(true)

	File = lipgloss.NewStyle().
		Foreground(White)

	Success = lipgloss.NewStyle().
		Foreground(Green)

	Error = lipgloss.NewStyle().
		Foreground(Red)

	Spinner = lipgloss.NewStyle().
		Foreground(Pink)
)
