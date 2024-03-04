package tests

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_WrongName(t *testing.T) {
	rootCmd := setUp()
	_, err := executeCommand(rootCmd, "remove", "ThisIsANameNoSiteShouldHave")
	assert.Error(t, err)
}

func Test_RemovalOfExistingInstall(t *testing.T) {
	rootCmd := setUp()
	_, err := executeCommand(rootCmd, "get", "distribution", "exampleName")
	assert.NoError(t, err)

	_, err = executeCommand(rootCmd, "remove", "exampleName")
	assert.NoError(t, err)
}

func Test_RemovalOfMultipleSites(t *testing.T) {
	rootCmd := setUp()
	_, err := executeCommand(rootCmd, "get", "distribution", "exampleName1")
	assert.NoError(t, err)

	_, err = executeCommand(rootCmd, "get", "distribution", "exampleName2")
	assert.NoError(t, err)

	_, err = executeCommand(rootCmd, "remove", "exampleName1", "exampleName2")
	assert.NoError(t, err)
}
