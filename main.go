package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vyrx-dev/toofan/internal/tui"
)

var version = "dev"

func main() {
	if len(os.Args) == 2 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Println("toofan " + version)
		return
	}

	p := tea.NewProgram(tui.New(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
