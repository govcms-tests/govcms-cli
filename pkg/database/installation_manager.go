package database

import (
	"context"
	"database-sqlc/database"
	"database/sql"
	"errors"
	"fmt"
	"github.com/spf13/afero"
	"os"
)

const schema = `CREATE TABLE IF NOT EXISTS installations(
    name TEXT PRIMARY KEY,
    path TEXT UNIQUE NOT NULL,
    type TEXT NOT NULL,
    FOREIGN KEY (type) REFERENCES installation_type (name)
);

CREATE TABLE IF NOT EXISTS installation_type (
  name TEXT PRIMARY KEY
);`

func NewDatabase(path string) (*sql.DB, error) {
	if _, err := os.Stat(path); err != nil {
		file, err := os.Create(path)
		checkError(err)
		err = file.Close()
		if err != nil {
			return nil, fmt.Errorf("error creating database")
		}
	}
	db, err := sql.Open("sqlite3", path)
	checkError(err)

	schemaStatement, err := db.Prepare(schema)
	checkError(err)

	_, err = schemaStatement.Exec()
	checkError(err)

	return db, nil
}

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
	path, err := im.queries.GetPath(context.Background(), name)
	if err != nil {
		return fmt.Errorf("no installation found with name %s", name)
	}

	err = im.queries.DeleteInstallation(context.Background(), name)
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

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
