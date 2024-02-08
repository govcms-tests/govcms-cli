package initializer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/govcms-tests/govcms-cli/pkg/settings"
	"github.com/govcms-tests/govcms-cli/pkg/utils"
)

// Setup initializes configuration for the CLI tool
func Setup() error {
	// Load configuration from settings.go
	appConfig, err := settings.LoadConfig()
	if err != nil {
		return fmt.Errorf("error loading configuration: %v", err)
	}

	// If Workspace is empty, return an error
	if appConfig.Workspace == "" {
		return fmt.Errorf("workspace is not set in configuration")
	}

	// Checks and creates the govcms folder
	if err := checkGovCMSFolder(); err != nil {
		return fmt.Errorf("error checking govcms folder: %v", err)
	}

	return nil
}

// Checks and creates the folder specified by the Workspace field in the settings package
func checkGovCMSFolder() error {
	// Load configuration from settings.go
	appConfig, err := settings.LoadConfig()
	if err != nil {
		return fmt.Errorf("error loading configuration: %v", err)
	}

	// If Workspace is empty, return an error
	if appConfig.Workspace == "" {
		return fmt.Errorf("workspace is not set in configuration")
	}

	// Construct the folder path within the workspace
	govcmsFolder := filepath.Join(appConfig.Workspace)

	// Create the folder if it doesn't exist
	if _, err := os.Stat(govcmsFolder); os.IsNotExist(err) {
		err := os.MkdirAll(govcmsFolder, 0755)
		if err != nil {
			return fmt.Errorf("error creating 'govcms' folder: %v", err)
		}
		fmt.Println("Created 'govcms' folder at", govcmsFolder)
	}

	// Construct the absolute path for the assets folder
	assetsFolder := filepath.Join(appConfig.Workspace, "assets")
	// Copy the assets directory using the copyDir function
	if err := utils.CopyDir("assets", assetsFolder); err != nil {
		return fmt.Errorf("error copying assets: %v", err)
	}

	return nil
}
