package utils

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
)

// ExecuteCommand starts a command, captures its output, and waits for it to finish
func ExecuteCommand(cmd *exec.Cmd) {
	// Create pipes for capturing stdout and stderr
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("Error creating stdout pipe:", err)
		return
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println("Error creating stderr pipe:", err)
		return
	}

	// Start the command
	err = cmd.Start()
	if err != nil {
		fmt.Println("Error starting command:", err)
		return
	}

	// Print output from stdout
	go PrintOutput(stdout)
	// Print output from stderr
	go PrintOutput(stderr)

	// Wait for the command to finish
	err = cmd.Wait()
	if err != nil {
		fmt.Println("Error running command:", err)
		return
	}
}

// PrintOutput reads from a reader and prints the output
func PrintOutput(reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading output:", err)
	}
}
