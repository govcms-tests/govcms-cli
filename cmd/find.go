package cmd

import (
	"fmt"
	"github.com/govcms-tests/govcms-cli/pkg/data"
	"github.com/govcms-tests/govcms-cli/pkg/settings"
	"github.com/spf13/cobra"
	"io/fs"
	"path/filepath"
	"strings"
)

var findCmd = &cobra.Command{
	Use:   "find [root_path]",
	Short: "Find all GovCMS installations from a given path",
	Long:  "Find all GovCMS installations from a given path. An 'installation' is defined to be any directory containing a 'govcms.info.yml' file",
	//Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var specifiedPath string
		// If no argument is provided, assume the user wants to check from current location
		if len(args) < 1 {
			specifiedPath = "."
		} else {
			specifiedPath = args[0]
		}
		allInstalls := findAllInstallations(specifiedPath)
		allInstallURLs := FindAllInstallPaths(specifiedPath)

		if len(allInstallURLs) == 0 {
			fmt.Println("No installations found.")
			return
		}
		fmt.Println("Found GovCMS installations at:")
		fmt.Println(strings.Join(allInstallURLs, "\n"))

		for _, install := range allInstalls {
			err := installationManager.CreateInstallation(install.Name, install.Path, install.Type)
			if err != nil {
				return
			}
		}

	},
}

// FindAllInstallPaths returns a list of absolute url strings of all locally stored GovCMS installations
func FindAllInstallPaths(root string) []string {
	var allPaths []string
	recursiveSearchForGovcms(root, &allPaths)
	// Then check specified govcms workspace folder
	cfg, _ := settings.LoadConfig()
	recursiveSearchForGovcms(cfg.Workspace, &allPaths)
	return allPaths
}

func recursiveSearchForGovcms(root string, allPaths *[]string) {
	filepath.WalkDir(root, func(path string, file fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if file.Name() == "govcms.info.yml" {
			parentDir := filepath.Dir(path)
			absPath, _ := filepath.Abs(parentDir)
			*allPaths = append(*allPaths, absPath)
			return filepath.SkipDir
		}
		return nil
	})
}

func findAllInstallations(rootPath string) []data.Installation {
	var allInstalls []data.Installation
	allPaths := FindAllInstallPaths(rootPath)
	for _, path := range allPaths {
		name := filepath.Base(path)
		res := "distribution"
		install := data.Installation{Name: name, Path: path, Type: res}
		allInstalls = append(allInstalls, install)
	}
	return allInstalls
}
