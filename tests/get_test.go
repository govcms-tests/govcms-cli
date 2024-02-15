package tests

import (
	"github.com/govcms-tests/govcms-cli/cmd"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_GetBothFlags(t *testing.T) {
	_, err := executeCommand(cmd.RootCmd, "get", "--pr=1", "--branch=3.x-develop")
	assert.Error(t, err)
}

func Test_ValidGet(t *testing.T) {
	_, err := executeCommand(cmd.RootCmd, "get", "distribution", "asdfasdfasdf")
	assert.NoError(t, err)
}
