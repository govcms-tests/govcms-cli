package utils

import (
	"fmt"
	"os/exec"
	"runtime"
)

// Returns the name of the current operating system
func GetOperatingSystem() string {
	switch os := runtime.GOOS; os {
	case "darwin":
		return "macOS"
	case "linux":
		return "Linux"
	case "windows":
		return "Windows"
	default:
		return os
	}
}

// Opens the default web browser with the specified URL
func OpenBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default:
		return fmt.Errorf("unsupported platform")
	}

	err := cmd.Start()
	if err != nil {
		return err
	}
	return nil
}
