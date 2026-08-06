package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/vyrx-dev/toofan/internal/game"
	"github.com/vyrx-dev/toofan/internal/theme"
)

var botCounts = []int{1, 2, 3, 4, 5}
var botStrategies = []string{"personalized", "easy", "medium", "hard"}
var botColorHexes = []string{"#ff7eb6", "#42be65", "#33b1ff", "#8a3ffc", "#ff7a00"}

func (m model) handleBotPicker(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.botPickStep == 0 {
		return m.handleBotCountPicker(msg)
	}
	return m.handleBotDiffPicker(msg)
}

func (m model) handleBotCountPicker(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.botCountCur > 0 {
			m.botCountCur--
		}
	case "down", "j":
		if m.botCountCur < len(botCounts)-1 {
			m.botCountCur++
		}
	case "enter":
		m.botPickStep = 1
		m.botDiffCur = 0
	case "esc":
		m.pickingBots = false
	}
	return m, nil
}

func (m model) handleBotDiffPicker(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.botDiffCur > 0 {
			m.botDiffCur--
		}
	case "down", "j":
		if m.botDiffCur < len(botStrategies)-1 {
			m.botDiffCur++
		}
	case "enter":
		m.pickingBots = false
		m.botPickStep = 0
		cfg := game.BotConfig{
			Count:    botCounts[m.botCountCur],
			Strategy: botStrategies[m.botDiffCur],
		}
		avg := game.GetUserAvgWPM()
		m.bots = game.NewBots(cfg, avg)
		m.botCfg = cfg
		m.game.Reset(m.mode, m.lang, m.difficulty)
	case "esc":
		m.botPickStep = 0
	}
	return m, nil
}

func (m model) viewBotPicker(p theme.Palette) string {
	if m.botPickStep == 0 {
		names := make([]string, len(botCounts))
		for i, c := range botCounts {
			if c == 1 {
				names[i] = "1 bot"
			} else {
				names[i] = fmt.Sprintf("%d bots", c)
			}
		}
		return renderList(p, "race against", names, nil, m.botCountCur)
	}

	labels := make([]string, len(botStrategies))
	hints := make([]string, len(botStrategies))
	dim := lipgloss.NewStyle().Foreground(p.Foreground)
	for i, s := range botStrategies {
		labels[i] = s
		switch s {
		case "personalized":
			avg := game.GetUserAvgWPM()
			if avg > 0 {
				hints[i] = dim.Render(fmt.Sprintf("~%.0f wpm", avg))
			} else {
				hints[i] = dim.Render("not enough data")
			}
		case "easy":
			hints[i] = dim.Render("20-40 wpm")
		case "medium":
			hints[i] = dim.Render("40-80 wpm")
		case "hard":
			hints[i] = dim.Render("80-140 wpm")
		}
	}
	return renderList(p, "bot difficulty", labels, hints, m.botDiffCur)
}

func viewRaceBar(p theme.Palette, bots []game.Bot, userProgress float64, barWidth int) string {
	if barWidth < 20 {
		barWidth = 20
	}

	hi := lipgloss.NewStyle().Foreground(p.Accent)
	dim := lipgloss.NewStyle().Foreground(p.Foreground)
	val := lipgloss.NewStyle().Foreground(p.Typed)

	type racer struct {
		id       int
		name     string
		progress float64
		isUser   bool
	}

	racers := make([]racer, 0, len(bots)+1)
	racers = append(racers, racer{id: -1, name: "you", progress: userProgress, isUser: true})
	for _, b := range bots {
		racers = append(racers, racer{id: b.ID, name: b.Name, progress: b.Progress})
	}

	// sort descending by progress
	for i := 0; i < len(racers); i++ {
		for j := i + 1; j < len(racers); j++ {
			if racers[j].progress > racers[i].progress {
				racers[i], racers[j] = racers[j], racers[i]
			}
		}
	}

	nameWidth := 8
	trackWidth := barWidth - nameWidth - 8
	if trackWidth < 10 {
		trackWidth = 10
	}

	var rows []string
	for _, r := range racers {
		filled := int(r.progress * float64(trackWidth))
		if filled > trackWidth {
			filled = trackWidth
		}
		empty := trackWidth - filled

		pct := fmt.Sprintf("%3.0f%%", r.progress*100)

		var nameStr, barStr, pctStr string
		if r.isUser {
			nameStr = hi.Render(fmt.Sprintf("%-*s", nameWidth, r.name))
			barStr = hi.Render(strings.Repeat("━", filled)) + dim.Render(strings.Repeat("─", empty))
			pctStr = val.Render(pct)
		} else {
			botStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(botColorHexes[r.id%len(botColorHexes)]))
			nameStr = botStyle.Render(fmt.Sprintf("%-*s", nameWidth, r.name))
			barStr = botStyle.Render(strings.Repeat("━", filled)) + dim.Render(strings.Repeat("─", empty))
			pctStr = dim.Render(pct)
		}

		rows = append(rows, nameStr+" "+barStr+" "+pctStr)
	}

	return strings.Join(rows, "\n")
}

func viewBotResults(p theme.Palette, placements []game.BotPlacement) string {
	dim := lipgloss.NewStyle().Foreground(p.Foreground)
	hi := lipgloss.NewStyle().Foreground(p.Accent).Bold(true)

	ordinals := []string{"", "1st", "2nd", "3rd", "4th", "5th", "6th"}

	var rows []string
	rows = append(rows, "", hi.Render("race results"), "")

	for _, pl := range placements {
		ord := ""
		if pl.Rank < len(ordinals) {
			ord = ordinals[pl.Rank]
		} else {
			ord = fmt.Sprintf("%dth", pl.Rank)
		}

		wpmStr := fmt.Sprintf("%.0f wpm", pl.WPM)

		var line string
		if pl.IsUser {
			marker := hi.Render(" <")
			line = hi.Render(fmt.Sprintf("  %-4s  %-8s  %s", ord, pl.Name, wpmStr)) + marker
		} else {
			botStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(botColorHexes[pl.ID%len(botColorHexes)]))
			// Add placeholder space to match the " <" marker width
			line = dim.Render(fmt.Sprintf("  %-4s  ", ord)) + botStyle.Render(fmt.Sprintf("%-8s", pl.Name)) + dim.Render(fmt.Sprintf("  %s  ", wpmStr))
		}
		rows = append(rows, line)
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}
