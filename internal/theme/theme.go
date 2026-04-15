package theme

import "github.com/charmbracelet/lipgloss"

// Palette holds all the colors for a theme
type Palette struct {
	Name       string
	Background lipgloss.Color
	Foreground lipgloss.Color // untyped text, hints
	Typed      lipgloss.Color // correctly typed
	Error      lipgloss.Color // mistakes
	Cursor     lipgloss.Color // current character
	Accent     lipgloss.Color // highlights, timer, active elements
	Success    lipgloss.Color // personal best, positive feedback
}

var All = []Palette{TokyoNight, Gruvbox, Sakura, Monkeytype, Monochrome, Forest, Espresso, Lumon, Mars, Void, Everforest, Omarchy}

var Current = TokyoNight

// Next cycles themes
func Next() {
	for i, t := range All {
		if t.Name == Current.Name {
			Current = All[(i+1)%len(All)]
			return
		}
	}
}

// ByName finds a theme, defaults to TokyoNight
func ByName(name string) Palette {
	for _, t := range All {
		if t.Name == name {
			return t
		}
	}
	return TokyoNight
}
