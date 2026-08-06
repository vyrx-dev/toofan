package main

import (
	"fmt"
	"log"
	"sort"
	"sync"
	"time"
)

type player struct {
	client   *client
	name     string
	progress float64
	wpm      float64
	finished bool
}

type room struct {
	id         string
	pin        string
	hub        *hub
	mu         sync.Mutex
	players    map[*client]*player
	maxPlayers int
	started    bool
	counting   bool
	finished   bool
	startTime  time.Time
	text       string
	closed     bool
	difficulty string
	mode       string
	lang       string
	duration   int
	host       string
	autoStart  bool
}

func newRoom(h *hub, id string, pin string, size int, diff, mode, lang string, dur int, host string, autoStart bool) *room {
	rm := &room{
		id:         id,
		pin:        pin,
		hub:        h,
		players:    make(map[*client]*player),
		maxPlayers: size,
		difficulty: diff,
		mode:       mode,
		lang:       lang,
		duration:   dur,
		host:       host,
		autoStart:  autoStart,
	}
	log.Printf("room created id=%s private=%t size=%d mode=%s lang=%s difficulty=%s duration=%ds host=%s auto_start=%t", id, pin != "", size, mode, lang, diff, dur, host, autoStart)
	return rm
}

func (r *room) addPlayer(c *client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.players[c] = &player{client: c, name: c.name}
	log.Printf("room=%s player_join name=%s players=%d/%d", r.id, c.name, len(r.players), r.maxPlayers)
}

func (r *room) removePlayer(c *client) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.players, c)
	if r.host == c.name {
		r.host = ""
		for _, p := range r.players {
			r.host = p.name
			break
		}
		if r.host != "" {
			log.Printf("room=%s host_reassigned host=%s", r.id, r.host)
		}
	}
	log.Printf("room=%s player_leave name=%s players=%d/%d", r.id, c.name, len(r.players), r.maxPlayers)
}

// playerCount returns the number of players. It locks r.mu so callers never
// read r.players while holding only the hub lock.
func (r *room) playerCount() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return len(r.players)
}

// isStarted reports whether the race has begun. See playerCount.
func (r *room) isStarted() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.started
}

func (r *room) broadcast(msg ServerMsg) {
	data := marshal(msg)
	r.mu.Lock()
	defer r.mu.Unlock()
	for c := range r.players {
		select {
		case c.events <- data:
		default:
		}
	}
}

func (r *room) broadcastLobby() {
	// Grab the online count before taking r.mu: hub.onlineCount locks h.mu, and
	// serveJoin locks h.mu then r.mu. Acquiring h.mu while holding r.mu would
	// invert that order and deadlock.
	online := r.hub.onlineCount()

	r.mu.Lock()
	state := "lobby"
	if r.started {
		state = "racing"
	} else if r.counting {
		state = "countdown"
	}

	timeLeft := 0
	if r.started && r.duration > 0 {
		elapsed := time.Since(r.startTime)
		timeLeft = r.duration - int(elapsed.Seconds())
		if timeLeft < 0 {
			timeLeft = 0
		}
	}

	msg := ServerMsg{
		Type: "joined",
		Payload: JoinMsg{
			Room:       r.id,
			Players:    r.playerNamesLocked(),
			Online:     online,
			Host:       r.host,
			AutoStart:  r.autoStart,
			CanStart:   r.canStartLocked(),
			Difficulty: r.difficulty,
			Mode:       r.mode,
			Lang:       r.lang,
			Duration:   r.duration,
			IsPrivate:  r.pin != "",
			State:      state,
			Text:       r.text,
			TimeLeft:   timeLeft,
		},
	}
	r.mu.Unlock()
	r.broadcast(msg)
}

func (r *room) playerNamesLocked() []string {
	names := make([]string, 0, len(r.players))
	for _, p := range r.players {
		names = append(names, p.name)
	}
	return names
}

func (r *room) maybeStart() {
	r.mu.Lock()
	count := len(r.players)
	alreadyStarted := r.started
	maxP := r.maxPlayers
	autoStart := r.autoStart
	r.mu.Unlock()

	if !autoStart || alreadyStarted || count < maxP {
		return
	}

	go r.startCountdown()
}

func (r *room) canStartLocked() bool {
	if r.started || r.counting {
		return false
	}
	return len(r.players) >= 2
}

func (r *room) requestStart(name string) error {
	r.mu.Lock()
	if name != r.host {
		r.mu.Unlock()
		return fmt.Errorf("only host can start")
	}
	if !r.canStartLocked() {
		r.mu.Unlock()
		return fmt.Errorf("not enough players to start")
	}
	r.mu.Unlock()
	go r.startCountdown()
	return nil
}

func (r *room) startCountdown() {
	r.mu.Lock()
	if r.started || r.counting {
		r.mu.Unlock()
		return
	}
	r.counting = true
	r.mu.Unlock()

	r.broadcast(ServerMsg{
		Type:    "countdown",
		Payload: CountdownMsg{Seconds: 3},
	})

	time.Sleep(3 * time.Second)

	r.mu.Lock()
	r.text = generateText(r.mode, r.lang, r.difficulty)
	r.started = true
	r.counting = false
	r.startTime = time.Now()
	textPreview := r.text
	if len(textPreview) > 80 {
		textPreview = textPreview[:80] + "..."
	}
	log.Printf("room=%s race_start players=%d mode=%s lang=%s difficulty=%s duration=%ds text_len=%d text_preview=%q",
		r.id, len(r.players), r.mode, r.lang, r.difficulty, r.duration, len(r.text), textPreview)
	r.mu.Unlock()

	r.broadcast(ServerMsg{
		Type:    "start",
		Payload: StartMsg{
			Text:       r.text,
			Difficulty: r.difficulty,
			Mode:       r.mode,
			Lang:       r.lang,
			Duration:   r.duration,
		},
	})

	if r.duration > 0 {
		go func() {
			time.Sleep(time.Duration(r.duration+5) * time.Second)
			r.forceFinish()
		}()
	}
}

func (r *room) forceFinish() {
	r.mu.Lock()
	if r.closed {
		r.mu.Unlock()
		return
	}
	
	// Only finish if not already finished by allDone in broadcastProgress.
	// The check-and-set happen under r.mu so the race can't be finished twice.
	if r.finished {
		r.mu.Unlock()
		return
	}
	r.finished = true

	players := make([]PlayerProgress, 0, len(r.players))
	for _, p := range r.players {
		players = append(players, PlayerProgress{
			Name:     p.name,
			Progress: p.progress,
			WPM:      p.wpm,
			Finished: p.finished,
		})
	}
	r.mu.Unlock()

	// sort by progress desc
	sort.Slice(players, func(i, j int) bool {
		return players[i].Progress > players[j].Progress
	})

	r.broadcast(ServerMsg{Type: "finish", Payload: FinishMsg{Placements: players}})
	log.Printf("room=%s race_finish reason=timeout players=%d", r.id, len(players))
}

func (r *room) updateProgress(name string, progress float64, wpm float64) {
	r.mu.Lock()

	for _, p := range r.players {
		if p.name == name {
			p.progress = progress
			p.wpm = wpm
			if progress >= 1.0 {
				p.finished = true
			}
			break
		}
	}
	r.mu.Unlock()
	r.broadcastProgress()
}

func (r *room) broadcastProgress() {
	r.mu.Lock()
	defer r.mu.Unlock()

	players := make([]PlayerProgress, 0, len(r.players))
	allDone := len(r.players) > 0
	for _, p := range r.players {
		players = append(players, PlayerProgress{
			Name:     p.name,
			Progress: p.progress,
			WPM:      p.wpm,
			Finished: p.finished,
		})
		if !p.finished {
			allDone = false
		}
	}

	// sort by progress desc
	sort.Slice(players, func(i, j int) bool {
		return players[i].Progress > players[j].Progress
	})

	data := marshal(ServerMsg{Type: "progress", Payload: ProgressMsg{Players: players}})
	for c := range r.players {
		select {
		case c.events <- data:
		default:
		}
	}

	// Send the finish message exactly once: guard with r.finished under r.mu
	// so forceFinish can't double-finish the same race.
	if allDone && !r.finished {
		r.finished = true
		finishData := marshal(ServerMsg{Type: "finish", Payload: FinishMsg{Placements: players}})
		for c := range r.players {
			select {
			case c.events <- finishData:
			default:
			}
		}
		log.Printf("room=%s race_finish reason=all_done players=%d", r.id, len(players))
	}
}

func (r *room) close() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.closed = true
}
