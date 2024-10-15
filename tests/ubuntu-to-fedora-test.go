package ubuntuswaptest

import (
    "os"
    "testing"
    "io/ioutil"
    "path/filepath"
    "github.com/stretchr/testify/assert"
    "go_proj/ubuntuswap"
)

// TestCloneOmakubRepo tests the repository cloning logic
func TestCloneOmakubRepo(t *testing.T) {
    // Create a temporary directory for cloning
    tempDir, err := ioutil.TempDir("", "omakub-test-repo")
    if err != nil {
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
    tempDir, err := ioutil.TempDir("", "replace-test")
    if err != nil {
        t.Fatalf("Failed to create temp dir: %v", err)
    }
    defer os.RemoveAll(tempDir) // Clean up

    // Create a mock Ubuntu script
mockScript := Get:1 https://packages.microsoft.com/repos/microsoft-ubuntu-focal-prod focal InRelease [3632 B]
Get:2 https://dl.yarnpkg.com/debian stable InRelease [17.1 kB]
Get:3 https://packages.microsoft.com/repos/microsoft-ubuntu-focal-prod focal/main amd64 Packages [317 kB]
Get:4 https://repo.anaconda.com/pkgs/misc/debrepo/conda stable InRelease [3961 B]
Get:5 http://archive.ubuntu.com/ubuntu focal InRelease [265 kB]
Get:6 https://packages.microsoft.com/repos/microsoft-ubuntu-focal-prod focal/main all Packages [2942 B]
Get:7 http://security.ubuntu.com/ubuntu focal-security InRelease [128 kB]
Get:8 https://dl.yarnpkg.com/debian stable/main all Packages [11.8 kB]
Get:9 https://dl.yarnpkg.com/debian stable/main amd64 Packages [11.8 kB]
Get:10 https://repo.anaconda.com/pkgs/misc/debrepo/conda stable/main amd64 Packages [4557 B]
Get:12 http://archive.ubuntu.com/ubuntu focal-updates InRelease [128 kB]
Get:11 https://packagecloud.io/github/git-lfs/ubuntu focal InRelease [28.0 kB]
Get:13 http://security.ubuntu.com/ubuntu focal-security/multiverse amd64 Packages [30.9 kB]
Get:14 http://archive.ubuntu.com/ubuntu focal-backports InRelease [128 kB]
Get:15 http://security.ubuntu.com/ubuntu focal-security/main amd64 Packages [4027 kB]
Get:16 http://archive.ubuntu.com/ubuntu focal/multiverse amd64 Packages [177 kB]
Get:18 http://archive.ubuntu.com/ubuntu focal/main amd64 Packages [1275 kB]
Get:17 https://packagecloud.io/github/git-lfs/ubuntu focal/main amd64 Packages [3690 B]
Get:19 http://archive.ubuntu.com/ubuntu focal/universe amd64 Packages [11.3 MB]
Get:20 http://security.ubuntu.com/ubuntu focal-security/restricted amd64 Packages [4036 kB]
Get:21 http://security.ubuntu.com/ubuntu focal-security/universe amd64 Packages [1274 kB]
Get:22 http://archive.ubuntu.com/ubuntu focal/restricted amd64 Packages [33.4 kB]
Get:23 http://archive.ubuntu.com/ubuntu focal-updates/restricted amd64 Packages [4235 kB]
Get:24 http://archive.ubuntu.com/ubuntu focal-updates/main amd64 Packages [4527 kB]
Get:25 http://archive.ubuntu.com/ubuntu focal-updates/universe amd64 Packages [1566 kB]
Get:26 http://archive.ubuntu.com/ubuntu focal-updates/multiverse amd64 Packages [33.5 kB]
Get:27 http://archive.ubuntu.com/ubuntu focal-backports/main amd64 Packages [55.2 kB]
Get:28 http://archive.ubuntu.com/ubuntu focal-backports/universe amd64 Packages [28.6 kB]
Fetched 33.7 MB in 3s (10.7 MB/s)
Reading package lists...
Building dependency tree...
Reading state information...
18 packages can be upgraded. Run 'apt list --upgradable' to see them.
Reading package lists...
Building dependency tree...
Reading state information...
curl is already the newest version (7.68.0-1ubuntu2.24).
0 upgraded, 0 newly installed, 0 to remove and 18 not upgraded.
Error: must run as root
    scriptPath := filepath.Join(tempDir, "test-script.sh")
    err = ioutil.WriteFile(scriptPath, []byte(mockScript), 0644)
    if err != nil {
        t.Fatalf("Failed to write mock script: %v", err)
    }

    // Run the replacement function
    err = ubuntuswap.ReplaceUbuntuWithFedora(tempDir)
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
