package govcms

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/govcms-tests/govcms-cli/pkg/utils"
)

func CopyDistributionFiles(path string) {
	sourceDir := "assets/distribution"
	destDir := path

	filepath.Walk(sourceDir, func(srcPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(sourceDir, srcPath)
		if err != nil {
			return err
		}

		destPath := filepath.Join(destDir, relPath)
		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		sourceFile, err := os.Open(srcPath)
		if err != nil {
			return err
		}
		defer sourceFile.Close()

		destFile, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, sourceFile)
		return err
	})
}

func StartContainer(path string) {
	// Create a command to run docker-compose
	dockerComposeCmd := exec.Command("docker-compose", "-f", path+"/docker-compose.yml", "up", "-d")

	// Execute the command using the ExecuteCommand function from exec.go
	utils.ExecuteCommand(dockerComposeCmd)
}

// StopContainer stops the Docker container for the GovCMS distribution
func StopContainer(path string) {
	// Run docker compose command to stop the container
	dockerComposeCmd := exec.Command("docker-compose", "-f", path+"/docker-compose.yml", "down")

	// Execute the command using the ExecuteCommand function from exec.go
	utils.ExecuteCommand(dockerComposeCmd)
}
