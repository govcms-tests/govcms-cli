package database

import (
	"context"
	"database-sqlc/database"
	"database/sql"
	"errors"
	"github.com/spf13/afero"
	"os"
)

type InstallationManager struct {
	queries *database.Queries
	appFs   afero.Fs
}

func NewInstallationManager(db *sql.DB, fs afero.Fs) *InstallationManager {
	im := new(InstallationManager)

	im.queries = database.New(db)
	im.appFs = fs

	im.Sync()
	return im
}

//======================================================================================================================

func (im *InstallationManager) Sync() {
	listOfPaths, _ := im.queries.ListPaths(context.Background())
	for _, path := range listOfPaths {
		im.removePathIfMissing(path)
	}
}

func (im *InstallationManager) removePathIfMissing(path string) {
	if exists, _ := im.dirExists(path); !exists {
		_ = im.queries.DeletePath(context.Background(), path)
	}
}

func (im *InstallationManager) dirExists(path string) (bool, error) {
	_, err := im.appFs.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

//======================================================================================================================

// CreateInstallation adds a new installation to the database
func (im *InstallationManager) CreateInstallation(name string, path string, installationType string) error {
	_, err := im.queries.CreateInstallation(context.Background(), database.CreateInstallationParams{
		Name: name,
		Path: path,
		Type: installationType,
	})
	return err
}

// DeleteInstallation removes a local installation from the database and the filesystem, if it exists
func (im *InstallationManager) DeleteInstallation(name string) error {
	// Get installation path so that we can later delete it from the filesystem
	path, _ := im.queries.GetPath(context.Background(), name)
	//
	err := im.queries.DeleteInstallation(context.Background(), name)
	if err != nil {
		return err
	}
	// Only remove installation from filesystem if installation was successfully removed from database
	return im.appFs.RemoveAll(path)
}

// GetAllPaths returns a list of paths of all installations
func (im *InstallationManager) GetAllPaths() ([]string, error) {
	paths, err := im.queries.ListPaths(context.Background())
	return paths, err
}

func (im *InstallationManager) GetPath(name string) (string, error) {
	path, err := im.queries.GetPath(context.Background(), name)
	return path, err
}
