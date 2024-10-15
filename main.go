package main

import (
	"fmt"
	"omakub-fedora/cmd"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Initialize the TUI model
	m := cmd.NewModel()

	// Start the TUI
	p := tea.NewProgram(m)
	if err := p.Start(); err != nil {
		fmt.Printf("Error starting TUI: %v\n", err)
		os.Exit(1)
	}
}
