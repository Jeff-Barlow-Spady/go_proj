package main

import (
	"fmt"
	"os"

	"ubuntu-to-fedora/cmd"

	tea "github.com/charmbracelet/bubbletea"
)

// programCreator allows us to mock program creation in tests
func main() {
	if _, err := runTUI(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

// runTUI runs the Bubble Tea TUI for the app
func runTUI() (tea.Model, error) {
	p := tea.NewProgram(cmd.InitialModel())
	return p.Run()
}
