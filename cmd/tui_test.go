package cmd

import (
	"fmt"
	"testing"
	"ubuntu-to-fedora/pkg/converter"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestInitialModel(t *testing.T) {
	model := InitialModel()

	// Since InitialModel now depends on git clone and file system,
	// we'll just test that it initializes without error
	if model.selected == nil {
		t.Error("Expected selected map to be initialized")
	}

	if model.repoDir == "" {
		t.Error("Expected repoDir to be set")
	}
}

func TestModelUpdate(t *testing.T) {
	tests := []struct {
		name          string
		msg           tea.Msg
		initialModel  Model
		expectedModel Model
		expectQuit    bool
	}{
		{
			name: "Quit with q",
			msg:  tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}},
			initialModel: Model{
				choices: []converter.AppScript{
					{Name: "Chrome", FilePath: "chrome.sh"},
					{Name: "VSCode", FilePath: "vscode.sh"},
					{Name: "Docker", FilePath: "docker.sh"},
				},
				selected: make(map[int]struct{}),
				repoDir:  "./test-repo",
			},
			expectedModel: Model{
				choices: []converter.AppScript{
					{Name: "Chrome", FilePath: "chrome.sh"},
					{Name: "VSCode", FilePath: "vscode.sh"},
					{Name: "Docker", FilePath: "docker.sh"},
				},
				selected: make(map[int]struct{}),
				quitting: true,
				repoDir:  "./test-repo",
			},
			expectQuit: true,
		},
		{
			name: "Quit with ctrl+c",
			msg:  tea.KeyMsg{Type: tea.KeyCtrlC},
			initialModel: Model{
				choices: []converter.AppScript{
					{Name: "Chrome", FilePath: "chrome.sh"},
					{Name: "VSCode", FilePath: "vscode.sh"},
					{Name: "Docker", FilePath: "docker.sh"},
				},
				selected: make(map[int]struct{}),
				repoDir:  "./test-repo",
			},
			expectedModel: Model{
				choices: []converter.AppScript{
					{Name: "Chrome", FilePath: "chrome.sh"},
					{Name: "VSCode", FilePath: "vscode.sh"},
					{Name: "Docker", FilePath: "docker.sh"},
				},
				selected: make(map[int]struct{}),
				quitting: true,
				repoDir:  "./test-repo",
			},
			expectQuit: true,
		},
		{
			name: "Move cursor down",
			msg:  tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}},
			initialModel: Model{
				choices: []converter.AppScript{
					{Name: "Chrome", FilePath: "chrome.sh"},
					{Name: "VSCode", FilePath: "vscode.sh"},
					{Name: "Docker", FilePath: "docker.sh"},
				},
				selected: make(map[int]struct{}),
				repoDir:  "./test-repo",
			},
			expectedModel: Model{
				choices: []converter.AppScript{
					{Name: "Chrome", FilePath: "chrome.sh"},
					{Name: "VSCode", FilePath: "vscode.sh"},
					{Name: "Docker", FilePath: "docker.sh"},
				},
				cursor:   1,
				selected: make(map[int]struct{}),
				repoDir:  "./test-repo",
			},
		},
		{
			name: "Move cursor up",
			msg:  tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}},
			initialModel: Model{
				choices: []converter.AppScript{
					{Name: "Chrome", FilePath: "chrome.sh"},
					{Name: "VSCode", FilePath: "vscode.sh"},
					{Name: "Docker", FilePath: "docker.sh"},
				},
				cursor:   1,
				selected: make(map[int]struct{}),
				repoDir:  "./test-repo",
			},
			expectedModel: Model{
				choices: []converter.AppScript{
					{Name: "Chrome", FilePath: "chrome.sh"},
					{Name: "VSCode", FilePath: "vscode.sh"},
					{Name: "Docker", FilePath: "docker.sh"},
				},
				cursor:   0,
				selected: make(map[int]struct{}),
				repoDir:  "./test-repo",
			},
		},
		{
			name: "Select item with space",
			msg:  tea.KeyMsg{Type: tea.KeySpace},
			initialModel: Model{
				choices: []converter.AppScript{
					{Name: "Chrome", FilePath: "chrome.sh"},
					{Name: "VSCode", FilePath: "vscode.sh"},
					{Name: "Docker", FilePath: "docker.sh"},
				},
				selected: make(map[int]struct{}),
				repoDir:  "./test-repo",
			},
			expectedModel: Model{
				choices: []converter.AppScript{
					{Name: "Chrome", FilePath: "chrome.sh"},
					{Name: "VSCode", FilePath: "vscode.sh"},
					{Name: "Docker", FilePath: "docker.sh"},
				},
				selected: map[int]struct{}{
					0: {},
				},
				repoDir: "./test-repo",
			},
		},
		{
			name: "Deselect item with space",
			msg:  tea.KeyMsg{Type: tea.KeySpace},
			initialModel: Model{
				choices: []converter.AppScript{
					{Name: "Chrome", FilePath: "chrome.sh"},
					{Name: "VSCode", FilePath: "vscode.sh"},
					{Name: "Docker", FilePath: "docker.sh"},
				},
				selected: map[int]struct{}{
					0: {},
				},
				repoDir: "./test-repo",
			},
			expectedModel: Model{
				choices: []converter.AppScript{
					{Name: "Chrome", FilePath: "chrome.sh"},
					{Name: "VSCode", FilePath: "vscode.sh"},
					{Name: "Docker", FilePath: "docker.sh"},
				},
				selected: make(map[int]struct{}),
				repoDir:  "./test-repo",
			},
		},
		{
			name: "Cursor at bottom boundary",
			msg:  tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}},
			initialModel: Model{
				choices: []converter.AppScript{
					{Name: "Chrome", FilePath: "chrome.sh"},
					{Name: "VSCode", FilePath: "vscode.sh"},
					{Name: "Docker", FilePath: "docker.sh"},
				},
				cursor:   2,
				selected: make(map[int]struct{}),
				repoDir:  "./test-repo",
			},
			expectedModel: Model{
				choices: []converter.AppScript{
					{Name: "Chrome", FilePath: "chrome.sh"},
					{Name: "VSCode", FilePath: "vscode.sh"},
					{Name: "Docker", FilePath: "docker.sh"},
				},
				cursor:   2, // Should not move past last item
				selected: make(map[int]struct{}),
				repoDir:  "./test-repo",
			},
		},
		{
			name: "Cursor at top boundary",
			msg:  tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}},
			initialModel: Model{
				choices: []converter.AppScript{
					{Name: "Chrome", FilePath: "chrome.sh"},
					{Name: "VSCode", FilePath: "vscode.sh"},
					{Name: "Docker", FilePath: "docker.sh"},
				},
				cursor:   0,
				selected: make(map[int]struct{}),
				repoDir:  "./test-repo",
			},
			expectedModel: Model{
				choices: []converter.AppScript{
					{Name: "Chrome", FilePath: "chrome.sh"},
					{Name: "VSCode", FilePath: "vscode.sh"},
					{Name: "Docker", FilePath: "docker.sh"},
				},
				cursor:   0, // Should not move past first item
				selected: make(map[int]struct{}),
				repoDir:  "./test-repo",
			},
		},
		{
			name: "Enter with no selection",
			msg:  tea.KeyMsg{Type: tea.KeyEnter},
			initialModel: Model{
				choices: []converter.AppScript{
					{Name: "Chrome", FilePath: "chrome.sh"},
					{Name: "VSCode", FilePath: "vscode.sh"},
					{Name: "Docker", FilePath: "docker.sh"},
				},
				selected: make(map[int]struct{}),
				repoDir:  "./test-repo",
			},
			expectedModel: Model{
				choices: []converter.AppScript{
					{Name: "Chrome", FilePath: "chrome.sh"},
					{Name: "VSCode", FilePath: "vscode.sh"},
					{Name: "Docker", FilePath: "docker.sh"},
				},
				selected: make(map[int]struct{}),
				repoDir:  "./test-repo",
			},
		},
		{
			name: "Enter with selection",
			msg:  tea.KeyMsg{Type: tea.KeyEnter},
			initialModel: Model{
				choices: []converter.AppScript{
					{Name: "Chrome", FilePath: "chrome.sh"},
					{Name: "VSCode", FilePath: "vscode.sh"},
					{Name: "Docker", FilePath: "docker.sh"},
				},
				selected: map[int]struct{}{
					0: {},
					1: {},
				},
				repoDir: "./test-repo",
			},
			expectedModel: Model{
				choices: []converter.AppScript{
					{Name: "Chrome", FilePath: "chrome.sh"},
					{Name: "VSCode", FilePath: "vscode.sh"},
					{Name: "Docker", FilePath: "docker.sh"},
				},
				selected: map[int]struct{}{
					0: {},
					1: {},
				},
				repoDir: "./test-repo",
			},
			expectQuit: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a deep copy of the initial model to prevent state leakage
			initialModel := Model{
				choices:  make([]converter.AppScript, len(tt.initialModel.choices)),
				cursor:   tt.initialModel.cursor,
				selected: make(map[int]struct{}),
				repoDir:  tt.initialModel.repoDir,
				quitting: tt.initialModel.quitting,
				err:      tt.initialModel.err,
			}
			copy(initialModel.choices, tt.initialModel.choices)
			for k, v := range tt.initialModel.selected {
				initialModel.selected[k] = v
			}

			model, cmd := initialModel.Update(tt.msg)
			resultModel := model.(Model)

			// Compare relevant fields
			assert.Equal(t, tt.expectedModel.quitting, resultModel.quitting, "Quitting state mismatch")
			assert.Equal(t, tt.expectedModel.cursor, resultModel.cursor, "Cursor position mismatch")
			assert.Equal(t, tt.expectedModel.selected, resultModel.selected, "Selection state mismatch")
			assert.Equal(t, tt.expectedModel.repoDir, resultModel.repoDir, "Repository directory mismatch")

			if tt.expectQuit {
				assert.NotNil(t, cmd, "Expected quit command, got nil")
			}
		})
	}
}

func TestModelView(t *testing.T) {
	tests := []struct {
		name        string
		model       Model
		contains    []string
		notContains []string
	}{
		{
			name: "Normal view",
			model: Model{
				choices: []converter.AppScript{
					{Name: "Chrome", FilePath: "chrome.sh"},
					{Name: "VSCode", FilePath: "vscode.sh"},
				},
				selected: make(map[int]struct{}),
			},
			contains: []string{
				"Which apps do you want to keep?",
				"Chrome",
				"VSCode",
				"↑/↓: navigate",
				"[ ]",
			},
		},
		{
			name: "View with selection",
			model: Model{
				choices: []converter.AppScript{
					{Name: "Chrome", FilePath: "chrome.sh"},
					{Name: "VSCode", FilePath: "vscode.sh"},
				},
				selected: map[int]struct{}{
					0: {},
				},
				cursor: 1,
			},
			contains: []string{
				"[x] Chrome",
				"[ ] VSCode",
				">",
			},
		},
		{
			name: "View with multiple selections",
			model: Model{
				choices: []converter.AppScript{
					{Name: "Chrome", FilePath: "chrome.sh"},
					{Name: "VSCode", FilePath: "vscode.sh"},
				},
				selected: map[int]struct{}{
					0: {},
					1: {},
				},
			},
			contains: []string{
				"[x] Chrome",
				"[x] VSCode",
			},
		},
		{
			name: "Quitting view",
			model: Model{
				quitting: true,
			},
			contains: []string{"Bye!"},
			notContains: []string{
				"Which apps",
				"navigate",
			},
		},
		{
			name: "Error view",
			model: Model{
				err: fmt.Errorf("test error"),
			},
			contains: []string{"Error:", "test error"},
			notContains: []string{
				"Which apps",
				"navigate",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			view := tt.model.View()

			// Check for required content
			for _, expected := range tt.contains {
				assert.Contains(t, view, expected, "View should contain %q", expected)
			}

			// Check for content that should not be present
			for _, unexpected := range tt.notContains {
				assert.NotContains(t, view, unexpected, "View should not contain %q", unexpected)
			}
		})
	}
}

func TestModelInit(t *testing.T) {
	model := InitialModel()
	cmd := model.Init()
	assert.Nil(t, cmd, "Expected nil command from Init")
}

func TestRunConversion(t *testing.T) {
	// Create a temporary test directory
	tempDir := t.TempDir()

	// Test successful conversion
	t.Run("Successful conversion", func(t *testing.T) {
		err := runConversion(tempDir)
		assert.NoError(t, err, "Expected no error for successful conversion")
	})

	// Test with invalid directory
	t.Run("Invalid directory", func(t *testing.T) {
		err := runConversion("/nonexistent/directory")
		assert.Error(t, err, "Expected error for invalid directory")
	})
}
