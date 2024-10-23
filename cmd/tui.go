package cmd

import (
	"fmt"

	"ubuntu-to-fedora/pkg/converter"

	tea "github.com/charmbracelet/bubbletea"
)

// Model represents the TUI state
type Model struct {
	choices  []string         // Available packages/apps
	selected map[int]struct{} // User-selected packages
	err      error            // Capture any errors that occur
	quitting bool             // Quit flag
}

// InitialModel returns a new model with initial state
func InitialModel() Model {
	return Model{
		choices: []string{
			"Chrome",
			"VSCode",
			"Docker",
			"Spotify",
			"Signal",
		},
		selected: make(map[int]struct{}),
	}
}

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles user input and updates the model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "up", "k":
			return m, nil
		case "down", "j":
			return m, nil
		case " ":
			// Toggle selection
			return m, nil
		case "enter":
			// Run the conversion
			err := runConversion()
			if err != nil {
				m.err = err
			}
			return m, tea.Quit
		}
	}
	return m, nil
}

// View renders the current state of the model
func (m Model) View() string {
	if m.quitting {
		return "Bye!\n"
	}

	if m.err != nil {
		return fmt.Sprintf("Error: %v\n", m.err)
	}

	s := "Which apps do you want to keep? Press Enter to confirm.\n\n"
	for i, choice := range m.choices {
		cursor := " "
		if _, ok := m.selected[i]; ok {
			cursor = "x"
		}
		s += fmt.Sprintf("[%s] %s\n", cursor, choice)
	}

	s += "\nPress space to select/unselect\n"
	s += "Press enter to confirm\n"
	s += "Press q to quit\n"

	return s
}

func runConversion() error {
	repoDir := "./omakub"
	err := converter.CloneOmakubRepo(repoDir)
	if err != nil {
		return fmt.Errorf("error cloning repository: %v", err)
	}

	err = converter.ReplaceUbuntuWithFedora(repoDir)
	if err != nil {
		return fmt.Errorf("error replacing Ubuntu-specific commands: %v", err)
	}

	fmt.Println("Conversion completed successfully.")
	return nil
}
