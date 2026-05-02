package game

import (
	"strings"
	"time"

	"github.com/vyrx-dev/toofan/internal/lang"
)

type Stats struct {
	WPM      float64
	Raw      float64
	Accuracy float64
	Chars    int
	Mistakes int // total wrong keystrokes, including corrected ones
}

type Game struct {
	text      string
	Snippet   lang.Snippet
	input     string
	errors    map[int]bool // unfixed errors — used for live display coloring
	mistakeAt map[int]bool // every wrong keystroke ever — never cleared by backspace
	started   bool
	duration  int
	CodeMode  bool // true = snippet-based typing, false = standard words
	difficulty string
	GhostPos   int

	elapsed   time.Duration
	lastTyped time.Time
	LastTick  time.Time
}

// Accessors for fields that TUI needs to read
func (g *Game) Text() string         { return g.text }
func (g *Game) Input() string        { return g.input }
func (g *Game) Errors() map[int]bool { return g.errors }
func (g *Game) Started() bool        { return g.started }
func (g *Game) Duration() int        { return g.duration }
func (g *Game) Elapsed() time.Duration { return g.elapsed }
func (g *Game) GhostPosition() int   { return g.GhostPos }
func (g *Game) SetText(s string)     { g.text = normalizeTabs(s) }

func (g *Game) UpdateGhost(pbWPM float64) {
	if !g.started || pbWPM <= 0 {
		return
	}
	elapsed := g.elapsed.Seconds()
	g.GhostPos = int((pbWPM * 5.0 * elapsed) / 60.0)
}

func New(duration int, mode string, language string, difficulty string) *Game {
	g := &Game{
		duration:  duration, // 0 means infinite mode (tied to length of snippet)
		errors:    make(map[int]bool),
		mistakeAt: make(map[int]bool),
		difficulty: difficulty,
	}

	if mode == "code" {
		g.CodeMode = true
		g.Snippet = lang.RandomSnippet(language, difficulty)
		g.text = normalizeTabs(g.Snippet.Content)
	} else {
		words := lang.RandomWords(language, difficulty, 200)
		g.text = strings.Join(words, " ")
	}

	return g
}

func (g *Game) Tick(now time.Time) {
	if !g.started || g.Finished() {
		return
	}

	dt := now.Sub(g.LastTick)
	g.LastTick = now
	g.elapsed += dt
}

// normalizeTabs replaces tabs with spaces so position tracking and
// rendering stay perfectly aligned in the terminal.
func normalizeTabs(s string) string {
	return strings.ReplaceAll(s, "\t", "    ")
}

func isStartOfLine(s string) bool {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '\n' {
			return true
		}
		if s[i] != ' ' {
			return false
		}
	}
	return true
}

func (g *Game) TypeChar(ch rune) {
	if g.Finished() {
		return
	}

	if !g.started {
		g.started = true
		g.LastTick = time.Now()
	}

	g.lastTyped = time.Now()

	pos := len(g.input)
	if pos >= len(g.text) {
		return
	}

	if g.CodeMode && ch == '\n' && g.text[pos] != '\n' {
		return
	}

	g.input += string(ch)
	if g.text[pos] != byte(ch) {
		g.errors[pos] = true
		g.mistakeAt[pos] = true
	}

	// auto-skip newlines and leading indentation
	for len(g.input) < len(g.text) {
		next := g.text[len(g.input)]
		if next == '\n' {
			if g.CodeMode {
				break
			}
			g.input += "\n"
		} else if next == ' ' && isStartOfLine(g.input) {
			g.input += " "
		} else {
			break
		}
	}
}


func (g *Game) Backspace() {
	if len(g.input) == 0 {
		return
	}

	g.lastTyped = time.Now()

	// skip back over auto-inserted newlines and leading whitespace
	for len(g.input) > 0 {
		last := g.input[len(g.input)-1]
		if last == '\n' {
			g.input = g.input[:len(g.input)-1]
		} else if last == ' ' && isStartOfLine(g.input[:len(g.input)-1]) {
			g.input = g.input[:len(g.input)-1]
		} else {
			break
		}
	}

	if len(g.input) > 0 {
		pos := len(g.input) - 1
		g.input = g.input[:pos]
		delete(g.errors, pos)
		// mistakeAt is NOT cleared — tracks lifetime errors
	}
}

func (g *Game) TimeLeft() int {
	if !g.started {
		return g.duration
	}
	if g.duration == 0 {
		return int(g.elapsed.Seconds()) // Count up in infinite mode
	}
	left := g.duration - int(g.elapsed.Seconds())
	if left < 0 {
		return 0
	}
	return left
}

func (g *Game) Finished() bool {
	if !g.started {
		return false
	}
	if g.duration == 0 {
		// Infinite duration finishes when the text is completed
		return len(g.input) >= len(g.text)
	}
	return g.TimeLeft() == 0 || len(g.input) >= len(g.text)
}

func (g *Game) Stats() Stats {
	if !g.started {
		return Stats{}
	}

	mins := g.elapsed.Minutes()
	if mins < 0.001 {
		mins = 0.001
	}

	// count only real typed chars (skip auto-inserted newlines)
	total := 0
	unfixed := 0
	mistakes := 0
	for i := 0; i < len(g.input); i++ {
		if g.input[i] != '\n' {
			total++
			if g.errors[i] {
				unfixed++
			}
			if g.mistakeAt[i] {
				mistakes++
			}
		}
	}

	// WPM uses unfixed errors (corrections already penalized via time spent)
	raw := float64(total) / 5.0 / mins
	wpm := float64(total-unfixed) / 5.0 / mins
	if wpm < 0 {
		wpm = 0
	}

	// Accuracy uses ALL mistakes — correcting doesn't erase the miss
	acc := 100.0
	if total > 0 {
		acc = float64(total-mistakes) / float64(total) * 100
		if acc < 0 {
			acc = 0
		}
	}

	return Stats{WPM: wpm, Raw: raw, Accuracy: acc, Chars: total, Mistakes: mistakes}
}

// ErrorWords returns the words where the user made mistakes (word mode only).
// Includes corrected mistakes so users can see where they struggled.
func (g *Game) ErrorWords() []string {
	if g.CodeMode {
		return nil
	}
	words := strings.Split(g.text, " ")
	pos := 0
	var result []string
	seen := make(map[string]bool)
	for _, word := range words {
		if pos >= len(g.input) {
			break
		}
		hasErr := false
		for j := 0; j < len(word) && pos+j < len(g.input); j++ {
			if g.mistakeAt[pos+j] {
				hasErr = true
				break
			}
		}
		if hasErr && !seen[word] {
			result = append(result, word)
			seen[word] = true
		}
		pos += len(word) + 1 // +1 for space separator
	}
	return result
}

func (g *Game) Reset(mode string, language string, difficulty string) {
	g.difficulty = difficulty
	if mode == "code" {
		g.CodeMode = true
		g.Snippet = lang.RandomSnippet(language, difficulty)
		g.text = normalizeTabs(g.Snippet.Content)
	} else {
		g.CodeMode = false
		words := lang.RandomWords(language, difficulty, 200)
		g.text = strings.Join(words, " ")
	}
	g.input = ""
	g.errors = make(map[int]bool)
	g.mistakeAt = make(map[int]bool)
	g.started = false
	g.elapsed = 0
	g.lastTyped = time.Time{}
	g.LastTick = time.Time{}
}
