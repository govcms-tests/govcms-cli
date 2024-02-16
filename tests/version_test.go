package tests

import (
	"github.com/govcms-tests/govcms-cli/cmd"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Version(t *testing.T) {
	rootCmd := cmd.NewRootCmd(afero.NewMemMapFs())

	actualOutput, _ := executeCommand(rootCmd, "version")
	expectedOutput := "v0.1.0 -- HEAD\n"

	assert.Equal(t, expectedOutput, actualOutput, "Version does not match the expected version")
}
