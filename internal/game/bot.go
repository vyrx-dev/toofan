package game

import (
	"bufio"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var botNames = []string{
	"turbo", "ghost", "nitro", "flash",
	"zen", "byte", "pixel", "glitch",
}

type Bot struct {
	ID       int
	Name     string
	WPM      float64
	Accuracy float64
	Progress float64
	Finished bool
}

type BotConfig struct {
	Count    int
	Strategy string // "personalized", "easy", "medium", "hard"
}

func NewBots(cfg BotConfig, userAvgWPM float64) []Bot {
	if cfg.Count <= 0 {
		return nil
	}
	if cfg.Count > 5 {
		cfg.Count = 5
	}

	names := make([]string, len(botNames))
	copy(names, botNames)
	rand.Shuffle(len(names), func(i, j int) { names[i], names[j] = names[j], names[i] })

	bots := make([]Bot, cfg.Count)
	for i := range bots {
		wpm, acc := botStats(cfg.Strategy, i, cfg.Count, userAvgWPM)
		bots[i] = Bot{
			ID:       i,
			Name:     names[i%len(names)],
			WPM:      wpm,
			Accuracy: acc,
		}
	}
	return bots
}

func botStats(strategy string, index, total int, userAvg float64) (wpm, accuracy float64) {
	switch strategy {
	case "easy":
		wpm = 20 + rand.Float64()*20
		accuracy = 0.90 + rand.Float64()*0.05
	case "hard":
		wpm = 80 + rand.Float64()*60
		accuracy = 0.95 + rand.Float64()*0.04
	case "personalized":
		if userAvg < 10 {
			userAvg = 50
		}
		spread := userAvg * 0.4
		lo := userAvg - spread/2
		if lo < 15 {
			lo = 15
		}
		step := spread / float64(total)
		wpm = lo + step*float64(index) + (rand.Float64()-0.5)*10
		if wpm < 15 {
			wpm = 15
		}
		accuracy = 0.92 + rand.Float64()*0.06
	default: // medium
		wpm = 40 + rand.Float64()*40
		accuracy = 0.92 + rand.Float64()*0.05
	}
	return
}

func TickBots(bots []Bot, dt time.Duration, textLen int) {
	if textLen <= 0 {
		return
	}
	secs := dt.Seconds()
	for i := range bots {
		if bots[i].Finished {
			continue
		}

		// chars per second = (WPM * 5) / 60
		cps := (bots[i].WPM * 5.0) / 60.0

		// add jitter: ±15% variation per tick
		jitter := 1.0 + (rand.Float64()-0.5)*0.3
		chars := cps * secs * jitter

		bots[i].Progress += chars / float64(textLen)
		if bots[i].Progress >= 1.0 {
			bots[i].Progress = 1.0
			bots[i].Finished = true
		}
	}
}

func BotPlacements(bots []Bot, userProgress float64) []BotPlacement {
	type entry struct {
		id       int
		name     string
		progress float64
		isUser   bool
	}

	entries := make([]entry, 0, len(bots)+1)
	entries = append(entries, entry{id: -1, name: "you", progress: userProgress, isUser: true})
	for _, b := range bots {
		entries = append(entries, entry{id: b.ID, name: b.Name, progress: b.Progress})
	}

	// sort descending by progress
	for i := 0; i < len(entries); i++ {
		for j := i + 1; j < len(entries); j++ {
			if entries[j].progress > entries[i].progress {
				entries[i], entries[j] = entries[j], entries[i]
			}
		}
	}

	placements := make([]BotPlacement, len(entries))
	for i, e := range entries {
		wpm := 0.0
		if e.isUser {
			// we don't have user wpm here, but it's usually passed or calculated
		} else {
			for _, b := range bots {
				if b.ID == e.id {
					wpm = b.WPM
					break
				}
			}
		}
		placements[i] = BotPlacement{
			Rank:     i + 1,
			ID:       e.id,
			Name:     e.name,
			Progress: e.progress,
			WPM:      wpm,
			IsUser:   e.isUser,
		}
	}
	return placements
}

type BotPlacement struct {
	Rank     int
	ID       int
	Name     string
	Progress float64
	WPM      float64
	IsUser   bool
}

func GetUserAvgWPM() float64 {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return 0
	}
	path := filepath.Join(configDir, "toofan", "results.txt")
	f, err := os.Open(path)
	if err != nil {
		return 0
	}
	defer f.Close()

	var wpms []float64
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		parts := strings.Split(sc.Text(), "|")
		if len(parts) < 2 {
			continue
		}
		wpmStr := strings.TrimSpace(parts[1])
		wpmStr = strings.TrimSuffix(wpmStr, "wpm")
		wpmStr = strings.TrimSpace(wpmStr)
		if w, err := strconv.ParseFloat(wpmStr, 64); err == nil {
			wpms = append(wpms, w)
		}
	}

	if len(wpms) == 0 {
		return 0
	}

	n := 10
	if len(wpms) < n {
		n = len(wpms)
	}

	sum := 0.0
	for i := len(wpms) - n; i < len(wpms); i++ {
		sum += wpms[i]
	}
	return sum / float64(n)
}
