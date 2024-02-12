package cmd

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/govcms-tests/govcms-cli/data"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get [resource] [name]",
	Short: "get a GovCMS distribution, saas, or paas site",
	Long:  "get a GovCMS distribution, saas, or paas site.",
	//Args:      cobra.MatchAll(cobra.ExactArgs(2), cobra.OnlyValidArgs),
	//ValidArgs: []string{"distribution", "saas", "paas"},
	Run: func(cmd *cobra.Command, args []string) {
		// Validate the provided resource type
		resource := args[0]
		name := args[1]
		if resource != "distribution" && resource != "saas" && resource != "paas" {
			fmt.Println("Invalid resource type. Must be 'distribution', 'saas', or 'paas'")
			return
		}

		pathErr := os.Mkdir(name, os.ModePerm)
		if pathErr != nil {
			fmt.Println("Invalid path")
			return
		}
		// Define the target folder where repositories will be cloned
		targetFolder := name
		// Clone the corresponding repository
		repoURL := map[string]string{
			"distribution": "govCMS/GovCMS",
			"saas":         "govCMS/scaffold",
			"paas":         "govCMS/scaffold",
		}[resource]
		fmt.Printf("Cloning %s into %s\n", repoURL, targetFolder)
		_, err := git.PlainClone(targetFolder, false, &git.CloneOptions{
			URL:      "https://github.com/" + repoURL + ".git",
			Progress: os.Stdout,
		})
		if errors.Is(err, git.ErrRepositoryAlreadyExists) {
			fmt.Print("A repository with that name already exists")
			return
		}

		if err != nil {
			fmt.Printf("Error cloning repository: %s\n", err)
			return
		}

		absPath, err := filepath.Abs(targetFolder)

		res, _ := data.StringToInstallation(resource)

		install := data.Installation{
			Name:     name,
			Path:     absPath,
			Resource: res,
		}
		data.InsertInstall(install)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
}
