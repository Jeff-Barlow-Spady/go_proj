package converter_test

import (
	"os"
	"path/filepath"
	"testing"

	"ubuntu-to-fedora/pkg/converter"

	"github.com/stretchr/testify/assert"
)

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
sudo apt install nginx
sudo apt upgrade
add-apt-repository ppa:some/repo
sudo apt autoremove`

	scriptPath := filepath.Join(tempDir, "test-script.sh")
	err = os.WriteFile(scriptPath, []byte(mockScript), 0644)
	if err != nil {
		t.Fatalf("Failed to write mock script: %v", err)
	}

	// Run the replacement function
	err = converter.ReplaceUbuntuWithFedora(tempDir)
	assert.NoError(t, err, "Expected no error during replacement")

	// Read the modified script
	modifiedScript, err := os.ReadFile(scriptPath)
	if err != nil {
		t.Fatalf("Failed to read modified script: %v", err)
	}

	expected := `#!/bin/bash
sudo dnf update
sudo dnf install nginx
sudo dnf upgrade
sudo dnf config-manager --add-repo ppa:some/repo
sudo dnf autoremove`

	assert.Equal(t, expected, string(modifiedScript), "The script should have Fedora replacements")
}
