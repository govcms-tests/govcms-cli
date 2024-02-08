package govcms

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"

	"github.com/govcms-tests/govcms-cli/pkg/initializer"
)

// start will bring lift up
func Start() {
	fmt.Println("Start the GovCMS local development")

	initializer.Setup()

	// Run docker compose command
	dockerCmd := exec.Command("docker-compose", "-f", "assets/services/docker-compose.yml", "up", "-d")
	// Create pipes for capturing stdout and stderr
	stdout, err := dockerCmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error creating stdout pipe:", err)
		return
	}
	stderr, err := dockerCmd.StderrPipe()
	if err != nil {
		fmt.Println("Error creating stderr pipe:", err)
		return
	}

	// Start the command
	err = dockerCmd.Start()
	if err != nil {
		fmt.Println("Error starting docker compose:", err)
		return
	}

	// Print output from stdout
	go printOutput(stdout)
	// Print output from stderr
	go printOutput(stderr)

	// Wait for the command to finish
	err = dockerCmd.Wait()
	if err != nil {
		fmt.Println("Error running docker compose:", err)
		return
	}
}

// printOutput reads from a reader and prints the output
func printOutput(reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading output:", err)
	}
}
