package main

import (
    "fmt"
    "os"

    "github.com/charmbracelet/bubbletea"
    "ubuntu-to-fedora/cmd" // Importing the TUI logic
)

// Entry point
func main() {
    if _, err := runTUI(); err != nil {
        fmt.Printf("Error: %v\n", err)
        os.Exit(1)
    }
}

// runTUI runs the Bubble Tea TUI for the app
func runTUI() (tea.Model, error) {
    // Initialize the TUI
    p := tea.NewProgram(cmd.InitialModel())

    // Start the TUI
    return p.Start()
}
