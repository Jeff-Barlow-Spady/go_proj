package ubuntuswaptest

import (
    "os"
    "testing"
    "path/filepath"
    "github.com/stretchr/testify/assert"
    "omakub-fedora/ubuntuswap" // Use the correct relative import path
)

// TestCloneOmakubRepo tests the repository cloning logic
func TestCloneOmakubRepo(t *testing.T) {
    // Create a temporary directory for cloning
    tempDir, err := os.MkdirTemp("", "omakub-test-repo")
    if (err != nil) {
        t.Fatalf("Failed to create temp dir: %v", err)
    }
    defer os.RemoveAll(tempDir) // Clean up

    err = ubuntuswap.CloneOmakubRepo(tempDir)
    assert.NoError(t, err, "Expected no error during repo clone")
    // Check if the repo was cloned (by checking for a README.md or .git directory)
    assert.DirExists(t, filepath.Join(tempDir, ".git"), "Repository should be cloned")
}

// TestReplaceUbuntuWithFedora tests the command replacement logic
func TestReplaceUbuntuWithFedora(t *testing.T) {
    // Create a temporary directory with a mock script
    tempDir, err := os.MkdirTemp("", "replace-test")
    if err != nil {
        t.Fatalf("Failed to create temp dir: %v", err)
    }
    defer os.RemoveAll(tempDir) // Clean up

    // Create a mock Ubuntu script
    mockScript := `#!/bin/bash
    sudo apt update
    sudo apt install -y somepackage
    `
    scriptPath := filepath.Join(tempDir, "mock_script.sh")
    err = os.WriteFile(scriptPath, []byte(mockScript), 0644)
    if err != nil {
        t.Fatalf("Failed to write mock script: %v", err)
    }

    // Run the replacement function
    err = ubuntuswap.ReplaceUbuntuWithFedora(tempDir)
    assert.NoError(t, err, "Expected no error during command replacement")

    // Check if the script was modified correctly
    modifiedContent, err := os.ReadFile(scriptPath)
    if err != nil {
        t.Fatalf("Failed to read modified script: %v", err)
    }
    assert.Contains(t, string(modifiedContent), "sudo dnf update", "Expected Fedora command in modified script")
}
