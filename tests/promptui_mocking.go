package tests

import (
	"bytes"
	"io"
	"strings"
)

// Sets up a reader of type io.ReadCloser with the correct buffer of string arguments
// so that it can be fed to PromptUI's input reader.
// The string arguments must end in a new line
//
// See https://stackoverflow.com/a/69505423
func addInputsToPromptReader(reader *io.ReadCloser, args ...string) {
	ensureStringsEndWithNewlines(args)
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

func ensureStringsEndWithNewlines(strings []string) {
	for index := range strings {
		if lastCharacterIsNewline(strings[index]) {
			continue
		}
		strings[index] = strings[index] + "\n"
	}
}

func lastCharacterIsNewline(str string) bool {
	return strings.HasSuffix(str, "\n")
}

func pad(size int, buf *bytes.Buffer) {
	pu := make([]byte, 4096-size)
	for i := 0; i < 4096-size; i++ {
		pu[i] = 97
	}
	buf.Write(pu)
}
