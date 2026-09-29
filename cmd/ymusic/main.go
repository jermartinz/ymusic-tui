package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"

	"github.com/jermartinz/ymusic-tui/internal/tui"
)

func main() {
	program := tea.NewProgram(tui.NewModel())
	if _, err := program.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "ymusic: %v\n", err)
		os.Exit(1)
	}
}
