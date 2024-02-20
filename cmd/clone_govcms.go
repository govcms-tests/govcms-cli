package cmd

import (
	"errors"
	"fmt"
	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/govcms-tests/govcms-cli/pkg/config"
	"github.com/govcms-tests/govcms-cli/pkg/data"
	"github.com/govcms-tests/govcms-cli/pkg/settings"
	"github.com/govcms-tests/govcms-cli/pkg/utils"
	"io"
	"os"
	"path/filepath"
)

// Reader that allows for mocking interactions with the terminal, which is needed to test getGovcmsWithPrompt(),
// default is os.Stdin
var Reader io.ReadCloser

var memoryStorage storage.Storer

var abstractFilesystemAtRepoPath billy.Filesystem

/*
CloneGovCMS clones one of the GovCMS repo types given a correct type. It also can change the git repository's state
to a specific branch or pull request (but not both)
*/
func CloneGovCMS(name string, govcmsType string, branchName string, prNumber int) error {
	err := validateFlags(prNumber)
	if err != nil {
		return err
	}

	err = loadConfiguration()
	if err != nil {
		return err
	}

	// Construct the folder path within the workspace
	repoPath := getGovcmsRepoPath(name, prNumber, branchName)
	repoURL, err := getGovcmsURL(govcmsType)
	if err != nil {
		return err
	}

	// Clone the corresponding repository
	// TODO inject these dependencies
	memoryStorage = memory.NewStorage()
	abstractFilesystemAtRepoPath = utils.NewBillyFromAfero(AppFs, repoPath)
	err = attemptCloning(repoURL, branchName, prNumber)
	if err != nil {
		return err
	}

	// Add to database
	res, _ := data.StringToResource(govcmsType)
	err = local.InsertInstallation(data.Installation{Name: name, Path: repoPath, Resource: res})
	if err != nil {
		return err
	}

	return nil
}

func validateFlags(prNumber int) error {
	// Validate the provided PR number
	if prNumber < 0 {
		return fmt.Errorf("PR number must be a positive integer")
	}
	return nil
}

func loadConfiguration() error {
	// Load configuration from settings.go
	appConfig, err := settings.LoadConfig()
	if err != nil {
		return fmt.Errorf("error loading configuration: %v", err)
	}
	// If Workspace is empty, return an error
	if appConfig.Workspace == "" {
		return fmt.Errorf("workspace is not set in configuration")
	}
	return nil
}

func getGovcmsRepoPath(name string, prNumber int, branchName string) string {
	govcmsFolder := filepath.Join(appConfig.Workspace)

	if branchName != "" {
		name += fmt.Sprintf("_branch_%s", branchName)
	} else if prNumber != 0 {
		name += fmt.Sprintf("_pr_%d", prNumber)
	}
	fullPath := filepath.Join(govcmsFolder, name)
	fmt.Println(fullPath)
	return fullPath
}

func getGovcmsURL(govcmsType string) (string, error) {
	repoURL, ok := config.GovCMSReposList[govcmsType]
	if !ok {
		_ = fmt.Errorf("invalid govcmsType type")
		return "", errors.New("invalid type")
	}
	return repoURL, nil
}

func attemptCloning(repoURL string, branchName string, prNumber int) error {
	_, err := git.Clone(memoryStorage, abstractFilesystemAtRepoPath, &git.CloneOptions{
		URL:      "https://github.com/" + repoURL + ".git",
		Progress: os.Stdout,
	})
	if errors.Is(err, git.ErrRepositoryAlreadyExists) {
		return fmt.Errorf("repository with this name already exists at this location")
	}
	if err != nil {
		return fmt.Errorf("error cloning repository: %s", err)
	}
	return goToCorrectGitCommit(branchName, prNumber)
}

func goToCorrectGitCommit(branchName string, prNumber int) error {
	// Create local branch if needed
	if branchName != "" {
		err := utils.CreateLocalBranchIfNeeded(memoryStorage, abstractFilesystemAtRepoPath, branchName)
		if err != nil {
			return fmt.Errorf("error creating local branch: %s", err)
		}
		fmt.Println("Branch cloned successfully!")
		return nil
	}
	if prNumber != 0 {
		err := utils.GetPullRequest(memoryStorage, abstractFilesystemAtRepoPath, prNumber)
		if err != nil {
			return err
		}
		fmt.Println("Pull request cloned successfully!")
		return nil
	}
	return nil
}
