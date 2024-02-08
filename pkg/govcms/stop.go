package govcms

import (
	"fmt"
	"os/exec"
)

// Stop will stop and destroy the containers created by Start()
func Stop() {
	fmt.Println("Stopping the GovCMS local development")

	// Run docker compose down command
	dockerCmd := exec.Command("docker-compose", "-f", "assets/services/docker-compose.yml", "down")

	// Run the command and wait for it to finish
	err := dockerCmd.Run()
	if err != nil {
		fmt.Println("Error stopping docker compose:", err)
		return
	}

	fmt.Println("GovCMS local development stopped successfully")
}
