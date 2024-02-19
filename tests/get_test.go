package tests

import (
	"fmt"
	"github.com/govcms-tests/govcms-cli/cmd"
	"github.com/govcms-tests/govcms-cli/pkg/data"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func init() {
	data.Initialise()
}

var govcmsTypes = []string{
	"distribution",
	"saas",
	"paas",
	"lagoon",
	"tests",
}

func Test_InvalidFlagValues(t *testing.T) {
	rootCmd := cmd.NewRootCmd(afero.NewMemMapFs())
	// PR flag cannot be less than zero
	_, err := executeCommand(rootCmd, "get", "distribution", "test-site", "--pr=-1")
	assert.Error(t, err)
	assert.True(t, strings.Contains(fmt.Sprint(err), "PR number must be a positive integer"))
}

func Test_InvalidResourceType(t *testing.T) {
	rootCmd := cmd.NewRootCmd(afero.NewMemMapFs())
	// Note invalid resource type here
	_, err := executeCommand(rootCmd, "get", "someOtherType", "test-site")
	assert.Error(t, err)
}

func Test_GetWithBothFlags(t *testing.T) {
	rootCmd := cmd.NewRootCmd(afero.NewMemMapFs())

	_, err := executeCommand(rootCmd, "get", "--pr=1", "--branch=3.x-develop")
	assert.Error(t, err)
}

func Test_ValidGetNoFlags(t *testing.T) {
	rootCmd := cmd.NewRootCmd(afero.NewMemMapFs())
	for _, govcmsType := range govcmsTypes {
		_, err := executeCommand(rootCmd, "get", govcmsType, "testRepo")
		assert.NoError(t, err)
	}
}

func Test_ValidGetWithPrFlag(t *testing.T) {
	rootCmd := cmd.NewRootCmd(afero.NewMemMapFs())
	rootCmd.SetArgs([]string{})

	_, err := executeCommand(rootCmd, "get", "distribution", "testRepo", "--pr=1020")
	assert.NoError(t, err)
}

func Test_ValidGetWithBranchFlag(t *testing.T) {
	rootCmd := cmd.NewRootCmd(afero.NewMemMapFs())
	_, err := executeCommand(rootCmd, "get", "distribution", "testRepo", "--branch=3.x-develop")
	assert.NoError(t, err)
}

func Test_ValidGetWithPrompt(t *testing.T) {
	rootCmd := cmd.NewRootCmd(afero.NewMemMapFs())

	addInputsToPromptReader(&cmd.Reader, "testSite", "distribution")

	_, err := executeCommand(rootCmd, "get")
	assert.NoError(t, err)
}
