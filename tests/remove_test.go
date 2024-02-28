package tests

import (
	"database/sql"
	"github.com/govcms-tests/govcms-cli/cmd"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var db *sql.DB

func init() {
	// TODO create database in the mocked file system
	if _, err := os.Stat("test.db"); err == nil {
		_ = os.Remove("test.db")
	}
	db, _ = sql.Open("sqlite3", "test.db")
	//local = data.Initialise(db)
}

func Test_WrongName(t *testing.T) {
	rootCmd := cmd.NewRootCmd(afero.NewMemMapFs(), db)
	_, err := executeCommand(rootCmd, "remove", "ThisIsANameNoSiteShouldHave")
	assert.Error(t, err)
}

func Test_RemovalOfExistingInstall(t *testing.T) {
	rootCmd := cmd.NewRootCmd(afero.NewMemMapFs(), db)
	_, err := executeCommand(rootCmd, "get", "distribution", "exampleName")
	assert.NoError(t, err)

	_, err = executeCommand(rootCmd, "remove", "exampleName")
	assert.NoError(t, err)
}

func Test_RemovalOfMultipleSites(t *testing.T) {
	rootCmd := cmd.NewRootCmd(afero.NewMemMapFs(), db)
	_, err := executeCommand(rootCmd, "get", "distribution", "exampleName1")
	assert.NoError(t, err)

	_, err = executeCommand(rootCmd, "get", "distribution", "exampleName2")
	assert.NoError(t, err)

	_, err = executeCommand(rootCmd, "remove", "exampleName1", "exampleName2")
	assert.NoError(t, err)
}
