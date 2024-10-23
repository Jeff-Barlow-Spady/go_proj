package cmd

import (
	"fmt"
	"strings"

	"ubuntu-to-fedora/pkg/converter"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF75B5")).
			Bold(true).
			MarginBottom(1)

	itemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF"))

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF75B5")).
				Bold(true)

	checkedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			MarginTop(1)
)

// Model represents the TUI state
type Model struct {
	choices  []converter.AppScript
	cursor   int
	selected map[int]struct{}
	err      error
	quitting bool
	repoDir  string
}

// InitialModel returns a new model with initial state
func InitialModel() Model {
	repoDir := "./omakub"

	// Clone the repository first
	err := converter.CloneOmakubRepo(repoDir)
	if err != nil {
		return Model{
			err:      err,
			repoDir:  repoDir,
			selected: make(map[int]struct{}), // Initialize selected map even on error
		}
	}

	// Get available apps
	apps, err := converter.GetAvailableApps(repoDir)
	if err != nil {
		return Model{
			err:      err,
			repoDir:  repoDir,
			selected: make(map[int]struct{}), // Initialize selected map even on error
		}
	}

	return Model{
		choices:  apps,
		selected: make(map[int]struct{}),
		repoDir:  repoDir,
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
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = struct{}{}
			}
		case "enter":
			// Only process enter if there are selections
			if len(m.selected) > 0 {
				err := runConversion(m.repoDir)
				if err != nil {
					m.err = err
				}
				return m, tea.Quit
			}
			// If no selections, return the model unchanged
			return m, nil
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
		return errorStyle.Render(fmt.Sprintf("Error: %v\n", m.err))
	}

	s := titleStyle.Render("Which apps do you want to keep?")
	s += "\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = "x"
		}

		item := fmt.Sprintf("%s [%s] %s", cursor, checked, choice.Name)

		if m.cursor == i {
			s += selectedItemStyle.Render(item)
		} else {
			s += itemStyle.Render(item)
		}
		s += "\n"
	}

	help := strings.Join([]string{
		"↑/↓: navigate",
		"space: select/unselect",
		"enter: confirm",
		"q: quit",
	}, " • ")

	s += helpStyle.Render(help)
	return s
}

func runConversion(repoDir string) error {
	err := converter.ReplaceUbuntuWithFedora(repoDir)
	if err != nil {
		return fmt.Errorf("error replacing Ubuntu-specific commands: %v", err)
	}

	fmt.Println("Conversion completed successfully.")
	return nil
}
