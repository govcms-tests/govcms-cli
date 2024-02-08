package govcms

import (
	"os/exec"
)

// Checks all required dependencies.
func CheckRequirements() error {
	var errors []error

	if err := checkDocker(); err != nil {
		errors = append(errors, err)
	}
	if err := checkGit(); err != nil {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return errors[0] // Returning just the first error for simplicity
	}

	return nil
}

func checkDocker() error {
	return checkCommand("docker", "--version")
}

func checkGit() error {
	return checkCommand("git", "--version")
}

func checkCommand(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	return cmd.Run()
}
