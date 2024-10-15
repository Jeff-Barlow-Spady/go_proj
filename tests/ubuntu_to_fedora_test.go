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
    mockScript := Hit:1 https://download.docker.com/linux/ubuntu noble InRelease
Hit:2 https://typora.io/linux ./ InRelease
Hit:3 https://updates.signal.org/desktop/apt xenial InRelease
Hit:4 https://dl.google.com/linux/chrome/deb stable InRelease
Get:5 http://security.ubuntu.com/ubuntu noble-security InRelease [126 kB]
Hit:6 https://mise.jdx.dev/deb stable InRelease
Hit:7 http://ca.archive.ubuntu.com/ubuntu noble InRelease
Hit:8 https://cli.github.com/packages stable InRelease
Hit:9 https://packages.microsoft.com/repos/code stable InRelease
Get:10 http://ca.archive.ubuntu.com/ubuntu noble-updates InRelease [126 kB]
Hit:11 https://brave-browser-apt-release.s3.brave.com stable InRelease
Hit:12 https://ppa.launchpadcontent.net/agornostal/ulauncher/ubuntu noble InRelease
Hit:13 https://ppa.launchpadcontent.net/zhangsongcui3371/fastfetch/ubuntu noble InRelease
Hit:14 http://ca.archive.ubuntu.com/ubuntu noble-backports InRelease
Get:15 http://security.ubuntu.com/ubuntu noble-security/main amd64 Packages [410 kB]
Hit:16 http://repository.spotify.com stable InRelease
Get:17 http://ca.archive.ubuntu.com/ubuntu noble-updates/main amd64 Packages [592 kB]
Get:18 http://ca.archive.ubuntu.com/ubuntu noble-updates/main Translation-en [144 kB]
Get:19 http://ca.archive.ubuntu.com/ubuntu noble-updates/restricted amd64 Packages [385 kB]
Get:20 http://security.ubuntu.com/ubuntu noble-security/main Translation-en [90.4 kB]
Get:21 http://security.ubuntu.com/ubuntu noble-security/universe amd64 Packages [553 kB]
Get:22 http://ca.archive.ubuntu.com/ubuntu noble-updates/restricted Translation-en [74.4 kB]
Get:23 http://ca.archive.ubuntu.com/ubuntu noble-updates/universe amd64 Packages [697 kB]
Get:24 http://security.ubuntu.com/ubuntu noble-security/universe Translation-en [147 kB]
Get:25 http://ca.archive.ubuntu.com/ubuntu noble-updates/multiverse amd64 Packages [14.8 kB]
Fetched 3,360 kB in 4s (950 kB/s)
Reading package lists...
Building dependency tree...
Reading state information...
10 packages can be upgraded. Run 'apt list --upgradable' to see them.
Reading package lists...
Building dependency tree...
Reading state information...
curl is already the newest version (8.5.0-2ubuntu10.4).
0 upgraded, 0 newly installed, 0 to remove and 10 not upgraded.
Error: must run as root
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
    expected := 
    assert.Equal(t, expected, string(modifiedScript), "The script should have Fedora replacements")
}
