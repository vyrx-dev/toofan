package tui

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/vyrx-dev/toofan/internal/game"
	"github.com/vyrx-dev/toofan/internal/theme"
)

const (
	onlineOff        = 0
	onlineActionPick = 1
	onlineSizePick   = 2
	onlineRoomIDInput = 3
	onlinePinInput    = 4
	onlineUsername   = 5
	onlineConnecting = 6
	onlineLobby      = 7
	onlineCountdown  = 8
	onlineRacing     = 9
	onlineResults    = 10
	onlineConfigPick = 11
)

type joinResultMsg struct {
	err error
}

type onlineResultsDoneMsg struct{}
type startRaceResultMsg struct {
	err error
}

var onlineSizes = []int{2, 3, 4, 5, 6}
var onlineActions = []string{"No Rooms (Quick Match)", "Join Room", "Create Room"}

func (m model) handleOnline(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch m.raceState {
	case onlineActionPick:
		return m.handleActionPicker(msg)
	case onlineSizePick:
		return m.handleSizePicker(msg)
	case onlineRoomIDInput:
		return m.handleRoomIDInput(msg)
	case onlinePinInput:
		return m.handlePinInput(msg)
	case onlineConfigPick:
		return m.handleConfigPicker(msg)
	case onlineUsername:
		return m.handleUsernameInput(msg)
	case onlineLobby, onlineCountdown:
		if msg.String() == "esc" {
			m.disconnectRace()
			return m, nil
		}
		if m.raceState == onlineLobby && (msg.String() == "s" || msg.String() == "enter") {
			if m.raceClient != nil && m.isRaceHost && !m.onlineAutoStart && m.raceCanStart {
				return m, m.startRaceCmd()
			}
		}
	case onlineResults:
		if msg.String() == "esc" || msg.String() == "enter" {
			m.disconnectRace()
			return m, nil
		}
	}
	return m, nil
}

func (m model) handleActionPicker(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.onlineActionCur > 0 {
			m.onlineActionCur--
		}
	case "down", "j":
		if m.onlineActionCur < len(onlineActions)-1 {
			m.onlineActionCur++
		}
	case "enter":
		m.isCreating = (m.onlineActionCur == 2)
		m.onlineRoomID = ""
		m.onlineRoomIDBuf = ""
		m.onlinePin = ""
		m.onlinePinBuf = ""
		if m.isCreating {
			m.raceState = onlinePinInput
		} else if m.onlineActionCur == 1 {
			m.raceState = onlineRoomIDInput
		} else {
			m.raceState = onlineUsername
		}
	case "esc":
		m.pickingOnline = false
		m.raceState = onlineOff
	}
	return m, nil
}

func (m model) handleSizePicker(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.onlineSizeCur > 0 {
			m.onlineSizeCur--
		}
	case "down", "j":
		if m.onlineSizeCur < len(onlineSizes)-1 {
			m.onlineSizeCur++
		}
	case "enter":
		m.onlineSize = onlineSizes[m.onlineSizeCur]
		m.raceState = onlineConfigPick
		m.onlineConfigCur = 0
	case "esc":
		m.raceState = onlinePinInput
	}
	return m, nil
}

func (m model) getOnlineConfigOptions() []string {
	opts := []string{"Mode"}
	if m.mode == "code" {
		opts = append(opts, "Language")
	}
	opts = append(opts, "Difficulty", "Duration", "Start Policy", "Continue")
	return opts
}

func (m model) handleConfigPicker(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	opts := m.getOnlineConfigOptions()
	switch msg.String() {
	case "up", "k":
		if m.onlineConfigCur > 0 {
			m.onlineConfigCur--
		}
	case "down", "j":
		if m.onlineConfigCur < len(opts)-1 {
			m.onlineConfigCur++
		}
	case "enter":
		opt := opts[m.onlineConfigCur]
		switch opt {
		case "Mode":
			if m.mode == "words" {
				m.mode = "code"
			} else {
				m.mode = "words"
				m.lang = "english"
			}
			// Reset cursor if it was on a hidden option
			if m.onlineConfigCur >= len(m.getOnlineConfigOptions()) {
				m.onlineConfigCur = len(m.getOnlineConfigOptions()) - 1
			}
		case "Language":
			m.pickingLang = true
		case "Difficulty":
			m.pickingDifficulty = true
		case "Duration":
			m.duration = nextDur(m.duration)
		case "Start Policy":
			m.onlineAutoStart = !m.onlineAutoStart
		case "Continue":
			m.raceState = onlineUsername
		}
	case "esc":
		m.raceState = onlineSizePick
	}
	return m, nil
}

func (m model) handleRoomIDInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.raceState = onlineActionPick
		m.onlineRoomIDBuf = ""
		return m, nil
	case "enter":
		m.onlineRoomID = strings.ToUpper(strings.TrimSpace(m.onlineRoomIDBuf))
		m.raceState = onlinePinInput
		return m, nil
	case "backspace":
		if len(m.onlineRoomIDBuf) > 0 {
			m.onlineRoomIDBuf = m.onlineRoomIDBuf[:len(m.onlineRoomIDBuf)-1]
		}
	default:
		for _, r := range msg.Runes {
			if len(m.onlineRoomIDBuf) < 8 {
				m.onlineRoomIDBuf += string(r)
			}
		}
	}
	return m, nil
}

func (m model) handlePinInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		if m.isCreating {
			m.raceState = onlineActionPick
		} else {
			m.raceState = onlineRoomIDInput
		}
		m.onlinePinBuf = ""
		return m, nil
	case "enter":
		m.onlinePin = strings.TrimSpace(m.onlinePinBuf)
		if m.isCreating {
			m.raceState = onlineSizePick
		} else {
			m.raceState = onlineUsername
		}
		return m, nil
	case "backspace":
		if len(m.onlinePinBuf) > 0 {
			m.onlinePinBuf = m.onlinePinBuf[:len(m.onlinePinBuf)-1]
		}
	default:
		for _, r := range msg.Runes {
			if len(m.onlinePinBuf) < 8 {
				m.onlinePinBuf += string(r)
			}
		}
	}
	return m, nil
}

func (m model) handleUsernameInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		if m.isCreating {
			m.raceState = onlineConfigPick
		} else {
			m.raceState = onlinePinInput
		}
		m.usernameBuf = ""
		return m, nil
	case "enter":
		name := strings.TrimSpace(m.usernameBuf)
		if name == "" && m.username != "" {
			name = m.username
		}
		if len(name) < 2 || len(name) > 16 {
			return m, nil
		}
		m.username = name
		m.usernameBuf = ""
		m.save()
		m.raceState = onlineConnecting

		serverURL := m.serverURL
		if serverURL == "" {
			serverURL = game.DefaultServerURL
		}
		m.raceClient = game.NewRaceClient(serverURL, m.username)
		
		return m, m.joinRaceCmd()
	case "backspace":
		if len(m.usernameBuf) > 0 {
			m.usernameBuf = m.usernameBuf[:len(m.usernameBuf)-1]
		}
	default:
		for _, r := range msg.Runes {
			if len(m.usernameBuf) < 16 {
				m.usernameBuf += string(r)
			}
		}
	}
	return m, nil
}

func (m *model) disconnectRace() {
	if m.raceClient != nil {
		m.raceClient.Close()
		m.raceClient = nil
	}
	m.pickingOnline = false
	m.raceState = onlineOff
	m.racePlayers = nil
	m.raceText = ""
	m.onlineRoomIDBuf = ""
	m.onlinePinBuf = ""
	m.isRaceHost = false
	m.raceHostName = ""
	m.raceCanStart = false
}

func (m model) listenRaceMsg() tea.Cmd {
	if m.raceClient == nil {
		return nil
	}
	rc := m.raceClient
	return func() tea.Msg {
		msg, ok := <-rc.Messages()
		if !ok {
			return raceServerMsg{Msg: game.ServerMsg{Type: "disconnected"}}
		}
		return raceServerMsg{Msg: msg}
	}
}

func (m model) joinRaceCmd() tea.Cmd {
	return func() tea.Msg {
		err := m.raceClient.Join(m.onlineRoomID, m.onlinePin, m.isCreating, m.onlineSize, m.difficulty, m.mode, m.lang, m.duration, m.onlineAutoStart)
		return joinResultMsg{err: err}
	}
}

func (m model) startRaceCmd() tea.Cmd {
	return func() tea.Msg {
		if m.raceClient == nil {
			return startRaceResultMsg{err: fmt.Errorf("not connected")}
		}
		return startRaceResultMsg{err: m.raceClient.StartRace()}
	}
}

func (m model) handleRaceServerMsg(msg game.ServerMsg) (model, tea.Cmd) {
	if m.raceClient == nil {
		return m, nil
	}

	switch msg.Type {
	case "joined":
		var payload game.LobbyPayload
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			return m, m.listenRaceMsg()
		}
		m.raceClient.SetRoom(payload.Room)
		m.onlineRoomID = payload.Room
		m.onlineCount = payload.Online
		m.raceHostName = payload.Host
		m.onlineAutoStart = payload.AutoStart
		m.raceCanStart = payload.CanStart
		m.isRaceHost = payload.Host == m.username
		
		// sync settings from creator
		if !m.isCreating {
			m.difficulty = payload.Difficulty
			m.mode = payload.Mode
			if payload.Mode == "words" {
				m.lang = "english"
			} else {
				m.lang = payload.Lang
			}
			m.duration = payload.Duration
		}

		m.racePlayers = make([]game.RacePlayer, len(payload.Players))
		for i, name := range payload.Players {
			m.racePlayers[i] = game.RacePlayer{
				Name:   name,
				IsUser: name == m.username,
			}
		}
		
		// sync current state
		switch payload.State {
		case "countdown":
			m.raceState = onlineCountdown
		case "racing":
			m.raceText = payload.Text
			m.game.Reset(m.mode, m.lang, m.difficulty)
			if m.raceText != "" {
				m.game.SetText(m.raceText)
			}
			m.raceState = onlineRacing
			m.pickingOnline = false
		default:
			m.raceState = onlineLobby
		}

	case "countdown":
		var payload game.CountdownPayload
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			return m, m.listenRaceMsg()
		}
		m.raceState = onlineCountdown

	case "start":
		var payload game.StartPayload
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			return m, m.listenRaceMsg()
		}
		m.raceText = payload.Text
		
		// sync final settings
		m.difficulty = payload.Difficulty
		m.mode = payload.Mode
		if payload.Mode == "words" {
			m.lang = "english"
		} else {
			m.lang = payload.Lang
		}
		m.duration = payload.Duration
		
		m.game.Reset(m.mode, m.lang, m.difficulty)
		if payload.Text != "" {
			m.game.SetText(payload.Text)
		}
		m.raceState = onlineRacing
		m.pickingOnline = false

	case "progress":
		var payload game.ProgressPayload
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			return m, m.listenRaceMsg()
		}
		for i := range payload.Players {
			if m.raceClient != nil && payload.Players[i].Name == m.raceClient.Name() {
				payload.Players[i].IsUser = true
			}
		}
		m.racePlayers = payload.Players
		m.raceCanStart = false

	case "finish":
		// Only transition to results if we were actually part of the race
		if m.raceState != onlineRacing && m.raceState != onlineCountdown {
			return m, m.listenRaceMsg()
		}

		var payload game.FinishPayload
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			return m, m.listenRaceMsg()
		}
		for i := range payload.Placements {
			if m.raceClient != nil && payload.Placements[i].Name == m.raceClient.Name() {
				payload.Placements[i].IsUser = true
			}
		}
		m.racePlayers = payload.Placements
		m.raceState = onlineResults
		m.active = screenResults
		m.finishedAt = time.Now()
		return m, tea.Tick(4*time.Second, func(time.Time) tea.Msg { return onlineResultsDoneMsg{} })

	case "online":
		var payload game.OnlinePayload
		if err := json.Unmarshal(msg.Payload, &payload); err != nil {
			return m, m.listenRaceMsg()
		}
		m.onlineCount = payload.Count

	case "disconnected":
		m.disconnectRace()
		m.message = "disconnected from server"
		return m, nil
	}

	return m, m.listenRaceMsg()
}

func (m model) handleJoinResult(msg joinResultMsg) (model, tea.Cmd) {
	if msg.err != nil {
		m.message = "failed: " + msg.err.Error()
		m.msgTime = time.Now()
		m.pickingOnline = false
		m.raceState = onlineOff
		if m.raceClient != nil {
			m.raceClient.Close()
		}
		m.raceClient = nil
		return m, nil
	}
	m.raceState = onlineLobby
	return m, m.listenRaceMsg()
}

func (m model) handleStartRaceResult(msg startRaceResultMsg) (model, tea.Cmd) {
	if msg.err != nil {
		m.message = msg.err.Error()
		m.msgTime = time.Now()
	}
	return m, nil
}

func (m model) viewOnline(p theme.Palette) string {
	switch m.raceState {
	case onlineOff, onlineActionPick:
		return m.viewActionPicker(p)
	case onlineSizePick:
		return m.viewSizePicker(p)
	case onlineRoomIDInput:
		return m.viewRoomIDPrompt(p)
	case onlinePinInput:
		return m.viewPinPrompt(p)
	case onlineConfigPick:
		return m.viewConfig(p)
	case onlineUsername:
		return m.viewUsernamePrompt(p)
	case onlineConnecting:
		return m.viewConnecting(p)
	case onlineLobby:
		return m.viewLobby(p)
	case onlineCountdown:
		return m.viewRaceCountdown(p)
	case onlineResults:
		return m.viewOnlineResults(p)
	}
	return ""
}

func (m model) viewConfig(p theme.Palette) string {
	dim := lipgloss.NewStyle().Foreground(p.Foreground)
	hi := lipgloss.NewStyle().Foreground(p.Accent)
	val := lipgloss.NewStyle().Foreground(p.Typed)

	opts := m.getOnlineConfigOptions()

	var lines []string
	lines = append(lines, hi.Render("configure room"), "")

	for i, opt := range opts {
		prefix := "  "
		style := dim
		if i == m.onlineConfigCur {
			prefix = hi.Render("● ")
			style = val
		}

		var setting string
		switch opt {
		case "Mode":
			setting = m.mode
		case "Language":
			setting = m.lang
		case "Difficulty":
			setting = m.difficulty
		case "Duration":
			setting = fmt.Sprintf("%ds", m.duration)
		case "Start Policy":
			if m.onlineAutoStart {
				setting = "auto when full"
			} else {
				setting = "host starts"
			}
		case "Continue":
			setting = ""
		}

		if setting != "" {
			lines = append(lines, prefix+opt+": "+style.Render(setting))
		} else {
			lines = append(lines, prefix+opt)
		}
	}

	lines = append(lines, "", dim.Render("↑↓ move · enter toggle/select · esc back"))
	return lipgloss.JoinVertical(lipgloss.Center, lines...)
}

func mInput(p theme.Palette, prompt, value, hint string) string {
	dim := lipgloss.NewStyle().Foreground(p.Foreground)
	hi := lipgloss.NewStyle().Foreground(p.Accent)
	val := lipgloss.NewStyle().Foreground(p.Typed)
	cur := lipgloss.NewStyle().Foreground(p.Background).Background(p.Cursor)

	var inputLine string
	if value == "" {
		inputLine = cur.Render(" ")
	} else {
		inputLine = val.Render(value) + cur.Render(" ")
	}

	lines := []string{
		hi.Render(prompt),
		"",
		dim.Render(hint),
		"",
		inputLine,
		"",
		dim.Render("enter to continue · esc to go back"),
	}
	return lipgloss.JoinVertical(lipgloss.Center, lines...)
}

func (m model) viewActionPicker(p theme.Palette) string {
	return renderList(p, "multiplayer", onlineActions, nil, m.onlineActionCur)
}

func (m model) viewSizePicker(p theme.Palette) string {
	labels := make([]string, len(onlineSizes))
	for i, size := range onlineSizes {
		labels[i] = fmt.Sprintf("%d players", size)
	}
	return renderList(p, "room size", labels, nil, m.onlineSizeCur)
}

func (m model) viewRoomIDPrompt(p theme.Palette) string {
	return mInput(p, "join room", strings.ToUpper(m.onlineRoomIDBuf), "enter room id (e.g. ABCD)")
}

func (m model) viewPinPrompt(p theme.Palette) string {
	label := "room pin"
	if m.isCreating {
		label = "create pin (optional)"
	}
	return mInput(p, label, m.onlinePinBuf, "enter a pin or leave blank for public")
}

func (m model) viewUsernamePrompt(p theme.Palette) string {
	display := m.usernameBuf
	if m.username != "" && display == "" {
		display = m.username
	}
	return mInput(p, "multiplayer", display, "enter a username (2-16 chars)")
}

func (m model) viewConnecting(p theme.Palette) string {
	dim := lipgloss.NewStyle().Foreground(p.Foreground)
	hi := lipgloss.NewStyle().Foreground(p.Accent)

	return lipgloss.JoinVertical(lipgloss.Center,
		hi.Render("multiplayer"),
		"",
		dim.Render("connecting..."),
	)
}

func (m model) viewLobby(p theme.Palette) string {
	dim := lipgloss.NewStyle().Foreground(p.Foreground)
	hi := lipgloss.NewStyle().Foreground(p.Accent)
	val := lipgloss.NewStyle().Foreground(p.Typed)

	roomLabel := "room " + hi.Render(m.onlineRoomID)
	if m.onlinePin != "" {
		roomLabel += dim.Render(" (private)")
	}

	modeInfo := fmt.Sprintf("%s · %s · %ds", m.mode, m.difficulty, m.duration)
	if m.mode == "code" {
		modeInfo = fmt.Sprintf("%s · %s · %s · %ds", m.mode, m.lang, m.difficulty, m.duration)
	}

	lines := []string{
		hi.Render("lobby"),
		roomLabel,
		"",
		dim.Render(modeInfo),
		"",
		dim.Render("waiting for players..."),
		"",
	}

	if m.onlineCount > 0 {
		lines = append(lines, val.Render(fmt.Sprintf("%d", m.onlineCount))+" "+dim.Render("online"))
	}

	if len(m.racePlayers) > 0 {
		lines = append(lines, "")
		for _, pl := range m.racePlayers {
			prefix := "  "
			if pl.IsUser {
				prefix = hi.Render("> ")
			}
			lines = append(lines, prefix+val.Render(pl.Name))
		}
	}

	policy := "autostart: when room is full"
	if !m.onlineAutoStart {
		policy = "autostart: off (host starts)"
	}
	lines = append(lines, "", dim.Render(policy))
	if m.raceHostName != "" {
		lines = append(lines, dim.Render("host: "+m.raceHostName))
	}
	if m.isRaceHost && !m.onlineAutoStart && m.raceCanStart {
		lines = append(lines, "", hi.Render("press s to start race"))
	}

	lines = append(lines, "", dim.Render("esc to leave"))

	return lipgloss.JoinVertical(lipgloss.Center, lines...)
}

func (m model) viewRaceCountdown(p theme.Palette) string {
	hi := lipgloss.NewStyle().Foreground(p.Accent).Bold(true)
	dim := lipgloss.NewStyle().Foreground(p.Foreground)

	return lipgloss.JoinVertical(lipgloss.Center,
		hi.Render("get ready"),
		"",
		dim.Render("race starting..."),
	)
}

func (m model) viewOnlineResults(p theme.Palette) string {
	dim := lipgloss.NewStyle().Foreground(p.Foreground)
	hi := lipgloss.NewStyle().Foreground(p.Accent).Bold(true)

	if len(m.racePlayers) == 0 {
		return lipgloss.JoinVertical(lipgloss.Center,
			"",
			hi.Render("race finished"),
			"",
			dim.Render("waiting for server results..."),
		)
	}

	ordinals := []string{"", "1st", "2nd", "3rd", "4th", "5th", "6th"}

	var rows []string
	for i, pl := range m.racePlayers {
		ord := ""
		rank := i + 1
		if rank < len(ordinals) {
			ord = ordinals[rank]
		} else {
			ord = fmt.Sprintf("%dth", rank)
		}

		wpmStr := fmt.Sprintf("%.0f wpm", pl.WPM)
		var row string
		if pl.IsUser {
			marker := hi.Render(" <")
			row = hi.Render(fmt.Sprintf("  %-4s  %-12s  %s", ord, pl.Name, wpmStr)) + marker
		} else {
			colorIdx := stringToColorIndex(pl.Name)
			playerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(botColorHexes[colorIdx]))
			row = dim.Render(fmt.Sprintf("  %-4s  ", ord)) + playerStyle.Render(fmt.Sprintf("%-12s", pl.Name)) + dim.Render(fmt.Sprintf("  %s  ", wpmStr))
		}
		rows = append(rows, row)
	}

	title := "race results"
	// If the race is still technically "ongoing" (we finished locally), show that
	if m.raceClient != nil && m.raceState != onlineResults {
		title = "race finished (waiting for others)"
	}

	return lipgloss.JoinVertical(lipgloss.Center,
		"",
		hi.Render(title),
		"",
		lipgloss.JoinVertical(lipgloss.Left, rows...),
	)
}

func viewOnlineRaceBar(p theme.Palette, players []game.RacePlayer, barWidth int) string {
	if barWidth < 20 {
		barWidth = 20
	}

	hi := lipgloss.NewStyle().Foreground(p.Accent)
	dim := lipgloss.NewStyle().Foreground(p.Foreground)
	val := lipgloss.NewStyle().Foreground(p.Typed)

	type racer struct {
		name     string
		progress float64
		isUser   bool
	}

	racers := make([]racer, len(players))
	for i, pl := range players {
		racers[i] = racer{name: pl.Name, progress: pl.Progress, isUser: pl.IsUser}
	}

	// sort descending by progress
	for i := 0; i < len(racers); i++ {
		for j := i + 1; j < len(racers); j++ {
			if racers[j].progress > racers[i].progress {
				racers[i], racers[j] = racers[j], racers[i]
			}
		}
	}

	nameWidth := 12
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
			colorIdx := stringToColorIndex(r.name)
			playerStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(botColorHexes[colorIdx]))
			nameStr = playerStyle.Render(fmt.Sprintf("%-*s", nameWidth, r.name))
			barStr = playerStyle.Render(strings.Repeat("━", filled)) + dim.Render(strings.Repeat("─", empty))
			pctStr = dim.Render(pct)
		}

		rows = append(rows, nameStr+" "+barStr+" "+pctStr)
	}

	return strings.Join(rows, "\n")
}
