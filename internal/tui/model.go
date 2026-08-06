package tui

import (
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/vyrx-dev/toofan/internal/game"
	"github.com/vyrx-dev/toofan/internal/theme"
)

type raceServerMsg struct {
	Msg game.ServerMsg
}

var durations = []int{0, 15, 30, 60, 120}
var difficulties = []string{"easy", "medium", "hard"}

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
	mode     string // "words" or "code"
	lang     string
	difficulty string

	width, height int

	pickingDur    bool
	durCur        int
	pickingLang   bool
	langCur       int
	pickingLesson bool
	lessonCur     int
	pickingTheme  bool
	themeCur      int
	pickingDifficulty bool
	diffCur           int
	showHelp      bool

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

	bots        []game.Bot
	botCfg      game.BotConfig
	pickingBots bool
	botCountCur int
	botDiffCur  int
	botPickStep int

	raceClient    *game.RaceClient
	racePlayers   []game.RacePlayer
	raceState     int
	raceText      string
	onlineCount   int
	onlineActionCur int
	onlineRoomID  string
	onlineRoomIDBuf string
	onlinePin     string
	onlinePinBuf  string
	onlineSizeCur int
	onlineSize    int
	onlineConfigCur int
	isCreating    bool
	onlineAutoStart bool
	isRaceHost    bool
	raceHostName  string
	raceCanStart  bool
	serverURL     string
	pickingOnline bool
	username      string
	usernameBuf   string
	botLastTick   time.Time
}

func New() model {
	duration, mode, language, difficulty, th, username, serverURL := game.LoadConfig()
	theme.Current = theme.ByName(th)

	return model{
		game:       game.New(duration, mode, language, difficulty),
		duration:   duration,
		mode:       mode,
		lang:       language,
		difficulty: difficulty,
		username:   username,
		serverURL:  serverURL,
		onlineAutoStart: true,
	}
}

type tick time.Time

func (m model) isPaused() bool {
	return m.pickingDur || m.pickingLang || m.pickingLesson || m.pickingTheme || m.pickingBots || m.pickingOnline || m.showHelp
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
				m.botLastTick = time.Time(msg)
			} else {
				m.game.Tick(time.Time(msg))

				if len(m.bots) > 0 && !m.botLastTick.IsZero() {
					dt := time.Time(msg).Sub(m.botLastTick)
					game.TickBots(m.bots, dt, len(m.game.Text()))
				}
				m.botLastTick = time.Time(msg)

				if m.raceState == onlineRacing && m.raceClient != nil {
					prog := 0.0
					if len(m.game.Text()) > 0 {
						prog = float64(len(m.game.Input())) / float64(len(m.game.Text()))
					}
					wpm := m.game.Stats().WPM
					m.raceClient.SendProgress(prog, wpm)
				}

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

					if m.raceState == onlineRacing {
					if m.raceClient != nil {
						m.raceClient.SendProgress(1.0, m.result.WPM)
					}
						m.raceState = onlineResults
					}
					m.active = screenResults
					m.finishedAt = time.Now()
				}
			}
		}
		return m, tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
			return tick(t)
		})

	case raceServerMsg:
		m, cmd := m.handleRaceServerMsg(msg.Msg)
		return m, cmd
	
	case joinResultMsg:
		return m.handleJoinResult(msg)
	case startRaceResultMsg:
		return m.handleStartRaceResult(msg)

	case onlineResultsDoneMsg:
		if m.raceState == onlineResults {
			m.disconnectRace()
			m.game = game.New(m.duration, m.mode, m.lang, m.difficulty)
			m.active = screenTyping
		}
		return m, nil

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
				m.game = game.New(m.duration, m.mode, m.lang, m.difficulty)
				m.save()
				return m, nil
			case "esc":
				m.pickingDur = false
				return m, nil
			default:
				m.pickingDur = false
				// fallthrough to handleTyping so the key is typed
			}
		}
		if m.pickingDifficulty {
			switch msg.String() {
			case "up", "k", "left", "h":
				if m.diffCur > 0 {
					m.diffCur--
				}
				return m, nil
			case "down", "j", "right", "l":
				if m.diffCur < len(difficulties)-1 {
					m.diffCur++
				}
				return m, nil
			case "enter":
				m.difficulty = difficulties[m.diffCur]
				m.pickingDifficulty = false
				m.game = game.New(m.duration, m.mode, m.lang, m.difficulty)
				m.save()
				return m, nil
			case "esc":
				m.pickingDifficulty = false
				return m, nil
			default:
				m.pickingDifficulty = false
			}
		}
		if m.pickingBots {
			return m.handleBotPicker(msg)
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
		if m.pickingOnline {
			return m.handleOnline(msg)
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
	game.SaveConfig(m.duration, m.mode, m.lang, m.difficulty, theme.Current.Name, m.username, m.serverURL)
}

func nextDur(cur int) int {
	for i, d := range durations {
		if d == cur {
			return durations[(i+1)%len(durations)]
		}
	}
	return 30
}
