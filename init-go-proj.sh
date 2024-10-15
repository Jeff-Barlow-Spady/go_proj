#!/bin/bash

# Set project name and directories
PROJECT_NAME="omakub-fedora"
CMD_DIR="cmd"
TESTS_DIR="tests"

echo "Initializing Go project..."

# Step 1: Initialize the Go module
go mod init $PROJECT_NAME

# Step 2: Create necessary directories
mkdir -p $CMD_DIR
mkdir -p $TESTS_DIR

# Step 3: Create main.go file
cat <<EOL > main.go
package main

import (
    "fmt"
    "ubuntu_to_fedora"
)

func main() {
    repoDir := "./omakub"

    // Clone the repository
    err := ubuntu_to_fedora.CloneOmakubRepo(repoDir)
    if err != nil {
        fmt.Printf("Error during repository cloning: %v\\n", err)
        return
    }

    // Replace Ubuntu-specific commands
    err = ubuntu_to_fedora.ReplaceUbuntuWithFedora(repoDir)
    if err != nil {
        fmt.Printf("Error during command replacement: %v\\n", err)
        return
    }

    fmt.Println("All operations completed successfully.")
}
EOL

# Step 4: Create the core logic file (ubuntu_to_fedora.go)
cat <<EOL > ubuntu_to_fedora.go
package ubuntu_to_fedora

import (
    "fmt"
    "os"
    "os/exec"
    "path/filepath"
    "github.com/go-git/go-git/v5"
    "errors"
    "io/ioutil"
    "strings"
)

// CloneOmakubRepo clones the repository to the specified directory with improved error handling.
func CloneOmakubRepo(destDir string) error {
    repoURL := "https://github.com/omakub/omakub.git" // Replace with actual repo URL

    // Check if the destination directory exists and is non-empty
    if _, err := os.Stat(destDir); !os.IsNotExist(err) {
        files, err := os.ReadDir(destDir)
        if err != nil {
            return fmt.Errorf("failed to read destination directory: %v", err)
        }
        if len(files) > 0 {
            return fmt.Errorf("destination directory %s is not empty", destDir)
        }
    }

    // Check if Git is installed
    if _, err := exec.LookPath("git"); err != nil {
        return errors.New("git is not installed on this system")
    }

    fmt.Printf("Cloning omakub repository to %s...\\n", destDir)

    _, err := git.PlainClone(destDir, false, &git.CloneOptions{
        URL: repoURL,
    })
    if err != nil {
        return fmt.Errorf("failed to clone the repository: %v", err)
    }

    fmt.Println("Repository cloned successfully!")
    return nil
}

// ReplaceUbuntuWithFedora walks through the directory and replaces Ubuntu-specific commands in shell scripts.
func ReplaceUbuntuWithFedora(dir string) error {
    err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return fmt.Errorf("error accessing the path %s: %v", path, err)
        }

        // Only process regular files that are shell scripts
        if info.Mode().IsRegular() && strings.HasSuffix(info.Name(), ".sh") {
            fmt.Printf("Processing file: %s\\n", path)
            return replaceCommandsInFile(path)
        }

        return nil
    })

    if err != nil {
        return fmt.Errorf("error processing directory %s: %v", dir, err)
    }

    fmt.Println("Replacement completed successfully.")
    return nil
}

// replaceCommandsInFile reads a file and replaces Ubuntu commands with Fedora equivalents
func replaceCommandsInFile(filePath string) error {
    content, err := ioutil.ReadFile(filePath)
    if err != nil {
        return fmt.Errorf("failed to read file: %v", err)
    }

    original := string(content)

    // Replace Ubuntu-specific commands with Fedora equivalents
    replacements := map[string]string{
        "sudo apt update":       "sudo dnf update",
        "sudo apt upgrade":      "sudo dnf upgrade",
        "sudo apt install":      "sudo dnf install",
        "add-apt-repository":    "sudo dnf config-manager --add-repo",
        "sudo apt autoremove":   "sudo dnf autoremove",
    }

    modified := original
    for ubuntuCmd, fedoraCmd := range replacements {
        modified = strings.ReplaceAll(modified, ubuntuCmd, fedoraCmd)
    }

    if modified != original {
        fmt.Printf("Modifying file: %s\\n", filePath)
        err = ioutil.WriteFile(filePath, []byte(modified), 0644)
        if err != nil {
            return fmt.Errorf("failed to write modified file: %v", err)
        }
    } else {
        fmt.Printf("No Ubuntu-specific commands found in %s\\n", filePath)
    }

    return nil
}
EOL

# Step 5: Create a TUI skeleton (cmd/tui.go)
cat <<EOL > $CMD_DIR/tui.go
package cmd

import (
    "fmt"
    "github.com/charmbracelet/bubbletea"
    "ubuntu_to_fedora"
)

// Improved TUI model with graceful error handling
type model struct {
    choices  []string   // Available packages/apps
    selected map[int]struct{} // User-selected packages
    err      error      // Capture any errors that occur
    quitting bool       // Quit flag
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
            err := runConversion()
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
    if m.quitting {
        return "Bye!\\n"
    }

    if m.err != nil {
        return fmt.Sprintf("Error: %v\\n", m.err)
    }

    s := "Which apps do you want to keep? Press Enter to confirm.\\n\\n"
    for i, choice := range m.choices {
        selected := "[ ] "
        if _, ok := m.selected[i]; ok {
            selected = "[x] "
        }
        s += fmt.Sprintf("%s%s\\n", selected, choice)
    }

    return s
}

func runConversion() error {
    repoDir := "./omakub"
    err := ubuntu_to_fedora.CloneOmakubRepo(repoDir)
    if err != nil {
        return fmt.Errorf("error cloning repository: %v", err)
    }

    err = ubuntu_to_fedora.ReplaceUbuntuWithFedora(repoDir)
    if err != nil {
        return fmt.Errorf("error replacing Ubuntu-specific commands: %v", err)
    }

    fmt.Println("Conversion completed successfully.")
    return nil
}
EOL

# Step 6: Create test files in tests/
cat <<EOL > $TESTS_DIR/ubuntu_to_fedora_test.go
package ubuntu_to_fedora_test

import (
    "os"
    "testing"
    "io/ioutil"
    "path/filepath"
    "github.com/stretchr/testify/assert"
    "ubuntu_to_fedora"
)

// TestCloneOmakubRepo tests the repository cloning logic
func TestCloneOmakubRepo(t *testing.T) {
    // Create a temporary directory for cloning
    tempDir, err := ioutil.TempDir("", "omakub-test-repo")
    if err != nil {
        t.Fatalf("Failed to create temp dir: %v", err)
    }
    defer os.RemoveAll(tempDir) // Clean up

    err = ubuntu_to_fedora.CloneOmakubRepo(tempDir)
    assert.NoError(t, err, "Expected no error during repo clone")
    // Check if the repo was cloned (by checking for a README.md or .git directory)
    assert.DirExists(t, filepath.Join(tempDir, ".git"), "Repository should be cloned")
}

// TestReplaceUbuntuWithFedora tests the command replacement logic
func TestReplaceUbuntuWithFedora(t *testing.T) {
    // Create a temporary directory with a mock script
    tempDir, err := ioutil.TempDir("", "replace-test")
    if err != nil {
        t.Fatalf("Failed to create temp dir: %v", err)
    }
    defer os.RemoveAll(tempDir) // Clean up

    // Create a mock Ubuntu script
    mockScript := `#!/bin/bash
    sudo apt update
    sudo apt install curl
    add-apt-repository ppa:example/ppa
    `
    scriptPath := filepath.Join(tempDir, "test-script.sh")
    err = ioutil.WriteFile(scriptPath, []byte(mockScript), 0644)
    if err != nil {
        t.Fatalf("Failed to write mock script: %v", err)
    }

    // Run the replacement function
    err = ubuntu_to_fedora.ReplaceUbuntuWithFedora(tempDir)
    assert.NoError(t, err, "Expected no error during replacement")

    // Read the modified script
    modifiedScript, err := ioutil.ReadFile(scriptPath)
    if err != nil {
        t.Fatalf("Failed to read modified script: %v", err)
    }

    // Check if the replacements were made correctly
    expected := `#!/bin/bash
    sudo dnf update
    sudo dnf install curl
    sudo dnf config-manager --add-repo ppa:example/ppa
    `
    assert.Equal(t, expected, string(modifiedScript), "The script should have Fedora replacements")
}
EOL

# Step 7: Install necessary Go modules
go get github.com/go-git/go-git/v5
go get github.com/charmbracelet/bubbletea
go get github.com/stretchr/testify/assert

echo "Go project initialized successfully!"

