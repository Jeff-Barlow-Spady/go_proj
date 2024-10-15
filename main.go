package main

import (
	"fmt"
)

func main() {
	repoDir := "./omakub"

	// Clone the repository
	err := ubuntu_to_fedora.CloneOmakubRepo(repoDir)
	if err != nil {
		fmt.Printf("Error during repository cloning: %v\n", err)
		return
	}

	// Replace Ubuntu-specific commands
	err = ubuntu_to_fedora.ReplaceUbuntuWithFedora(repoDir)
	if err != nil {
		fmt.Printf("Error during command replacement: %v\n", err)
		return
	}

	fmt.Println("All operations completed successfully.")
}
