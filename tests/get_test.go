package tests

import (
	"github.com/govcms-tests/govcms-cli/cmd"
	"github.com/govcms-tests/govcms-cli/pkg/data"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	data.Initialise()
}

func Test_GetWithBothFlags(t *testing.T) {
	rootCmd := cmd.NewRootCmd(afero.NewMemMapFs())

	_, err := executeCommand(rootCmd, "get", "--pr=1", "--branch=3.x-develop")
	assert.Error(t, err)
}

func Test_ValidGetNoFlags(t *testing.T) {
	rootCmd := cmd.NewRootCmd(afero.NewMemMapFs())

	_, err := executeCommand(rootCmd, "get", "distribution", "notruwan")
	assert.NoError(t, err)
}

func Test_ValidGetWithPrFlag(t *testing.T) {
	rootCmd := cmd.NewRootCmd(afero.NewMemMapFs())
	rootCmd.SetArgs([]string{})

	_, err := executeCommand(rootCmd, "get", "distribution", "testRepo", "--pr=1020")
	assert.NoError(t, err)
}

func Test_ValidGetWithBranchFlag(t *testing.T) {
	rootCmd := cmd.NewRootCmd(afero.NewMemMapFs())
	rootCmd.SetArgs([]string{})

	_, err := executeCommand(rootCmd, "get", "distribution", "testRepo", "--branch=3.x-develop")
	assert.NoError(t, err)
}
