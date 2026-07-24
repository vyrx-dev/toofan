package tui

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/vyrx-dev/toofan/internal/game"
	"github.com/vyrx-dev/toofan/internal/lang"
	"github.com/vyrx-dev/toofan/internal/theme"
)

func (m model) handleTyping(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.pickingDur = true
		m.durCur = 0
		for i, d := range durations {
			if d == m.duration {
				m.durCur = i
			}
		}
		return m, nil

	case "tab":
		m.game = game.New(m.duration, m.mode, m.lang, m.difficulty)
		if m.activeRace != nil {
			m.game.SetText(m.activeRace.Text)
		}

	case "ctrl+d":
		if m.mode == "words" {
			m.activeRace = nil
			m.pickingDifficulty = true
			m.diffCur = 0
			for i, d := range difficulties {
				if d == m.difficulty {
					m.diffCur = i
				}
			}
		}
		return m, nil

	case "ctrl+w":
		if !m.game.Started() {
			m.activeRace = nil
			if m.mode == "words" {
				m.mode = "code"
			} else {
				m.mode = "words"
			}
			m.game = game.New(m.duration, m.mode, m.lang, m.difficulty)
			m.save()
		} else {
			m.game.BackspaceWord()
		}

	case "ctrl+l":
		if m.mode == "code" && !m.game.Started() {
			m.activeRace = nil
			m.pickingLang = true
			m.langCur = 0
			for i, name := range lang.Names {
				if name == m.lang {
					m.langCur = i
				}
			}
		}

	case "ctrl+o":
		if m.mode == "code" && !m.game.Started() {
			m.activeRace = nil
			m.pickingLesson = true
			m.lessonCur = 0
		}

	case "ctrl+t":
		m.pickingTheme = true
		m.themeCur = 0
		for i, t := range theme.All {
			if t.Name == theme.Current.Name {
				m.themeCur = i
			}
		}

	case "ctrl+h":
		if !m.game.Started() {
			m.showHelp = true
			return m, nil
		}

	case "?":
		if !m.game.Started() {
			m.showHelp = true
			return m, nil
		}
		m.game.TypeChar('?')

	case "ctrl+p":
		if !m.game.Started() {
			m.prof = loadProfile()
			m.active = screenProfile
			return m, nil
		}
	case "ctrl+g":
		if !m.game.Started() {
			m.races = game.LoadRaces()
			if len(m.races) == 0 {
				m.message = "no saved races yet"
				m.msgTime = time.Now()
				return m, nil
			}
			m.raceCur = 0
			m.pickingRace = true
			return m, nil
		}

	case "enter":
		m.game.TypeChar('\n')

	case "backspace":
		m.game.Backspace()

	default:
		for _, r := range msg.Runes {
			m.game.TypeChar(r)
		}
	}

	return m, nil
}

func (m model) viewTyping(p theme.Palette) string {
	if m.pickingDur {
		return m.viewDurPicker(p)
	}
	if m.showHelp {
		return m.viewHelp(p)
	}
	if m.pickingLang {
		return m.viewPicker(p)
	}
	if m.pickingLesson {
		return m.viewLessonPicker(p)
	}
	if m.pickingTheme {
		return m.viewThemePicker(p)
	}
	if m.pickingDifficulty {
		return m.viewDifficultyPicker(p)
	}

	dim := lipgloss.NewStyle().Foreground(p.Foreground)
	hi := lipgloss.NewStyle().Foreground(p.Accent)

	textWidth := min(m.width-8, 72)
	textWidth = max(textWidth, 40)

	lines := splitLines(m.game.Text(), textWidth, m.game.CodeMode)
	curLine := cursorLine(lines, len(m.game.Input()))
	ghostPos := -1
	if m.activeRace != nil && m.game.Started() {
		ghostPos = ghostLenAt(m.activeRace.Points, int(m.game.Elapsed().Milliseconds()))
	}

	// word mode: 3 lines (monkeytype style), code mode: 7 lines (full snippet)
	visible := 3
	if m.game.CodeMode {
		visible = 7
	}

	// Paginate the display so lines don't instantly scroll up every time you hit enter.
	// This keeps code snippets fixed in view until you finish the whole block.
	top := (curLine / visible) * visible
	bot := top + visible
	if bot > len(lines) {
		bot = len(lines)
	}

	text := lipgloss.NewStyle().
		Padding(0, 2).
		Render(colorText(m.game, p, lines, top, bot, ghostPos))

	// timer at top
	var topLine string
	if m.game.Started() {
		timeLeft := m.game.TimeLeft()
		if m.game.Elapsed().Seconds() >= 3 {
			wpm := m.game.Stats().WPM
			topLine = hi.Render(fmt.Sprintf("%d", timeLeft)) + dim.Render(fmt.Sprintf("   %.0f wpm", wpm))
		} else {
			topLine = hi.Render(fmt.Sprintf("%d", timeLeft))
		}

	} else {
		if m.duration == 0 {
			topLine = hi.Render("∞")
		} else {
			topLine = hi.Render(fmt.Sprintf("%d", m.duration))
		}
	}

	var out []string
	out = append(out, topLine, "")

	if m.mode == "code" && m.game.Snippet.Topic != "" {
		topicLine := dim.Render("/* " + m.game.Snippet.Topic + " */")
		topicLine = lipgloss.NewStyle().Padding(0, 2).Render(topicLine)
		out = append(out, topicLine, "")
	}

	out = append(out, text)

	// progress bar when typing
	if m.game.Started() {
		ratio := min(m.game.Stats().WPM/200.0, 1.0) // rough estimate
		if m.game.Duration() > 0 {
			ratio = min(float64(m.game.TimeLeft())/float64(m.game.Duration()), 1.0)
			ratio = 1.0 - ratio // invert for progress
		} else {
			if len(m.game.Text()) > 0 {
				ratio = min(float64(len(m.game.Input()))/float64(len(m.game.Text())), 1.0)
			}
		}
		barWidth := textWidth - 4
		filled := int(ratio * float64(barWidth))
		bar := lipgloss.NewStyle().Foreground(p.Accent).Render(strings.Repeat("━", filled)) +
			lipgloss.NewStyle().Foreground(p.Foreground).Render(strings.Repeat("─", barWidth-filled))
		out = append(out, "", bar)

	} else {
		// before test: subtle info line
		var modeLabel string
		if m.mode == "code" {
			modeLabel = "code (" + m.lang + ")"
		} else {
			modeLabel = "words"
		}
		infoStr := modeLabel + " · ? help"
		if m.activeRace != nil {
			infoStr += fmt.Sprintf(" · old %.0f", m.activeRace.Stats.WPM)
		}
		info := dim.Render(infoStr)
		out = append(out, "", info)
	}

	body := lipgloss.JoinVertical(lipgloss.Center, out...)

	return body
}

func (m model) viewHelp(p theme.Palette) string {
	dim := lipgloss.NewStyle().Foreground(p.Foreground)
	hi := lipgloss.NewStyle().Foreground(p.Accent)
	val := lipgloss.NewStyle().Foreground(p.Typed).Padding(0, 1)

	lines := []string{
		hi.Render("keybinds"),
		"",
		val.Render("ctrl+w") + dim.Render("    toggle words/code"),
		val.Render("ctrl+l") + dim.Render("    change language (code mode only)"),
		val.Render("ctrl+o") + dim.Render("    change lesson (code mode only)"),
		val.Render("ctrl+t") + dim.Render("    change theme"),
		val.Render("ctrl+p") + dim.Render("    open profile"),
		val.Render("ctrl+d") + dim.Render("    change difficulty (words mode only)"),
		val.Render("ctrl+g") + dim.Render("    race against previous run"),
		val.Render("ctrl+c") + dim.Render("    exit toofan"),
		val.Render("tab") + dim.Render("       restart test"),
		val.Render("esc") + dim.Render("       configure duration"),
		val.Render("e") + dim.Render("         view error words (results screen)"),
		val.Render("?") + dim.Render("         show this help"),
		"",
		dim.Render("press any key to close"),
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

func ghostLenAt(points []game.RacePoint, elapsedMS int) int {
	if len(points) == 0 {
		return 0
	}
	out := 0
	for _, p := range points {
		if p.MS > elapsedMS {
			break
		}
		out = p.Len
	}
	return out
}
