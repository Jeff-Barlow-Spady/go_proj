package main

import (
	"fmt"
	"omakub-fedora/ubuntuswap"
)

func main() {
	repoDir := "./omakub"

	// Clone the repository
	err := ubuntuswap.CloneOmakubRepo(repoDir)
	if err != nil {
		fmt.Printf("Error during repository cloning: %v\n", err)
		return
	}

	// Replace Ubuntu-specific commands
	err = ubuntuswap.ReplaceUbuntuWithFedora(repoDir)
	if err != nil {
		fmt.Printf("Error during command replacement: %v\n", err)
		return
	}

	fmt.Println("All operations completed successfully.")
}
