package cmd

import (
	"fmt"
	"os"
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

	cursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF75B5"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			MarginTop(1)

	containerStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			Padding(1, 2)

	welcomeStyle = lipgloss.NewStyle().
			Align(lipgloss.Center).
			Padding(1, 2)
)

// Model represents the TUI state
type Model struct {
	choices     []converter.AppScript
	cursor      int
	selected    map[int]struct{}
	err         error
	quitting    bool
	repoDir     string
	help        string
	width       int
	height      int
	windowStart int
	windowSize  int
	showWelcome bool
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(tea.EnterAltScreen, tea.ClearScreen)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.showWelcome {
		switch msg.(type) {
		case tea.KeyMsg:
			m.showWelcome = false
		}
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		headerHeight := 5 // Adjusted for header size
		footerHeight := 4 // Adjusted for footer size
		m.windowSize = m.height - headerHeight - footerHeight
		if m.windowSize > len(m.choices) {
			m.windowSize = len(m.choices)
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				if m.cursor < m.windowStart {
					m.windowStart--
				}
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
				if m.cursor >= m.windowStart+m.windowSize {
					m.windowStart++
				}
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

	m.help = "Press 'q' to quit, 'space' to select, 'enter' to confirm, 'up' and 'down' to navigate"
	return m, nil
}

// InitialModel returns a new model with initial state
func InitialModel() Model {
	repoDir := "./omakub"

	// Clone the repository first
	err := converter.CloneOmakubRepo(repoDir)
	if err != nil {
		return Model{
			err:         err,
			repoDir:     repoDir,
			selected:    make(map[int]struct{}), // Initialize selected map even on error
			showWelcome: true,
			windowSize:  10, // Default window size
		}
	}

	// Get available apps
	apps, err := converter.GetAvailableApps(repoDir)
	if err != nil {
		return Model{
			err:         err,
			repoDir:     repoDir,
			selected:    make(map[int]struct{}), // Initialize selected map even on error
			showWelcome: true,
			windowSize:  10, // Default window size
		}
	}

	return Model{
		choices:     apps,
		selected:    make(map[int]struct{}),
		repoDir:     repoDir,
		showWelcome: true,
		windowSize:  10, // Default window size
	}
}

// View renders the current state of the model
func (m Model) View() string {
	if m.quitting {
		return "Thanks for using. Bye!\n"
	}

	if m.err != nil {
		return errorStyle.Render(fmt.Sprintf("Error: %v\n", m.err))
	}

	if m.showWelcome {
		welcome := `
 ___  ___   ________   ___  ___   ________                              _________   ________                             ________  _______    ________   ________   ________   ________
|\  \|\  \ |\   __  \ |\  \|\  \ |\   ___  \                           |\___   ___\|\   __  \                           |\  _____\|\  ___ \  |\   ___ \ |\   __  \ |\   __  \ |\   __  \
\ \  \\\  \\ \  \|\ /_\ \  \\\  \\ \  \\ \  \        ____________      \|___ \  \_|\ \  \|\  \        ____________      \ \  \__/ \ \   __/| \ \  \_|\ \\ \  \|\  \\ \  \|\  \\ \  \|\  \
 \ \  \\\  \\ \   __  \\ \  \\\  \\ \  \\ \  \      |\____________\         \ \  \  \ \  \\\  \      |\____________\     \ \   __\ \ \  \_|/__\ \  \ \\ \\ \  \\\  \\ \   _  _\\ \   __  \
  \ \  \\\  \\ \  \|\  \\ \  \\\  \\ \  \\ \  \     \|____________|          \ \  \  \ \  \\\  \     \|____________|      \ \  \_|  \ \  \_|\ \\ \  \_\\ \\ \  \\\  \\ \  \\  \|\ \  \ \  \
   \ \_______\\ \_______\\ \_______\\ \__\\ \__\                              \ \__\  \ \_______\                          \ \__\    \ \_______\\ \_______\\ \_______\\ \__\\ _\ \ \__\ \__\
    \|_______| \|_______| \|_______| \|__| \|__|                               \|__|   \|_______|                           \|__|     \|_______| \|_______| \|_______| \|__|\|__| \|__|\|__|

Welcome to the Ubuntu to Fedora Converter!
Press any key to continue...
		`
		return welcomeStyle.Render(welcome)
	}

	s := titleStyle.Render("Select which apps you want to keep. It is a large list of shell scripts in the omakub directory. Work in progress.")
	s += "\n\n"

	// Handle window size not yet set
	if m.windowSize == 0 {
		m.windowSize = len(m.choices)
	}

	// Only display items in the current window
	end := m.windowStart + m.windowSize
	if end > len(m.choices) {
		end = len(m.choices)
	}

	for i := m.windowStart; i < end; i++ {
		choice := m.choices[i]
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

	additionalHelp := "Select the applications you wish to keep. Unselected applications will be converted to Fedora equivalents."
	s += "\n" + helpStyle.Render(additionalHelp)

	help := strings.Join([]string{
		"↑/↓: navigate",
		"space: select/unselect",
		"enter: confirm",
		"q: quit",
	}, " • ")

	s += "\n" + helpStyle.Render(help)
	return containerStyle.Render(s)
}

func runConversion(repoDir string) error {
	err := converter.ReplaceUbuntuWithFedora(repoDir)
	if err != nil {
		return fmt.Errorf("error replacing Ubuntu-specific commands: %v", err)
	}

	fmt.Println("Conversion completed successfully.")
	return nil
}

func main() {
	m := InitialModel()
	p := tea.NewProgram(m, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
