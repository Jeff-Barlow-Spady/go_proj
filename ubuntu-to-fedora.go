package ubuntu-to-fedora

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

    fmt.Printf("Cloning omakub repository to %s...\n", destDir)

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
            fmt.Printf("Processing file: %s\n", path)
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
        fmt.Printf("Modifying file: %s\n", filePath)
        err = ioutil.WriteFile(filePath, []byte(modified), 0644)
        if err != nil {
            return fmt.Errorf("failed to write modified file: %v", err)
        }
    } else {
        fmt.Printf("No Ubuntu-specific commands found in %s\n", filePath)
    }

    return nil
}
