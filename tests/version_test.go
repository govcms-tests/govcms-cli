package tests

import (
	"github.com/govcms-tests/govcms-cli/cmd"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Version(t *testing.T) {
	output, _ := executeCommand(cmd.RootCmd, "version")
	expectedVersion := "v0.1.0 -- HEAD"

	assert.Equal(t, output, expectedVersion, "Version does not match the expected version")
}
