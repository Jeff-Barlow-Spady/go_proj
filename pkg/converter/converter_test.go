package converter_test

import (
	"os"
	"path/filepath"
	"testing"

	"ubuntu-to-fedora/pkg/converter"

	"github.com/stretchr/testify/assert"
)

// TestGetAvailableApps tests the app discovery functionality
func TestGetAvailableApps(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "apps-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test files
	files := map[string]string{
		"chrome.sh":             "#!/bin/bash\necho 'Installing Chrome'",
		"visual_studio_code.sh": "#!/bin/bash\necho 'Installing VSCode'",
		"not-a-script.txt":      "This is not a script",
		"docker-ce.sh":          "#!/bin/bash\necho 'Installing Docker'",
		".hidden.sh":            "#!/bin/bash\necho 'Hidden script'", // Test hidden file
	}

	for name, content := range files {
		err := os.WriteFile(filepath.Join(tempDir, name), []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", name, err)
		}
	}

	// Test app discovery
	apps, err := converter.GetAvailableApps(tempDir)
	assert.NoError(t, err, "Expected no error during app discovery")

	// Verify number of apps (should only count visible .sh files)
	assert.Equal(t, 3, len(apps), "Expected 3 apps to be discovered")

	// Create a map for easier testing
	appMap := make(map[string]string)
	for _, app := range apps {
		appMap[app.Name] = app.FilePath
	}

	// Test specific apps
	expectedApps := map[string]bool{
		"Chrome":             false,
		"Visual Studio Code": false,
		"Docker Ce":          false,
	}

	for _, app := range apps {
		if _, ok := expectedApps[app.Name]; ok {
			expectedApps[app.Name] = true
		} else {
			t.Errorf("Unexpected app found: %s", app.Name)
		}
	}

	// Verify all expected apps were found
	for appName, found := range expectedApps {
		assert.True(t, found, "Expected to find app: %s", appName)
	}

	// Test empty directory
	emptyDir, err := os.MkdirTemp("", "empty-test")
	if err != nil {
		t.Fatalf("Failed to create empty temp dir: %v", err)
	}
	defer os.RemoveAll(emptyDir)

	emptyApps, err := converter.GetAvailableApps(emptyDir)
	assert.NoError(t, err, "Expected no error for empty directory")
	assert.Empty(t, emptyApps, "Expected no apps in empty directory")

	// Test non-existent directory
	_, err = converter.GetAvailableApps("/nonexistent/directory")
	assert.Error(t, err, "Expected error for non-existent directory")

	// Test directory with no permissions
	noPermDir, err := os.MkdirTemp("", "noperm-test")
	if err != nil {
		t.Fatalf("Failed to create no-permission temp dir: %v", err)
	}
	defer os.RemoveAll(noPermDir)

	err = os.Chmod(noPermDir, 0000)
	if err != nil {
		t.Fatalf("Failed to change directory permissions: %v", err)
	}

	_, err = converter.GetAvailableApps(noPermDir)
	assert.Error(t, err, "Expected error for directory without read permissions")
}

// TestReplaceUbuntuWithFedora tests the command replacement logic
func TestReplaceUbuntuWithFedora(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "Basic commands",
			input: `#!/bin/bash
sudo apt update
sudo apt install nginx
sudo apt upgrade
add-apt-repository ppa:some/repo
sudo apt autoremove`,
			expected: `#!/bin/bash
sudo dnf update
sudo dnf install nginx
sudo dnf upgrade
sudo dnf config-manager --add-repo ppa:some/repo
sudo dnf autoremove`,
		},
		{
			name: "Complex commands with options",
			input: `#!/bin/bash
sudo apt-get update -y
sudo apt install -y --no-install-recommends nginx
apt list --installed
sudo apt-get install -y docker-ce docker-ce-cli`,
			expected: `#!/bin/bash
sudo dnf update -y
sudo dnf install -y --no-install-recommends nginx
dnf list --installed
sudo dnf install -y docker-ce docker-ce-cli`,
		},
		{
			name: "Mixed content",
			input: `#!/bin/bash
# Update system
sudo apt update
echo "Installing packages..."
sudo apt install package1 package2
# Some other commands
ls -la
wget https://example.com/file
sudo apt autoremove`,
			expected: `#!/bin/bash
# Update system
sudo dnf update
echo "Installing packages..."
sudo dnf install package1 package2
# Some other commands
ls -la
wget https://example.com/file
sudo dnf autoremove`,
		},
		{
			name: "No Ubuntu commands",
			input: `#!/bin/bash
echo "Hello World"
ls -la
wget https://example.com/file`,
			expected: `#!/bin/bash
echo "Hello World"
ls -la
wget https://example.com/file`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary directory
			tempDir, err := os.MkdirTemp("", "replace-test")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tempDir)

			// Create test script
			scriptPath := filepath.Join(tempDir, "test-script.sh")
			err = os.WriteFile(scriptPath, []byte(tt.input), 0644)
			if err != nil {
				t.Fatalf("Failed to write test script: %v", err)
			}

			// Run the replacement function
			err = converter.ReplaceUbuntuWithFedora(tempDir)
			assert.NoError(t, err, "Expected no error during replacement")

			// Read the modified script
			modifiedScript, err := os.ReadFile(scriptPath)
			if err != nil {
				t.Fatalf("Failed to read modified script: %v", err)
			}

			assert.Equal(t, tt.expected, string(modifiedScript), "The script should have correct Fedora replacements")
		})
	}

	// Test error cases
	t.Run("Non-existent directory", func(t *testing.T) {
		err := converter.ReplaceUbuntuWithFedora("/nonexistent/directory")
		assert.Error(t, err, "Expected error for non-existent directory")
	})

	t.Run("No permission to write", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "noperm-test")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		// Create read-only script
		scriptPath := filepath.Join(tempDir, "readonly.sh")
		err = os.WriteFile(scriptPath, []byte("sudo apt update"), 0444)
		if err != nil {
			t.Fatalf("Failed to create read-only script: %v", err)
		}

		err = converter.ReplaceUbuntuWithFedora(tempDir)
		assert.Error(t, err, "Expected error when writing to read-only file")
	})
}

// TestCloneOmakubRepo tests the repository cloning functionality
func TestCloneOmakubRepo(t *testing.T) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "clone-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test cloning to non-empty directory
	t.Run("Non-empty directory", func(t *testing.T) {
		err = os.WriteFile(filepath.Join(tempDir, "existing-file.txt"), []byte("test"), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		err = converter.CloneOmakubRepo(tempDir)
		assert.Error(t, err, "Expected error when cloning to non-empty directory")
		assert.Contains(t, err.Error(), "not empty", "Expected 'not empty' error message")
	})

	// Test cloning to directory without write permissions
	t.Run("No write permission", func(t *testing.T) {
		noPermDir, err := os.MkdirTemp("", "noperm-test")
		if err != nil {
			t.Fatalf("Failed to create no-permission temp dir: %v", err)
		}
		defer os.RemoveAll(noPermDir)

		err = os.Chmod(noPermDir, 0555) // Read + execute only
		if err != nil {
			t.Fatalf("Failed to change directory permissions: %v", err)
		}

		err = converter.CloneOmakubRepo(noPermDir)
		assert.Error(t, err, "Expected error when cloning to directory without write permission")
	})

	// Test cloning to new directory
	t.Run("Successful clone", func(t *testing.T) {
		newDir := filepath.Join(tempDir, "new-dir")
		err = converter.CloneOmakubRepo(newDir)
		// Note: This test might be skipped if there's no internet connection
		if err != nil {
			t.Skipf("Clone test skipped (possibly no internet): %v", err)
		}
	})
}
