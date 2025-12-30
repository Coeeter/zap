package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/coeeter/zap/internal/scan"
)

type Mode int

const (
	ModeList Mode = iota
	ModePreview
)

type Item struct {
	Result   scan.Result
	Selected bool
}

type Model struct {
	Mode          Mode
	Items         []Item
	Cursor        int
	Width         int
	Height        int
	PreviewRoot   *PreviewNode
	PreviewCursor int
	PreviewNodes  []*PreviewNode
	LastKey       string
	Quitting      bool
	ToDelete      []string
	DeleteCalled  bool
}

type Result struct {
	ToDelete        []string
	DeleteConfirmed bool
}

func NewModel(results []scan.Result) Model {
	items := make([]Item, len(results))
	for i, r := range results {
		items[i] = Item{Result: r, Selected: false}
	}
	return Model{
		Mode:   ModeList,
		Items:  items,
		Cursor: 0,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil
	case tea.KeyMsg:
		if m.Mode == ModePreview {
			return m.updatePreview(msg)
		}
		return m.updateList(msg)
	}
	return m, nil
}

func (m Model) updateList(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	switch key {
	case "q", "esc", "ctrl+c":
		m.Quitting = true
		return m, tea.Quit
	case "up", "k", "ctrl+p":
		m.MoveUp()
		m.LastKey = ""
	case "down", "j", "ctrl+n":
		m.MoveDown()
		m.LastKey = ""
	case "g":
		if m.LastKey == "g" {
			m.MoveToTop()
			m.LastKey = ""
		} else {
			m.LastKey = "g"
		}
	case "G", "end":
		m.MoveToBottom()
		m.LastKey = ""
	case "home":
		m.MoveToTop()
		m.LastKey = ""
	case " ":
		m.ToggleCurrent()
		m.MoveDown()
		m.LastKey = ""
	case "a":
		m.SelectAll()
		m.LastKey = ""
	case "A":
		m.DeselectAll()
		m.LastKey = ""
	case "i":
		m.InvertSelection()
		m.LastKey = ""
	case "v", "l", "tab":
		m.EnterPreview()
		m.LastKey = ""
	case "enter":
		selected := m.GetSelectedPaths()
		if len(selected) > 0 {
			m.ToDelete = selected
			m.DeleteCalled = true
			m.Quitting = true
			return m, tea.Quit
		}
		m.LastKey = ""
	default:
		m.LastKey = ""
	}

	return m, nil
}

func (m Model) View() string {
	if m.Quitting {
		return ""
	}
	if m.Mode == ModePreview {
		return m.viewPreview()
	}
	return m.viewList()
}

func (m Model) viewList() string {
	cwd, _ := os.Getwd()

	title := fmt.Sprintf("Found %d folder(s)", len(m.Items))
	if count := m.SelectedCount(); count > 0 {
		title = fmt.Sprintf("Found %d folder(s) • %d selected", len(m.Items), count)
	}

	hint := "↑↓/jk move • space select • a all • v preview • enter delete • q quit"

	var content strings.Builder

	visibleHeight := m.Height - 6
	if visibleHeight < 1 {
		visibleHeight = 10
	}

	start := 0
	if m.Cursor >= visibleHeight {
		start = m.Cursor - visibleHeight + 1
	}
	end := min(start+visibleHeight, len(m.Items))

	for i := start; i < end; i++ {
		item := m.Items[i]

		cursor := "  "
		if i == m.Cursor {
			cursor = Cursor.Render("▸ ")
		}

		checkbox := "○"
		if item.Selected {
			checkbox = Selected.Render("●")
		}

		relPath, err := filepath.Rel(cwd, item.Result.Path)
		if err != nil {
			relPath = item.Result.Path
		}

		line := fmt.Sprintf("%s %s", checkbox, relPath)
		if item.Selected {
			line = Selected.Render(line)
		} else if i == m.Cursor {
			line = Cursor.Render(line)
		}

		content.WriteString(cursor)
		content.WriteString(line)
		content.WriteString("\n")
	}

	if len(m.Items) > visibleHeight {
		scrollInfo := Dim.Render(fmt.Sprintf("  %d-%d of %d", start+1, end, len(m.Items)))
		content.WriteString(scrollInfo)
		content.WriteString("\n")
	}

	return Title.Render(title) + "\n" + content.String() + Hint.Render(hint)
}

func (m Model) SelectedCount() int {
	count := 0
	for _, item := range m.Items {
		if item.Selected {
			count++
		}
	}
	return count
}

func (m Model) GetSelectedPaths() []string {
	var paths []string
	for _, item := range m.Items {
		if item.Selected {
			paths = append(paths, item.Result.Path)
		}
	}
	return paths
}

func (m *Model) ToggleCurrent() {
	if len(m.Items) > 0 && m.Cursor < len(m.Items) {
		m.Items[m.Cursor].Selected = !m.Items[m.Cursor].Selected
	}
}

func (m *Model) SelectAll() {
	for i := range m.Items {
		m.Items[i].Selected = true
	}
}

func (m *Model) DeselectAll() {
	for i := range m.Items {
		m.Items[i].Selected = false
	}
}

func (m *Model) InvertSelection() {
	for i := range m.Items {
		m.Items[i].Selected = !m.Items[i].Selected
	}
}

func (m *Model) MoveUp() {
	if m.Cursor > 0 {
		m.Cursor--
	}
}

func (m *Model) MoveDown() {
	if m.Cursor < len(m.Items)-1 {
		m.Cursor++
	}
}

func (m *Model) MoveToTop() {
	m.Cursor = 0
}

func (m *Model) MoveToBottom() {
	if len(m.Items) > 0 {
		m.Cursor = len(m.Items) - 1
	}
}

func RunSelector(results []scan.Result) (Result, error) {
	if len(results) == 0 {
		return Result{}, nil
	}

	model := NewModel(results)

	p := tea.NewProgram(model)
	finalModel, err := p.Run()
	if err != nil {
		return Result{}, fmt.Errorf("error running TUI: %w", err)
	}

	m := finalModel.(Model)
	return Result{
		ToDelete:        m.ToDelete,
		DeleteConfirmed: m.DeleteCalled,
	}, nil
}
