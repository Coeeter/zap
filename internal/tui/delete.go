package tui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type DeleteStatus struct {
	Path    string
	Status  string
	Error   error
	RelPath string
}

type DeleteModel struct {
	Items     []DeleteStatus
	Done      bool
	StartTime time.Time
	EndTime   time.Time
	spinner   spinner.Model
}

type DeleteResult struct {
	Deleted int
	Errors  []error
	Elapsed time.Duration
}

type deleteCompleteMsg struct {
	index int
	err   error
}

func NewDeleteModel(paths []string) DeleteModel {
	cwd, _ := os.Getwd()
	items := make([]DeleteStatus, len(paths))
	for i, p := range paths {
		relPath, err := filepath.Rel(cwd, p)
		if err != nil {
			relPath = p
		}
		items[i] = DeleteStatus{
			Path:    p,
			Status:  "pending",
			RelPath: relPath,
		}
	}

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = Spinner

	return DeleteModel{
		Items:     items,
		StartTime: time.Now(),
		spinner:   s,
	}
}

func (m DeleteModel) Init() tea.Cmd {
	cmds := []tea.Cmd{m.spinner.Tick}
	for i := range m.Items {
		cmds = append(cmds, m.startDelete(i))
	}
	return tea.Batch(cmds...)
}

func (m DeleteModel) startDelete(index int) tea.Cmd {
	return func() tea.Msg {
		path := m.Items[index].Path
		var err error

		if runtime.GOOS == "windows" {
			err = os.RemoveAll(path)
		} else {
			cmd := exec.Command("rm", "-rf", path)
			err = cmd.Run()
		}

		return deleteCompleteMsg{index: index, err: err}
	}
}

func (m DeleteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		if !m.Done {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}

	case deleteCompleteMsg:
		if msg.index < len(m.Items) {
			if msg.err != nil {
				m.Items[msg.index].Status = "error"
				m.Items[msg.index].Error = msg.err
			} else {
				m.Items[msg.index].Status = "done"
			}
		}

		allDone := true
		for _, item := range m.Items {
			if item.Status == "pending" || item.Status == "deleting" {
				allDone = false
				break
			}
		}
		if allDone {
			m.Done = true
			m.EndTime = time.Now()
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m DeleteModel) View() string {
	var b strings.Builder

	if m.Done {
		deleted := 0
		var errors []string
		for _, item := range m.Items {
			switch item.Status {
			case "done":
				deleted++
			case "error":
				errors = append(errors, fmt.Sprintf("  ✗ %s: %v", item.RelPath, item.Error))
			}
		}

		elapsed := m.EndTime.Sub(m.StartTime).Round(time.Millisecond)

		b.WriteString(Title.Render("Deletion complete"))
		b.WriteString("\n\n")

		for _, item := range m.Items {
			if item.Status == "done" {
				b.WriteString(Success.Render(fmt.Sprintf("  ✓ %s", item.RelPath)))
			} else {
				b.WriteString(Error.Render(fmt.Sprintf("  ✗ %s", item.RelPath)))
			}
			b.WriteString("\n")
		}

		b.WriteString("\n")
		summary := fmt.Sprintf("Deleted %d/%d folder(s) in %v", deleted, len(m.Items), elapsed)
		if len(errors) > 0 {
			b.WriteString(Error.Render(summary))
		} else {
			b.WriteString(Success.Render(summary))
		}
	} else {
		b.WriteString(Title.Render("Deleting..."))
		b.WriteString("\n\n")

		for _, item := range m.Items {
			switch item.Status {
			case "pending", "deleting":
				fmt.Fprintf(&b, "  %s %s\n", m.spinner.View(), item.RelPath)
			case "done":
				b.WriteString(Success.Render(fmt.Sprintf("  ✓ %s", item.RelPath)))
				b.WriteString("\n")
			case "error":
				b.WriteString(Error.Render(fmt.Sprintf("  ✗ %s", item.RelPath)))
				b.WriteString("\n")
			}
		}
	}

	return b.String()
}

func RunDelete(paths []string) (DeleteResult, error) {
	if len(paths) == 0 {
		return DeleteResult{}, nil
	}

	model := NewDeleteModel(paths)

	p := tea.NewProgram(model)
	finalModel, err := p.Run()
	if err != nil {
		return DeleteResult{}, fmt.Errorf("error running delete TUI: %w", err)
	}

	m := finalModel.(DeleteModel)

	result := DeleteResult{
		Elapsed: m.EndTime.Sub(m.StartTime),
	}
	for _, item := range m.Items {
		if item.Status == "done" {
			result.Deleted++
		} else if item.Error != nil {
			result.Errors = append(result.Errors, item.Error)
		}
	}

	return result, nil
}
