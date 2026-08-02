package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	port := flag.Int("port", 8525, "port to listen on")
	flag.Parse()

	h := newHub()

	http.HandleFunc("/race/join", h.serveJoin)

	http.HandleFunc("/race/progress", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POST only", http.StatusMethodNotAllowed)
			return
		}

		var update ProgressUpdate
		if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
			log.Printf("progress: bad json from %s: %v", r.RemoteAddr, err)
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}

		rm := h.getRoom(update.Room)
		if rm == nil {
			http.Error(w, "room not found", http.StatusNotFound)
			return
		}

		log.Printf("progress room=%s name=%s progress=%.3f wpm=%.1f", update.Room, update.Name, update.Progress, update.WPM)
		rm.updateProgress(update.Name, update.Progress, update.WPM)
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/race/start", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POST only", http.StatusMethodNotAllowed)
			return
		}
		var req StartRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("start: bad json from %s: %v", r.RemoteAddr, err)
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}
		rm := h.getRoom(req.Room)
		if rm == nil {
			http.Error(w, "room not found", http.StatusNotFound)
			return
		}
		if err := rm.requestStart(req.Name); err != nil {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		log.Printf("manual_start room=%s host=%s", req.Room, req.Name)
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/race/online", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		json.NewEncoder(w).Encode(OnlineMsg{Count: h.onlineCount()})
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "ok")
	})

	addr := fmt.Sprintf(":%d", *port)
	log.Printf("toofan race server listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
