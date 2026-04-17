package tui

import (
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/vyrx-dev/toofan/internal/game"
	"github.com/vyrx-dev/toofan/internal/theme"
)

var durations = []int{0, 15, 30, 60, 120}

type screen int

const (
	screenTyping screen = iota
	screenResults
	screenProfile
)

type model struct {
	active   screen
	game     *game.Game
	duration int
	mode     string
	lang     string

	width, height int

	pickingDur    bool
	durCur        int
	pickingLang   bool
	langCur       int
	pickingLesson bool
	lessonCur     int
	pickingTheme  bool
	themeCur      int
	showHelp      bool

	pickingDiff bool
	diffCur     int
	difficulty  string

	result        game.Stats
	pb            float64
	gotNewPB      bool
	finishedAt    time.Time
	showingErrors bool

	prof profileData

	message string
	msgTime time.Time

	pickingRestore bool
	backups        []string
	restoreCur     int
}

func New() model {
	duration, mode, language, th := game.LoadConfig()
	theme.Current = theme.ByName(th)

	return model{
		game:       game.NewWithDifficulty(duration, mode, language, "easy"),
		duration:   duration,
		mode:       mode,
		lang:       language,
		difficulty: "easy",
	}
}

type tick time.Time

func (m model) isPaused() bool {
	return m.pickingDur || m.pickingLang || m.pickingLesson || m.pickingTheme || m.pickingDiff || m.showHelp
}

func (m model) Init() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return tick(t)
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tick:
		if m.active == screenTyping && m.game.Started() {
			if m.isPaused() {
				m.game.LastTick = time.Time(msg)
			} else {
				m.game.Tick(time.Time(msg))
				if m.game.Finished() {
					m.result = m.game.Stats()
					m.pb = game.GetPB(m.duration, m.mode)
					m.gotNewPB = m.result.WPM > m.pb

					durToSave := m.duration
					if durToSave == 0 {
						durToSave = m.game.TimeLeft()
					}
					game.SaveResult(m.result, durToSave, m.mode, m.lang)
					if m.gotNewPB {
						game.SavePB(m.duration, m.mode, m.result.WPM)
					}

					m.active = screenResults
					m.finishedAt = time.Now()
				}
			}
		}
		return m, tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
			return tick(t)
		})

	case tea.KeyMsg:
		if m.message != "" {
			m.message = ""
		}

		if m.pickingRestore {
			switch msg.String() {
			case "up", "k":
				if m.restoreCur > 0 {
					m.restoreCur--
				}
			case "down", "j":
				if m.restoreCur < len(m.backups)-1 {
					m.restoreCur++
				}
			case "enter":
				src := m.backups[m.restoreCur]
				if err := game.RestoreBackup(src); err == nil {
					m.message = "imported " + filepath.Base(src)
					m.msgTime = time.Now()
					m.prof = loadProfile()
				}
				m.pickingRestore = false
			case "esc", "q":
				m.pickingRestore = false
			}
			return m, nil
		}

		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "ctrl+s":
			if dest, err := game.SaveBackup(); err == nil {
				m.message = "backup saved → " + dest
				m.msgTime = time.Now()
			}
			return m, nil
		case "ctrl+r":
			files, backupDir := game.ListBackups()
			if len(files) == 0 {
				m.message = "no backups found in " + backupDir
				m.msgTime = time.Now()
				return m, nil
			}
			m.backups = files
			m.pickingRestore = true
			m.restoreCur = 0
			return m, nil
		}
		if m.showHelp {
			m.showHelp = false
			return m, nil
		}
		if m.pickingDur {
			switch msg.String() {
			case "up", "k", "left", "h":
				if m.durCur > 0 {
					m.durCur--
				}
				return m, nil
			case "down", "j", "right", "l":
				if m.durCur < len(durations)-1 {
					m.durCur++
				}
				return m, nil
			case "enter":
				m.duration = durations[m.durCur]
				m.pickingDur = false
				m.game = game.NewWithDifficulty(m.duration, m.mode, m.lang, m.difficulty)
				m.save()
				return m, nil
			case "esc":
				m.pickingDur = false
				return m, nil
			default:
				m.pickingDur = false
			}
		}
		if m.pickingLang {
			return m.handlePicker(msg)
		}
		if m.pickingLesson {
			return m.handleLessonPicker(msg)
		}
		if m.pickingTheme {
			return m.handleThemePicker(msg)
		}
		if m.pickingDiff {
			switch msg.String() {
			case "up", "k", "left", "h":
				if m.diffCur > 0 {
					m.diffCur--
				}
				return m, nil
			case "down", "j", "right", "l":
				if m.diffCur < 2 {
					m.diffCur++
				}
				return m, nil
			case "enter":
				difficulties := []string{"easy", "medium", "hard"}
				m.difficulty = difficulties[m.diffCur]
				m.pickingDiff = false
				m.game = game.NewWithDifficulty(m.duration, m.mode, m.lang, m.difficulty)
				m.save()
				return m, nil
			case "esc":
				m.pickingDiff = false
				return m, nil
			}
		}

		switch m.active {
		case screenTyping:
			return m.handleTyping(msg)
		case screenResults:
			return m.handleResults(msg)
		case screenProfile:
			return m.handleProfile(msg)
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.width == 0 {
		return ""
	}

	p := theme.Current
	var body string

	if m.pickingRestore {
		var names []string
		for _, f := range m.backups {
			names = append(names, filepath.Base(f))
		}
		body = renderList(p, "restore backup", names, nil, m.restoreCur)
	} else {
		switch m.active {
		case screenTyping:
			body = m.viewTyping(p)
		case screenProfile:
			body = m.viewProfile(p)
		case screenResults:
			body = m.viewResults(p)
		}
	}

	if m.message != "" && time.Since(m.msgTime) < 5*time.Second {
		msgStyle := lipgloss.NewStyle().Foreground(p.Background).Background(p.Success).Padding(0, 2)
		body = lipgloss.JoinVertical(lipgloss.Center,
			body,
			"", "",
			msgStyle.Render(m.message),
		)
	}

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, body)
}

func (m model) save() {
	game.SaveConfig(m.duration, m.mode, m.lang, theme.Current.Name)
}

func nextDur(cur int) int {
	for i, d := range durations {
		if d == cur {
			return durations[(i+1)%len(durations)]
		}
	}
	return 30
}