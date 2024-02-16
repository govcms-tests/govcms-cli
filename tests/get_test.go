package tests

import (
	"github.com/govcms-tests/govcms-cli/cmd"
	"github.com/govcms-tests/govcms-cli/pkg/data"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GetBothFlags(t *testing.T) {
	rootCmd := cmd.NewRootCmd(afero.NewMemMapFs())

	_, err := executeCommand(rootCmd, "get", "--pr=1", "--branch=3.x-develop")
	assert.Error(t, err)
}

func Test_ValidGet(t *testing.T) {
	data.Initialise()

	rootCmd := cmd.NewRootCmd(afero.NewMemMapFs())

	_, err := executeCommand(rootCmd, "get", "distribution", "notruwan")
	assert.NoError(t, err)
}
