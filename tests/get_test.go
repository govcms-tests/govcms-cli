package tests

import (
	"bytes"
	"github.com/govcms-tests/govcms-cli/cmd"
	"github.com/govcms-tests/govcms-cli/pkg/data"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"io"
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

	addInputsToPromptReader(&cmd.Reader, "testSite\n", "distribution\n")

	_, err := executeCommand(rootCmd, "get")
	assert.NoError(t, err)
}

// Sets up a reader of type io.ReadCloser with the correct buffer of string arguments
// so that it can be fed to PromptUI's input reader.
// See https://stackoverflow.com/a/69505423
func addInputsToPromptReader(reader *io.ReadCloser, args ...string) {
	var b *bytes.Buffer
	for index, arg := range args {
		if index == 0 {
			b = bytes.NewBuffer([]byte(arg))
			pad(len(arg), b)
			*reader = io.NopCloser(b)
			continue
		}
		b.WriteString(arg)
		pad(len(arg), b)
	}
}

func pad(siz int, buf *bytes.Buffer) {
	pu := make([]byte, 4096-siz)
	for i := 0; i < 4096-siz; i++ {
		pu[i] = 97
	}
	buf.Write(pu)
}
