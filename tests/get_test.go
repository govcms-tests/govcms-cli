package tests

import (
	"fmt"
	"github.com/govcms-tests/govcms-cli/cmd"
	"github.com/stretchr/testify/assert"
	"strconv"
	"strings"
	"testing"
)

var govcmsTypes = []string{
	"distribution",
	"saas",
	"paas",
	"lagoon",
	"tests",
}

func Test_InvalidFlagValues(t *testing.T) {
	rootCmd := setUp()
	// PR flag cannot be less than zero
	_, err := executeCommand(rootCmd, "get", "distribution", "test-site", "--pr=-1")
	assert.Error(t, err)
	assert.True(t, strings.Contains(fmt.Sprint(err), "PR number must be a positive integer"))
}

func Test_InvalidResourceType(t *testing.T) {
	rootCmd := setUp()
	// Note invalid resource type here
	_, err := executeCommand(rootCmd, "get", "someOtherType", "test-site")
	assert.Error(t, err)
}

func Test_GetWithBothFlags(t *testing.T) {
	rootCmd := setUp()

	_, err := executeCommand(rootCmd, "get", "--pr=1", "--branch=3.x-develop")
	assert.Error(t, err)
}

func Test_ValidGetNoFlags(t *testing.T) {
	rootCmd := setUp()
	for index, govcmsType := range govcmsTypes {
		_, err := executeCommand(rootCmd, "get", govcmsType, "testRepoValidGetNoFlags"+strconv.Itoa(index))
		assert.NoError(t, err)
	}
}

func Test_ValidGetWithPrFlag(t *testing.T) {
	rootCmd := setUp()

	_, err := executeCommand(rootCmd, "get", "distribution", "testRepo2", "--pr=1020")
	assert.NoError(t, err)
}

func Test_ValidGetWithBranchFlag(t *testing.T) {
	rootCmd := setUp()
	_, err := executeCommand(rootCmd, "get", "distribution", "testRepo3", "--branch=3.x-develop")
	assert.NoError(t, err)
}

func Test_ValidGetWithPrompt(t *testing.T) {
	rootCmd := setUp()

	addInputsToPromptReader(&cmd.Reader, "testSite4", "distribution")

	_, err := executeCommand(rootCmd, "get")
	assert.NoError(t, err)
}
