package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/vyrx-dev/toofan/internal/game"
	"github.com/vyrx-dev/toofan/internal/theme"
)

// colorText renders visible lines with typed/error/cursor/untyped styling.
// ghostPositions maps character positions to their specific lipgloss style (for colored bot ghosts).
func colorText(g *game.Game, p theme.Palette, lines []string, top, bot int, ghostPositions map[int]lipgloss.Style) string {
	ok := lipgloss.NewStyle().Foreground(p.Typed)
	bad := lipgloss.NewStyle().Foreground(p.Error).Underline(true)
	cur := lipgloss.NewStyle().Foreground(p.Background).Background(p.Cursor)
	dim := lipgloss.NewStyle().Foreground(p.Foreground)

	typed := len(g.Input())

	// character offset at top of visible window
	startPos := 0
	for i := 0; i < top; i++ {
		startPos += len(lines[i]) + 1 // +1 for the \n or space separator
	}

	var out strings.Builder
	pos := startPos

	for i := top; i < bot && i < len(lines); i++ {
		if i > top {
			out.WriteByte('\n')
			pos++ // skip the \n in position tracking
		}
		for _, ch := range lines[i] {
			s := string(ch)
			ghostStyle, isGhost := ghostPositions[pos]
			switch {
			case pos < typed && g.Errors()[pos]:
				out.WriteString(bad.Render(s))
			case pos < typed:
				out.WriteString(ok.Render(s))
			case pos == typed:
				out.WriteString(cur.Render(s))
			case isGhost:
				out.WriteString(ghostStyle.Render(s))
			default:
				out.WriteString(dim.Render(s))
			}
			pos++
		}
	}
	return out.String()
}

// splitLines splits text for display.
// For code mode, splits on actual newlines (preserving indentation).
// For word mode, wraps at word boundaries.
func splitLines(text string, width int, codeMode bool) []string {
	if codeMode {
		return strings.Split(text, "\n")
	}
	return wrapText(text, width)
}

// wrapText breaks text into lines at word boundaries
func wrapText(text string, w int) []string {
	var lines []string
	var line strings.Builder
	lineLen := 0

	for _, word := range strings.Split(text, " ") {
		wordLen := len(word)
		space := 0
		if lineLen > 0 {
			space = 1
		}

		if lineLen+space+wordLen > w && lineLen > 0 {
			lines = append(lines, line.String())
			line.Reset()
			lineLen = 0
			space = 0
		}

		if space > 0 {
			line.WriteByte(' ')
			lineLen++
		}
		line.WriteString(word)
		lineLen += wordLen
	}

	if line.Len() > 0 {
		lines = append(lines, line.String())
	}
	return lines
}

// cursorLine returns which line the cursor is on
func cursorLine(lines []string, inputLen int) int {
	pos := 0
	for i, line := range lines {
		end := pos + len(line)
		if inputLen <= end {
			return i
		}
		pos = end + 1
	}
	return max(0, len(lines)-1)
}

// col renders text in a fixed-width column.
func col(w int, s string) string {
	return lipgloss.NewStyle().Width(w).Render(s)
}
