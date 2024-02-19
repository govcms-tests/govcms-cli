package cmd

import (
	"errors"
	"fmt"
	"github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/govcms-tests/govcms-cli/pkg/config"
	"github.com/govcms-tests/govcms-cli/pkg/data"
	"github.com/govcms-tests/govcms-cli/pkg/settings"
	"github.com/govcms-tests/govcms-cli/pkg/utils"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func NewGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [resource] [name]",
		Short: "get a GovCMS distribution, saas, or paas site",
		Long:  "get a GovCMS distribution, saas, or paas site.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if hasBothFlags(cmd) {
				return fmt.Errorf("cannot specify both --pr and --branch flags together")
			}
			if len(args) < 2 {
				err := getGovcmsWithPrompt()
				return err
			}
			err := getGovcmsWithoutPrompt(cmd, args)
			return err
		},
	}

	cmd.Flags().IntP("pr", "p", 0, "Github PR number")
	cmd.Flags().StringP("branch", "b", "", "Git branch name")

	return cmd
}

//// getCmd represents the get command
//var getCmd = &cobra.Command{
//	Use:   "get [resource] [name]",
//	Short: "get a GovCMS distribution, saas, or paas site",
//	Long:  "get a GovCMS distribution, saas, or paas site.",
//	RunE: func(cmd *cobra.Command, args []string) error {
//		if hasBothFlags(cmd) {
//			return fmt.Errorf("cannot specify both --pr and --branch flags together")
//		}
//		if len(args) < 2 {
//			err := getGovcmsWithPrompt()
//			return err
//		}
//		err := getGovcmsWithoutPrompt(cmd, args)
//		return err
//	},
//}

//func init() {
//	getCmd.Flags().IntP("pr", "p", 0, "Github PR number")
//	getCmd.Flags().StringP("branch", "b", "", "Git branch name")
//}

func hasBothFlags(cmd *cobra.Command) bool {
	prNumber, _ := cmd.Flags().GetInt("pr")
	branch, _ := cmd.Flags().GetString("branch")
	log.Printf("Flags: pr - %v, branch - %v", prNumber, branch)
	return prNumber != 0 && branch != ""
}

func getGovcmsWithPrompt() error {
	name, govcmsType := getDetailsFromPrompt()
	err := Generate(name, govcmsType, 0, "")
	return err
}

func getDetailsFromPrompt() (string, string) {
	name := getNameFromPrompt()
	govcmsType := getTypeFromPrompt()
	return name, govcmsType
}

func getNameFromPrompt() string {
	prompt := promptui.Prompt{
		Label: "What would you like to call this installation?",
	}
	name, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		panic(err)
	}
	fmt.Printf("You chose %q\n", name)
	return name
}

func getTypeFromPrompt() string {
	prompt := promptui.Select{
		Label: "Which type would you like to install?",
		Items: []string{"Distribution", "SaaS", "PaaS", "Lagoon", "Tests", "Scaffold-Tooling"},
	}
	_, govcmsType, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		panic(err)
	}
	fmt.Printf("You choose %q\n", govcmsType)
	return strings.ToLower(govcmsType)
}

func getGovcmsWithoutPrompt(cmd *cobra.Command, args []string) error {
	govcmsType := args[0]
	name := args[1]
	prNumber, _ := cmd.Flags().GetInt("pr")
	branchName, _ := cmd.Flags().GetString("branch")
	// Call the generate function from the govcms package
	err := Generate(name, govcmsType, prNumber, branchName)
	return err
}

func Generate(name string, govcmsType string, prNumber int, branchName string) error {
	fmt.Printf("Cloning repo type %s as %s with PR=%v and Branch %v. \n", govcmsType, name, prNumber, branchName)

	err := validateFlags(prNumber, branchName)
	if err != nil {
		return err
	}

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

	// Get the repository URL from the repos package
	repoURL, ok := config.GovCMSReposList[govcmsType]
	if !ok {
		return fmt.Errorf("invalid govcmsType type")
	}

	// Clone the corresponding repository
	var repoPath string
	if branchName != "" {
		name = fmt.Sprintf("%s_branch_%s", govcmsType, branchName)
	} else if prNumber != 0 {
		name = fmt.Sprintf("%s_pr_%d", govcmsType, prNumber)
	}

	repoPath = filepath.Join(govcmsFolder, name)

	// Print the cloning message
	fmt.Printf("Cloning %s into %s\n", repoURL, repoPath)

	memoryStorage := memory.NewStorage()
	abstractFilesystemAtRepoPath := utils.NewBillyFromAfero(AppFs, repoPath)

	_, err = git.Clone(memoryStorage, abstractFilesystemAtRepoPath, &git.CloneOptions{
		URL:      "https://github.com/" + repoURL + ".git",
		Progress: os.Stdout,
	})

	if errors.Is(err, git.ErrRepositoryAlreadyExists) {
		return fmt.Errorf("repository with this name already exists at this location")
	}

	// Handle errors
	if err != nil {
		return fmt.Errorf("error cloning repository: %s", err)
	}

	fmt.Println("This print statement is executed")
	res, _ := data.StringToResource(govcmsType)
	data.InsertInstallation(data.Installation{Name: name, Path: repoPath, Resource: res})

	// Create local branch if needed
	if branchName != "" {
		err = utils.CreateLocalBranchIfNeeded(memoryStorage, abstractFilesystemAtRepoPath, branchName)
		if err != nil {
			return fmt.Errorf("error creating local branch: %s", err)
		}
		fmt.Println("Branch cloned successfully!")
	} else if prNumber != 0 {
		// Open the repository
		r, err := git.Open(memoryStorage, abstractFilesystemAtRepoPath)
		if err != nil {
			return fmt.Errorf("error opening repository: %s", err)
		}

		// Fetch the PR reference from the remote repository
		prRefSpec := fmt.Sprintf("refs/pull/%d/head:refs/remotes/origin/pr/%d", prNumber, prNumber)

		err = r.Fetch(&git.FetchOptions{
			RefSpecs: []gitconfig.RefSpec{gitconfig.RefSpec(prRefSpec)},
			Force:    true,
		})

		if err != nil && err != git.NoErrAlreadyUpToDate {
			return fmt.Errorf("error fetching latest changes from remote repository: %s", err)
		}

		// Get the PR reference from the fetched references
		prRefName := plumbing.ReferenceName(fmt.Sprintf("refs/remotes/origin/pr/%d", prNumber))
		prRef, err := r.Reference(prRefName, true)
		if err != nil {
			return fmt.Errorf("error getting PR reference: %s", err)
		}

		// Create a new local branch from the fetched PR branch
		localBranchRefName := plumbing.ReferenceName(fmt.Sprintf("refs/heads/pr-%d", prNumber))
		localBranchRef := plumbing.NewHashReference(localBranchRefName, prRef.Hash())
		if err := r.Storer.SetReference(localBranchRef); err != nil {
			return fmt.Errorf("error creating local branch: %s", err)
		}

		// Checkout the newly created local branch
		w, err := r.Worktree()
		if err != nil {
			return fmt.Errorf("error getting worktree: %s", err)
		}
		err = w.Checkout(&git.CheckoutOptions{
			Branch: localBranchRefName,
			Create: false, // Local branch has been created explicitly
		})
		if err != nil {
			return fmt.Errorf("error checking out PR branch: %s", err)
		}

		fmt.Printf("Checked out PR branch %s\n", localBranchRefName)
		fmt.Println("Pull request cloned successfully!")
	} else {
		fmt.Println("Default branch cloned successfully!")
	}

	return nil
}

func validateFlags(prNumber int, branch string) error {
	// Validate the provided PR number
	if prNumber < 0 {
		return fmt.Errorf("PR number must be a positive integer")
	}
	return nil
}
