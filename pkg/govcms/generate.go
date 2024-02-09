package govcms

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/govcms-tests/govcms-cli/pkg/config"
	"github.com/govcms-tests/govcms-cli/pkg/settings"
	"github.com/govcms-tests/govcms-cli/pkg/utils"
)

func Generate(resource string, prNumber int, branchName string) error {
	// Validate the provided PR number
	if prNumber < 0 {
		return fmt.Errorf("PR number must be a positive integer")
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
	repoURL, ok := config.GovCMSReposList[resource]
	if !ok {
		return fmt.Errorf("invalid resource type")
	}

	// Clone the corresponding repository
	var repoPath string
	if branchName != "" {
		// If a branch name is provided, set repoPath to the branch name folder
		repoPath = filepath.Join(govcmsFolder, fmt.Sprintf("%s_branch_%s", resource, branchName))
	} else if prNumber != 0 {
		// If a PR number is provided, clone to a folder with resource name plus PR number
		repoPath = filepath.Join(govcmsFolder, fmt.Sprintf("%s_pr_%d", resource, prNumber))
	} else {
		// If no PR number is provided, clone to a folder with just the resource name
		repoPath = filepath.Join(govcmsFolder, resource)
	}

	// Print the cloning message
	fmt.Printf("Cloning %s into %s\n", repoURL, repoPath)

	// Clone the repository
	_, err = git.PlainClone(repoPath, false, &git.CloneOptions{
		URL:      "https://github.com/" + repoURL + ".git",
		Progress: os.Stdout,
	})

	// Handle errors
	if err != nil && err != git.ErrRepositoryAlreadyExists {
		return fmt.Errorf("error cloning repository: %s", err)
	}

	// Create local branch if needed
	if branchName != "" {
		err = utils.CreateLocalBranchIfNeeded(repoPath, branchName)
		if err != nil {
			return fmt.Errorf("error creating local branch: %s", err)
		}
		fmt.Println("Branch cloned successfully!")
	} else if prNumber != 0 {
		// Open the repository
		r, err := git.PlainOpen(repoPath)
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