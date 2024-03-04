package tests

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Version(t *testing.T) {
	rootCmd := setUp()

	actualOutput, _ := executeCommand(rootCmd, "version")
	expectedOutput := "v0.1.0 -- HEAD\n"

	assert.Equal(t, expectedOutput, actualOutput, "Version does not match the expected version")
}
