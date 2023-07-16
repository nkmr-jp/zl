package zl

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestSetOutput(t *testing.T) {
	tests := []Output{PrettyOutput, ConsoleAndFileOutput, ConsoleOutput, FileOutput}
	for _, tt := range tests {
		t.Run(tt.String(), func(t *testing.T) {
			SetOutput(tt)
			assert.Equal(t, tt, outputType)
		})
	}
}

func TestSetOutputByString(t *testing.T) {
	tests := []struct {
		in  string
		out Output
	}{
		{"Pretty", PrettyOutput},
		{"ConsoleAndFile", ConsoleAndFileOutput},
		{"Console", ConsoleOutput},
		{"File", FileOutput},
		{"", PrettyOutput},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			SetOutputByString(tt.in)
			assert.Equal(t, tt.out, outputType)
		})
	}
}

func TestSetLevel(t *testing.T) {
	tests := []zapcore.Level{DebugLevel, InfoLevel, WarnLevel, ErrorLevel, FatalLevel}
	for _, tt := range tests {
		t.Run(tt.String(), func(t *testing.T) {
			SetLevel(tt)
			assert.Equal(t, tt, severityLevel)
			ResetGlobalLoggerSettings()
		})
	}
}

func TestSetLevelByString(t *testing.T) {
	tests := []struct {
		in  string
		out zapcore.Level
	}{
		{"debug", DebugLevel},
		{"info", InfoLevel},
		{"warn", WarnLevel},
		{"error", ErrorLevel},
		{"fatal", FatalLevel},
	}
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			SetLevelByString(tt.in)
			assert.Equal(t, tt.out, severityLevel)
			ResetGlobalLoggerSettings()
		})
	}
}

func TestSetSeparator(t *testing.T) {
	ResetGlobalLoggerSettings()
	SetSeparator(":")
	assert.Equal(t, ":", separator)
}
