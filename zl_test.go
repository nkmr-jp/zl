package zl

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResetGlobalLoggerSettings(t *testing.T) {
	outputType = ConsoleOutput
	ResetGlobalLoggerSettings()
	assert.Equal(t, PrettyOutput, outputType)
}

func TestResetCleanup(t *testing.T) {
	outputType = ConsoleOutput
	Cleanup()
	assert.Equal(t, PrettyOutput, outputType)
}
