package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoot(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOutput(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{})

	Execute()
	assert.Contains(t, buf.String(), "Usage:")
	assert.Contains(t, buf.String(), "Flags:")
}
