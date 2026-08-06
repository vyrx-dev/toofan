package tui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/vyrx-dev/toofan/internal/game"
	"github.com/vyrx-dev/toofan/internal/lang"
	"github.com/vyrx-dev/toofan/internal/theme"
)

// --- language picker ---

func (m model) handlePicker(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.langCur > 0 {
			m.langCur--
		}
	case "down", "j":
		if m.langCur < len(lang.Names)-1 {
			m.langCur++
		}
	case "enter":
		m.lang = lang.Names[m.langCur]
		m.pickingLang = false
		m.game = game.New(m.duration, m.mode, m.lang, m.difficulty)
		m.save()
	case "esc":
		m.pickingLang = false
	}
	return m, nil
}

func (m model) viewPicker(p theme.Palette) string {
	return renderList(p, "language", lang.Names, nil, m.langCur)
}

// --- lesson picker ---

func (m model) handleLessonPicker(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	snippets := lang.GetSnippets(m.lang)
	total := len(snippets) + 1 // +1 for "Random Snippet"

	switch msg.String() {
	case "up", "k":
		if m.lessonCur > 0 {
			m.lessonCur--
		}
	case "down", "j":
		if m.lessonCur < total-1 {
			m.lessonCur++
		}
	case "enter":
		m.pickingLesson = false
		m.game = game.New(m.duration, m.mode, m.lang, m.difficulty)
		if m.lessonCur > 0 && m.lessonCur <= len(snippets) {
			m.game.Snippet = snippets[m.lessonCur-1]
			m.game.SetText(m.game.Snippet.Content)
		}
		m.save()
	case "esc":
		m.pickingLesson = false
	}
	return m, nil
}

func (m model) viewLessonPicker(p theme.Palette) string {
	snippets := lang.GetSnippets(m.lang)
	names := make([]string, len(snippets)+1)
	names[0] = "Random Snippet"
	for i, s := range snippets {
		names[i+1] = s.Topic
	}
	return renderList(p, "lesson", names, nil, m.lessonCur)
}

// --- theme picker ---

func (m model) handleThemePicker(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.themeCur > 0 {
			m.themeCur--
		}
		// live preview
		theme.Current = theme.All[m.themeCur]
	case "down", "j":
		if m.themeCur < len(theme.All)-1 {
			m.themeCur++
		}
		// live preview
		theme.Current = theme.All[m.themeCur]
	case "enter":
		theme.Current = theme.All[m.themeCur]
		m.pickingTheme = false
		m.save()
	case "esc":
		// revert to saved theme
		_, _, _, _, th, _, _ := game.LoadConfig()
		theme.Current = theme.ByName(th)
		m.pickingTheme = false
	}
	return m, nil
}

func (m model) viewThemePicker(p theme.Palette) string {
	names := make([]string, len(theme.All))
	suffixes := make([]string, len(theme.All))
	for i, t := range theme.All {
		names[i] = t.Name

		c1 := lipgloss.NewStyle().Foreground(t.Foreground).Render("●")
		c2 := lipgloss.NewStyle().Foreground(t.Typed).Render("●")
		c3 := lipgloss.NewStyle().Foreground(t.Accent).Render("●")
		c4 := lipgloss.NewStyle().Foreground(t.Error).Render("●")

		suffixes[i] = c1 + " " + c2 + " " + c3 + " " + c4
	}
	return renderList(p, "theme", names, suffixes, m.themeCur)
}

// --- duration picker ---

func (m model) viewDurPicker(p theme.Palette) string {
	durs := []string{"∞", "15", "30", "60", "120"}
	return renderList(p, "time", durs, nil, m.durCur)
}

// --- difficulty picker ---

func (m model) viewDifficultyPicker(p theme.Palette) string {
	return renderList(p, "difficulty", difficulties, nil, m.diffCur)
}

// clean list — no borders, just highlighted selection, padded uniformly
func renderList(p theme.Palette, title string, items []string, suffixes []string, cur int) string {
	dim := lipgloss.NewStyle().Foreground(p.Foreground)
	hi := lipgloss.NewStyle().Foreground(p.Accent)

	maxWidth := 0
	for _, n := range items {
		if lipgloss.Width(n) > maxWidth {
			maxWidth = lipgloss.Width(n)
		}
	}

	sel := lipgloss.NewStyle().Foreground(p.Accent).Bold(true)

	// center all lists based on the longest string to keep alignment fixed
	rows := []string{hi.Render(title), ""}

	for i, name := range items {
		var display string

		padLen := maxWidth - lipgloss.Width(name)
		display = name + strings.Repeat(" ", padLen) // explicitly pad the name so suffixes perfectly align

		if suffixes != nil && i < len(suffixes) {
			display += "  " + suffixes[i]
		}

		if i == cur {
			rows = append(rows, sel.Render(" ● "+display))
		} else {
			rows = append(rows, dim.Render(" ○ "+display))
		}
	}

	rows = append(rows, "", dim.Render("↑↓ move · enter · esc"))

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}
