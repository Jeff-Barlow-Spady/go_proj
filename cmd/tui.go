package cmd

import (
    "fmt"
    "omakub-fedora/ubuntuswap"

    "github.com/charmbracelet/bubbletea"
)

// Improved TUI model with graceful error handling
type model struct {
    choices          []string         // Available packages/apps
    selected         map[int]struct{} // User-selected packages
    err              error            // Capture any errors that occur
    quitting         bool             // Quit flag
    destinationDir   string           // Destination directory for conversion
}

// NewModel creates a new TUI model with default values
func NewModel() model {
    return model{
        choices:        []string{"App1", "App2", "App3"}, // Example choices
        selected:       make(map[int]struct{}),
        destinationDir: "./omakub", // Example destination directory
    }
}

// Init function to satisfy the tea.Model interface
func (m model) Init() tea.Cmd {
    return nil
}

// Update function with graceful error handling
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c", "q":
            m.quitting = true
            return m, tea.Quit
        case "enter":
            // Run the conversion, but handle errors gracefully
            err := runConversion(m.destinationDir)
            if err != nil {
                m.err = err
            }
            return m, tea.Quit
        }
    }
    return m, nil
}

// View function that displays any errors
func (m model) View() string {
    if (m.quitting) {
        return "Bye!\n"
    }

    if (m.err != nil) {
        return fmt.Sprintf("Error: %v\n", m.err)
    }

    s := "Which apps do you want to keep? Press Enter to confirm.\n\n"
    for i, choice := range m.choices {
        selected := "[ ] "
        if _, ok := m.selected[i]; ok {
            selected = "[x] "
        }
        s += fmt.Sprintf("%s%s\n", selected, choice)
    }

    return s
}

func runConversion(destinationDir string) error {
    err := ubuntuswap.CloneOmakubRepo(destinationDir)
    if err != nil {
        return fmt.Errorf("error cloning repository: %v", err)
    }

    err = ubuntuswap.ReplaceUbuntuWithFedora(destinationDir)
    if err != nil {
        return fmt.Errorf("error replacing Ubuntu-specific commands: %v", err)
    }

    fmt.Println("Conversion completed successfully.")
    return nil
}
