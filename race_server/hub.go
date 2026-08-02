package main

import (
	crand "crypto/rand"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

type client struct {
	name   string
	room   string
	ip     string
	events chan []byte
	done   chan struct{}
	closed bool
}

type hub struct {
	mu      sync.Mutex
	rooms   map[string]*room
	clients map[*client]bool
}

func newHub() *hub {
	return &hub{
		rooms:   make(map[string]*room),
		clients: make(map[*client]bool),
	}
}

func (h *hub) onlineCount() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.clients)
}

func (h *hub) removeClient(c *client) {
	h.mu.Lock()
	if c.closed {
		h.mu.Unlock()
		return
	}
	c.closed = true
	delete(h.clients, c)
	close(c.done)
	log.Printf("client %s disconnected from %s", c.name, c.ip)

	r, ok := h.rooms[c.room]
	count := 0
	started := false
	if ok {
		r.removePlayer(c)

		r.mu.Lock()
		count = len(r.players)
		started = r.started
		r.mu.Unlock()

		if count == 0 {
			r.close()
			delete(h.rooms, r.id)
		}
	}
	h.mu.Unlock()

	if !ok || count == 0 {
		return
	}
	if !started {
		r.broadcastLobby()
	} else {
		r.broadcastProgress()
	}
}

func (h *hub) getRoom(id string) *room {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.rooms[id]
}

func generateRoomID() string {
	const chars = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	b := make([]byte, 4)
	if _, err := crand.Read(b); err != nil {
		// crypto/rand unavailable: fall back to math/rand (auto-seeded since Go 1.20).
		for i := range b {
			b[i] = chars[rand.Intn(len(chars))]
		}
		return string(b)
	}
	// 256 is an exact multiple of len(chars), so the modulo is uniform.
	for i := range b {
		b[i] = chars[int(b[i])%len(chars)]
	}
	return string(b)
}

// uniqueRoomID returns a room ID not already in use. Caller must hold h.mu.
func (h *hub) uniqueRoomID() string {
	for {
		id := generateRoomID()
		if _, exists := h.rooms[id]; !exists {
			return id
		}
	}
}

func (h *hub) serveJoin(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "name required", http.StatusBadRequest)
		return
	}

	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}
	if strings.Contains(ip, ",") {
		ip = strings.Split(ip, ",")[0]
	}
	if strings.Contains(ip, ":") {
		lastColon := strings.LastIndex(ip, ":")
		if lastColon > strings.LastIndex(ip, "]") {
			ip = ip[:lastColon]
		}
	}

	roomID := r.URL.Query().Get("room")
	pin := r.URL.Query().Get("pin")
	isCreate := r.URL.Query().Get("is_create") == "true"
	noRooms := r.URL.Query().Get("no_rooms") == "true"

	h.mu.Lock()
	var rm *room
	if isCreate {
		roomID = h.uniqueRoomID()
		size := 2
		fmt.Sscanf(r.URL.Query().Get("size"), "%d", &size)
		dur := 30
		fmt.Sscanf(r.URL.Query().Get("duration"), "%d", &dur)
		autoStart := r.URL.Query().Get("auto_start") != "false"

		rm = newRoom(h, roomID, pin, size,
			r.URL.Query().Get("difficulty"),
			r.URL.Query().Get("mode"),
			r.URL.Query().Get("lang"),
			dur,
			name,
			autoStart)
		h.rooms[roomID] = rm
		log.Printf("create lobby room=%s owner=%s private=%t auto_start=%t", roomID, name, pin != "", autoStart)
	} else if roomID != "" {
		rm = h.rooms[roomID]
		if rm == nil {
			h.mu.Unlock()
			http.Error(w, "room not found", http.StatusNotFound)
			return
		}
		if rm.pin != "" && rm.pin != pin {
			h.mu.Unlock()
			http.Error(w, "invalid pin", http.StatusForbidden)
			return
		}
		if rm.playerCount() >= rm.maxPlayers {
			h.mu.Unlock()
			http.Error(w, "room full", http.StatusForbidden)
			return
		}
	} else {
		// auto-match in public rooms.
		// no_rooms=true means queue for random 1v1 only (never join custom-sized rooms).
		for _, ex := range h.rooms {
			if ex.isStarted() || ex.pin != "" || ex.playerCount() >= ex.maxPlayers {
				continue
			}
			if noRooms && ex.maxPlayers != 2 {
				continue
			}
			if noRooms && (ex.mode != "words" || ex.lang != "english" || ex.difficulty != "medium" || ex.duration != 30) {
				continue
			}
				rm = ex
				break
		}
		if rm == nil {
			roomID = h.uniqueRoomID()
			// default quick queue room
			rm = newRoom(h, roomID, "", 2, "medium", "words", "english", 30, name, true)
			h.rooms[roomID] = rm
			log.Printf("auto-created queue room=%s no_rooms=%t", roomID, noRooms)
		}
	}

	c := &client{
		name:   name,
		ip:     ip,
		room:   rm.id,
		events: make(chan []byte, 64),
		done:   make(chan struct{}),
	}
	h.clients[c] = true
	h.mu.Unlock()

	rm.addPlayer(c)
	log.Printf("client %s joined room %s from %s no_rooms=%t", name, rm.id, ip, noRooms)

	defer h.removeClient(c)

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	rm.broadcastLobby()
	rm.maybeStart()

	ctx := r.Context()
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-c.done:
			return
		case <-ticker.C:
			fmt.Fprintf(w, ": heartbeat\n\n")
			flusher.Flush()
		case data, ok := <-c.events:
			if !ok {
				return
			}
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		}
	}
}
