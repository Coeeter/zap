package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

const (
	maxPreviewDepth = 3
	maxPreviewItems = 300
)

type PreviewNode struct {
	Name     string
	Path     string
	IsDir    bool
	Expanded bool
	Children []*PreviewNode
	Depth    int
}

func (m Model) updatePreview(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	key := msg.String()

	switch key {
	case "q", "esc":
		m.ExitPreview()
	case "up", "k", "ctrl+p":
		m.PreviewMoveUp()
	case "down", "j", "ctrl+n":
		m.PreviewMoveDown()
	case "enter", "l":
		m.ToggleExpand()
	case "h":
		if m.CollapseOrBack() {
			m.ExitPreview()
		}
	case "n":
		m.NextFolder()
	case "p":
		m.PrevFolder()
	}

	return m, nil
}

func (m Model) viewPreview() string {
	cwd, _ := os.Getwd()

	folderPath := ""
	if len(m.Items) > 0 && m.Cursor < len(m.Items) {
		relPath, err := filepath.Rel(cwd, m.Items[m.Cursor].Result.Path)
		if err != nil {
			relPath = m.Items[m.Cursor].Result.Path
		}
		folderPath = relPath
	}
	title := fmt.Sprintf("Preview: %s (%d/%d)", folderPath, m.Cursor+1, len(m.Items))
	hint := "↑↓/jk move • enter/l expand • h back • n/p folder • q quit"

	var content strings.Builder

	if len(m.PreviewNodes) == 0 {
		content.WriteString(Dim.Render("  (empty)"))
		content.WriteString("\n")
	} else {
		visibleHeight := m.Height - 6
		if visibleHeight < 1 {
			visibleHeight = 10
		}

		start := 0
		if m.PreviewCursor >= visibleHeight {
			start = m.PreviewCursor - visibleHeight + 1
		}
		end := min(start+visibleHeight, len(m.PreviewNodes))

		for i := start; i < end; i++ {
			node := m.PreviewNodes[i]

			cursor := "  "
			if i == m.PreviewCursor {
				cursor = Cursor.Render("▸ ")
			}

			indent := strings.Repeat("  ", node.Depth)

			var line string
			if node.IsDir {
				icon := "▶"
				if node.Expanded {
					icon = "▼"
				}
				name := fmt.Sprintf("%s %s/", icon, node.Name)
				if i == m.PreviewCursor {
					line = Cursor.Render(indent + name)
				} else {
					line = indent + Dir.Render(name)
				}
			} else {
				name := fmt.Sprintf("  %s", node.Name)
				if i == m.PreviewCursor {
					line = Cursor.Render(indent + name)
				} else {
					line = indent + File.Render(name)
				}
			}

			content.WriteString(cursor)
			content.WriteString(line)
			content.WriteString("\n")
		}

		if len(m.PreviewNodes) > visibleHeight {
			scrollInfo := Dim.Render(fmt.Sprintf("  %d-%d of %d items", start+1, end, len(m.PreviewNodes)))
			content.WriteString(scrollInfo)
			content.WriteString("\n")
		}
	}

	return Title.Render(title) + "\n" + content.String() + Hint.Render(hint)
}

func BuildPreviewTree(rootPath string) (*PreviewNode, bool) {
	root := &PreviewNode{
		Name:     filepath.Base(rootPath),
		Path:     rootPath,
		IsDir:    true,
		Expanded: true,
		Depth:    0,
	}

	itemCount := 1
	truncated := false

	var buildTree func(node *PreviewNode, depth int)
	buildTree = func(node *PreviewNode, depth int) {
		if depth >= maxPreviewDepth || itemCount >= maxPreviewItems {
			truncated = true
			return
		}

		entries, err := os.ReadDir(node.Path)
		if err != nil {
			return
		}

		sort.Slice(entries, func(i, j int) bool {
			iDir := entries[i].IsDir()
			jDir := entries[j].IsDir()
			if iDir != jDir {
				return iDir
			}
			return entries[i].Name() < entries[j].Name()
		})

		for _, entry := range entries {
			if itemCount >= maxPreviewItems {
				truncated = true
				return
			}

			child := &PreviewNode{
				Name:     entry.Name(),
				Path:     filepath.Join(node.Path, entry.Name()),
				IsDir:    entry.IsDir(),
				Expanded: false,
				Depth:    depth + 1,
			}
			node.Children = append(node.Children, child)
			itemCount++

			if entry.IsDir() && depth < 1 {
				child.Expanded = true
				buildTree(child, depth+1)
			}
		}
	}

	buildTree(root, 0)
	return root, truncated
}

func FlattenPreviewTree(root *PreviewNode) []*PreviewNode {
	if root == nil {
		return nil
	}

	var nodes []*PreviewNode
	var flatten func(node *PreviewNode)
	flatten = func(node *PreviewNode) {
		nodes = append(nodes, node)
		if node.Expanded && node.IsDir {
			for _, child := range node.Children {
				flatten(child)
			}
		}
	}
	flatten(root)
	return nodes
}

func (m *Model) ToggleExpand() {
	if m.PreviewCursor >= len(m.PreviewNodes) {
		return
	}
	node := m.PreviewNodes[m.PreviewCursor]
	if node.IsDir {
		if !node.Expanded && len(node.Children) == 0 {
			entries, err := os.ReadDir(node.Path)
			if err == nil {
				sort.Slice(entries, func(i, j int) bool {
					iDir := entries[i].IsDir()
					jDir := entries[j].IsDir()
					if iDir != jDir {
						return iDir
					}
					return entries[i].Name() < entries[j].Name()
				})
				for _, entry := range entries {
					child := &PreviewNode{
						Name:     entry.Name(),
						Path:     filepath.Join(node.Path, entry.Name()),
						IsDir:    entry.IsDir(),
						Expanded: false,
						Depth:    node.Depth + 1,
					}
					node.Children = append(node.Children, child)
				}
			}
		}
		node.Expanded = !node.Expanded
		m.PreviewNodes = FlattenPreviewTree(m.PreviewRoot)
	}
}

func (m *Model) CollapseOrBack() bool {
	if m.PreviewCursor >= len(m.PreviewNodes) {
		return true
	}
	node := m.PreviewNodes[m.PreviewCursor]
	if node.IsDir && node.Expanded {
		node.Expanded = false
		m.PreviewNodes = FlattenPreviewTree(m.PreviewRoot)
		return false
	}
	if m.PreviewCursor == 0 {
		return true
	}
	for i := m.PreviewCursor - 1; i >= 0; i-- {
		if m.PreviewNodes[i].Depth < node.Depth {
			m.PreviewCursor = i
			return false
		}
	}
	return true
}

func (m *Model) PreviewMoveUp() {
	if m.PreviewCursor > 0 {
		m.PreviewCursor--
	}
}

func (m *Model) PreviewMoveDown() {
	if m.PreviewCursor < len(m.PreviewNodes)-1 {
		m.PreviewCursor++
	}
}

func (m *Model) NextFolder() {
	if m.Cursor < len(m.Items)-1 {
		m.Cursor++
		m.EnterPreview()
	}
}

func (m *Model) PrevFolder() {
	if m.Cursor > 0 {
		m.Cursor--
		m.EnterPreview()
	}
}

func (m *Model) EnterPreview() {
	if len(m.Items) == 0 || m.Cursor >= len(m.Items) {
		return
	}
	path := m.Items[m.Cursor].Result.Path
	root, _ := BuildPreviewTree(path)
	m.PreviewRoot = root
	m.PreviewNodes = FlattenPreviewTree(root)
	m.PreviewCursor = 0
	m.Mode = ModePreview
}

func (m *Model) ExitPreview() {
	m.Mode = ModeList
	m.PreviewRoot = nil
	m.PreviewNodes = nil
	m.PreviewCursor = 0
}
