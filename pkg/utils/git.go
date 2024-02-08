package utils

import (
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func CreateLocalBranchIfNeeded(repoPath, branchName string) error {
	// Open the repository
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return fmt.Errorf("error opening repository: %s", err)
	}

	// Construct the local reference name based on branchName
	localRef := plumbing.ReferenceName("refs/heads/" + branchName)

	// Check if the branch already exists locally
	_, err = repo.Reference(localRef, false)
	if err == nil {
		fmt.Printf("Branch '%s' already exists locally.\n", branchName)

		// Checkout the branch
		w, err := repo.Worktree()
		if err != nil {
			return fmt.Errorf("error getting worktree: %s", err)
		}
		err = w.Checkout(&git.CheckoutOptions{
			Branch: localRef,
			Create: false,
		})
		if err != nil {
			return fmt.Errorf("error checking out branch '%s': %s", branchName, err)
		}
		fmt.Printf("Checked out branch '%s'.\n", branchName)

		return nil
	} else if err != plumbing.ErrReferenceNotFound {
		return fmt.Errorf("error checking if branch exists locally: %s", err)
	}

	// Fetch the latest changes from the remote repository
	err = repo.Fetch(&git.FetchOptions{})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return fmt.Errorf("error fetching latest changes from remote repository: %s", err)
	}

	// Get the reference for the remote branch
	branchRef := plumbing.ReferenceName("refs/remotes/origin/" + branchName)
	ref, err := repo.Reference(branchRef, true)
	if err != nil {
		return fmt.Errorf("error getting reference for remote branch: %s", err)
	}

	// Resolve the revision for the remote branch
	hash := ref.Hash()

	// Create a new local branch pointing to the same commit as the remote branch
	newRef := plumbing.NewHashReference(localRef, hash)
	err = repo.Storer.SetReference(newRef)
	if err != nil {
		return fmt.Errorf("error creating local branch: %s", err)
	}

	// Checkout the branch
	w, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("error getting worktree: %s", err)
	}
	err = w.Checkout(&git.CheckoutOptions{
		Branch: localRef,
		Create: false,
	})
	if err != nil {
		return fmt.Errorf("error checking out branch '%s': %s", branchName, err)
	}
	fmt.Printf("Checked out branch '%s'.\n", branchName)

	return nil
}
