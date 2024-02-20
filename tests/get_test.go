package tests

import (
	"database/sql"
	"fmt"
	"github.com/govcms-tests/govcms-cli/cmd"
	"github.com/govcms-tests/govcms-cli/pkg/data"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"os"
	"strconv"
	"strings"
	"testing"
)

var local data.LocalStorage

func init() {
	// TODO create database in the mocked file system
	if _, err := os.Stat("test.db"); err == nil {
		_ = os.Remove("test.db")
	}
	db, _ := sql.Open("sqlite3", "test.db")
	local = data.Initialise(db)
}

var govcmsTypes = []string{
	"distribution",
	"saas",
	"paas",
	"lagoon",
	"tests",
}

func Test_InvalidFlagValues(t *testing.T) {
	rootCmd := cmd.NewRootCmd(afero.NewMemMapFs(), local)
	// PR flag cannot be less than zero
	_, err := executeCommand(rootCmd, "get", "distribution", "test-site", "--pr=-1")
	assert.Error(t, err)
	assert.True(t, strings.Contains(fmt.Sprint(err), "PR number must be a positive integer"))
}

func Test_InvalidResourceType(t *testing.T) {
	rootCmd := cmd.NewRootCmd(afero.NewMemMapFs(), local)
	// Note invalid resource type here
	_, err := executeCommand(rootCmd, "get", "someOtherType", "test-site")
	assert.Error(t, err)
}

func Test_GetWithBothFlags(t *testing.T) {
	rootCmd := cmd.NewRootCmd(afero.NewMemMapFs(), local)

	_, err := executeCommand(rootCmd, "get", "--pr=1", "--branch=3.x-develop")
	assert.Error(t, err)
}

func Test_ValidGetNoFlags(t *testing.T) {
	rootCmd := cmd.NewRootCmd(afero.NewMemMapFs(), local)
	for index, govcmsType := range govcmsTypes {
		_, err := executeCommand(rootCmd, "get", govcmsType, "testRepoValidGetNoFlags"+strconv.Itoa(index))
		assert.NoError(t, err)
	}
}

func Test_ValidGetWithPrFlag(t *testing.T) {
	rootCmd := cmd.NewRootCmd(afero.NewMemMapFs(), local)

	_, err := executeCommand(rootCmd, "get", "distribution", "testRepo2", "--pr=1020")
	assert.NoError(t, err)
}

func Test_ValidGetWithBranchFlag(t *testing.T) {
	rootCmd := cmd.NewRootCmd(afero.NewMemMapFs(), local)
	_, err := executeCommand(rootCmd, "get", "distribution", "testRepo3", "--branch=3.x-develop")
	assert.NoError(t, err)
}

func Test_ValidGetWithPrompt(t *testing.T) {
	rootCmd := cmd.NewRootCmd(afero.NewMemMapFs(), local)

	addInputsToPromptReader(&cmd.Reader, "testSite4", "distribution")

	_, err := executeCommand(rootCmd, "get")
	assert.NoError(t, err)
}
