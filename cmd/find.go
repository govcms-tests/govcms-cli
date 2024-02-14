package cmd

import (
	"fmt"
	data2 "github.com/govcms-tests/govcms-cli/pkg/data"
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

		fmt.Println("Found GovCMS installations at:")
		fmt.Println(strings.Join(FindAllInstallPaths(specifiedPath), "\n"))

		allInstalls := findAllInstallations(specifiedPath)
		data2.InsertInstallations(allInstalls)
	},
}

func init() {
	rootCmd.AddCommand(findCmd)
}

func FindAllInstallPaths(root string) []string {
	var allPaths []string

	_ = filepath.WalkDir(root, func(path string, file fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if file.Name() == "govcms.info.yml" {
			parentDir := filepath.Dir(path)
			absPath, _ := filepath.Abs(parentDir)
			allPaths = append(allPaths, absPath)
			return filepath.SkipDir
		}
		return nil
	})
	return allPaths
}

func findAllInstallations(rootPath string) []data2.Installation {
	var allInstalls []data2.Installation
	allPaths := FindAllInstallPaths(rootPath)
	for _, path := range allPaths {
		name := filepath.Base(path)
		res := data2.DISTRIBUTION
		install := data2.Installation{Name: name, Path: path, Resource: res}
		allInstalls = append(allInstalls, install)
	}
	return allInstalls
}
