package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vyrx-dev/toofan/internal/tui"
)

var version = "dev"

func main() {
	v := flag.Bool("version", false, "Print current version")
	flag.Parse()

	if *v {
		fmt.Println("toofan " + version)
		os.Exit(0)
	}

	p := tea.NewProgram(tui.New(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
