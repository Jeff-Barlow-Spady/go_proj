package converter

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
)

type AppScript struct {
	Name     string
	FilePath string
}

func GetAvailableApps(dir string) ([]AppScript, error) {
	seen := make(map[string]bool)
	var apps []AppScript

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing path %s: %v", path, err)
		}

		if strings.HasPrefix(info.Name(), ".") {
			return nil
		}

		if info.Mode().IsRegular() && strings.HasSuffix(info.Name(), ".sh") {
			name := strings.TrimSuffix(info.Name(), ".sh")
			name = strings.ReplaceAll(name, "_", " ")
			name = strings.ReplaceAll(name, "-", " ")
			name = strings.Title(strings.ToLower(name))

			if !seen[name] {
				seen[name] = true
				apps = append(apps, AppScript{
					Name:     name,
					FilePath: path,
				})
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error scanning directory: %v", err)
	}

	return apps, nil
}

func CloneOmakubRepo(destDir string) error {
	repoURL := "https://github.com/basecamp/omakub.git"

	if _, err := os.Stat(destDir); !os.IsNotExist(err) {
		files, err := os.ReadDir(destDir)
		if err != nil {
			return fmt.Errorf("failed to read destination directory: %v", err)
		}
		if len(files) > 0 {
			return fmt.Errorf("destination directory %s is not empty", destDir)
		}
	}

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

func ReplaceUbuntuWithFedora(dir string) error {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error accessing the path %s: %v", path, err)
		}

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

func replaceCommandsInFile(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	original := string(content)

	replacements := map[string]string{
		"sudo apt update":     "sudo dnf update",
		"sudo apt upgrade":    "sudo dnf upgrade",
		"sudo apt install":    "sudo dnf install",
		"add-apt-repository":  "sudo dnf config-manager --add-repo",
		"sudo apt autoremove": "sudo dnf autoremove",
		"sudo apt-get":        "sudo dnf",
		"apt-get":             "dnf",
		"sudo apt":            "sudo dnf",
		"apt":                 "dnf",
	}

	modified := original
	keys := make([]string, 0, len(replacements))
	for k := range replacements {
		keys = append(keys, k)
	}
	for i := 0; i < len(keys)-1; i++ {
		for j := i + 1; j < len(keys); j++ {
			if len(keys[i]) < len(keys[j]) {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
	}

	for _, ubuntuCmd := range keys {
		modified = strings.ReplaceAll(modified, ubuntuCmd, replacements[ubuntuCmd])
	}

	if modified != original {
		fmt.Printf("Modifying file: %s\n", filePath)
		err = os.WriteFile(filePath, []byte(modified), 0644)
		if err != nil {
			return fmt.Errorf("failed to write modified file: %v", err)
		}
	} else {
		fmt.Printf("No Ubuntu-specific commands found in %s\n", filePath)
	}

	return nil
}
