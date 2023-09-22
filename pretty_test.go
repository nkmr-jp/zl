package zl

import (
	"bytes"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func Test_newPrettyLogger(t *testing.T) {
	tests := []struct {
		name           string
		setOutputType  Output
		setOmitKeys    []Key
		out            io.Writer
		expectedOutput *os.File
		expectedFlags  int
		expectedNil    bool
	}{
		{
			name:           "not PrettyOutput type",
			setOutputType:  ConsoleOutput,
			out:            os.Stderr,
			expectedOutput: os.Stderr,
			expectedNil:    true,
		},
		{
			name:           "not set options",
			setOutputType:  PrettyOutput,
			out:            os.Stderr,
			expectedOutput: os.Stderr,
			expectedFlags:  log.Ldate | log.Ltime | log.Lshortfile,
		},
		{
			name:           "set omitKeys",
			setOutputType:  PrettyOutput,
			setOmitKeys:    []Key{TimeKey},
			out:            os.Stderr,
			expectedOutput: os.Stderr,
			expectedFlags:  log.Lshortfile,
		},
		{
			name:           "set isStdOut and omitKeys",
			setOutputType:  PrettyOutput,
			setOmitKeys:    []Key{TimeKey},
			out:            os.Stdout,
			expectedOutput: os.Stdout,
			expectedFlags:  log.Lshortfile,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			outputType = tt.setOutputType
			omitKeys = tt.setOmitKeys

			// Execute
			logger := newPrettyLogger(tt.out, os.Stderr)

			// Assert
			if tt.expectedNil {
				assert.Nil(t, logger)
			} else {
				assert.Equal(t, logger.Logger.Writer(), tt.expectedOutput)
				assert.Equal(t, logger.Logger.Flags(), tt.expectedFlags)
			}
			ResetGlobalLoggerSettings()
		})
	}
}

// faultyWriter always returns an error when Write is called.
type faultyWriter struct{}

func (fw *faultyWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("forced writer error")
}

func Test_prettyLogger_log(t *testing.T) {
	tests := []struct {
		severityLevel zapcore.Level
		level         zapcore.Level
		message       string
		expectedMsg   string
	}{
		{
			zapcore.DebugLevel,
			zapcore.DebugLevel,
			"Debug Message",
			"\u001B[90mDEBUG\u001B[0m \u001B[2mDebug Message\u001B[0m",
		},
		{
			zapcore.DebugLevel,
			zapcore.InfoLevel,
			"Info Message",
			"\u001B[94mINFO\u001B[0m Info Message",
		},
		{
			zapcore.DebugLevel,
			zapcore.WarnLevel,
			"Warn Message",
			"\u001B[33mWARN\u001B[0m Warn Message",
		},
		{
			zapcore.DebugLevel,
			zapcore.ErrorLevel,
			"Error Message",
			"\u001B[31mERROR\u001B[0m Error Message",
		},
		{
			zapcore.DebugLevel,
			zapcore.FatalLevel,
			"Fatal Message",
			"\u001B[31mFATAL\u001B[0m Fatal Message",
		},
		{
			zapcore.ErrorLevel,
			zapcore.InfoLevel,
			"Debug Message",
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.severityLevel.String()+"_"+tt.level.String()+"_Level", func(t *testing.T) {
			outputType = PrettyOutput
			omitKeys = []Key{TimeKey}
			severityLevel = tt.severityLevel

			var buf bytes.Buffer
			logger := newPrettyLogger(&buf, os.Stderr)
			logger.log(tt.message, tt.level, nil)
			assert.Contains(t, buf.String(), tt.expectedMsg)
			if tt.expectedMsg == "" {
				assert.Empty(t, buf.String())
			}
			ResetGlobalLoggerSettings()
		})
	}

	t.Run("capture internal error", func(t *testing.T) {
		outputType = PrettyOutput
		severityLevel = zapcore.DebugLevel

		var buf bytes.Buffer
		l := newPrettyLogger(&faultyWriter{}, &buf)
		l.log("test message", zapcore.InfoLevel, nil)
		assert.Contains(t, buf.String(), "[INTERNAL ERROR] ")
	})
}

func Test_prettyLogger_logWithError(t *testing.T) {
	tests := []struct {
		severityLevel zapcore.Level
		level         zapcore.Level
		message       string
		err           error
		expectedMsg   string
	}{
		{
			zapcore.DebugLevel,
			zapcore.DebugLevel,
			"Debug Message",
			errors.New("some error"),
			"\u001B[90mDEBUG\u001B[0m \u001B[2mDebug Message",
		},
		{
			zapcore.DebugLevel,
			zapcore.InfoLevel,
			"Info Message",
			errors.New("some error"),
			"\u001B[94mINFO\u001B[0m Info Message",
		},
		{
			zapcore.DebugLevel,
			zapcore.WarnLevel,
			"Warn Message",
			errors.New("some error"),
			"\u001B[33mWARN\u001B[0m Warn Message",
		},
		{
			zapcore.DebugLevel,
			zapcore.ErrorLevel,
			"Error Message",
			errors.New("some error"),
			"\u001B[31mERROR\u001B[0m Error Message",
		},
		{
			zapcore.DebugLevel,
			zapcore.FatalLevel,
			"Fatal Message",
			errors.New("some error"),
			"\u001B[31mFATAL\u001B[0m Fatal Message",
		},
		{
			zapcore.ErrorLevel,
			zapcore.InfoLevel,
			"Debug Message",
			errors.New("some error"),
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.severityLevel.String()+"_"+tt.level.String()+"_Level", func(t *testing.T) {
			outputType = PrettyOutput
			omitKeys = []Key{TimeKey}
			severityLevel = tt.severityLevel

			var buf bytes.Buffer
			logger := newPrettyLogger(&buf, os.Stderr)
			logger.logWithError(tt.message, tt.level, tt.err, nil)
			assert.Contains(t, buf.String(), tt.expectedMsg)
			if tt.expectedMsg == "" {
				assert.Empty(t, buf.String())
			}
			ResetGlobalLoggerSettings()
		})
	}

	t.Run("capture internal error", func(t *testing.T) {
		outputType = PrettyOutput
		severityLevel = zapcore.DebugLevel

		var buf bytes.Buffer
		l := newPrettyLogger(&faultyWriter{}, &buf)
		l.logWithError("test message", zapcore.InfoLevel, errors.New("some error"), nil)
		assert.Contains(t, buf.String(), "[INTERNAL ERROR] ")
	})
}

func Test_prettyLogger_coloredLevel(t *testing.T) {
	// Initialize the prettyLogger instance
	logger := &prettyLogger{}

	// Create tt cases
	tests := []struct {
		level    zapcore.Level
		expected string
	}{
		{level: zapcore.FatalLevel, expected: "\u001B[31mFATAL\u001B[0m"},
		{level: zapcore.ErrorLevel, expected: "\u001B[31mERROR\u001B[0m"},
		{level: zapcore.WarnLevel, expected: "\u001B[33mWARN\u001B[0m"},
		{level: zapcore.InfoLevel, expected: "\u001B[94mINFO\u001B[0m"},
		{level: zapcore.DebugLevel, expected: "\u001B[90mDEBUG\u001B[0m"},
		{level: zapcore.Level(10 /* Undefined level */), expected: "\u001B[90m\u001B[0m"},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.level.String(), func(t *testing.T) {
			coloredString := logger.coloredLevel(tt.level).String()
			assert.Equal(t, tt.expected, coloredString)
			ResetGlobalLoggerSettings()
		})
	}
}

func Test_prettyLogger_showErrorReport(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		// Prepare expected string
		expected := "" +
			"\n" +
			"\n\u001B[1;31mERROR REPORT" +
			"\n\u001B[0m  \u001B[36mErrorCount\u001B[0m: 1" +
			"\n  \u001B[36mPID\u001B[0m: 123" +
			"\n\n" +
			"\n\u001B[1m1\u001B[0m. pretty_test.go:218: \u001B[31mERROR\u001B[0m SOME_ERROR \u001B[35msome error\u001B[0m" +
			"\n  \u001B[36mTimestamp\u001B[0m:\t2023-09-09T15:53:17.287179+09:00" +
			"\n  \u001B[36mLogFile\u001B[0m:\t/PATH/TO/PROJECT/ROOT/testdata/basic.jsonl:1" +
			"\n  \u001B[36mStackTrace\u001B[0m: " +
			"\n\tgithub.com/nkmr-jp/zl.Test_prettyLogger_showErrorReport" +
			"\n\t\t/PATH/TO/PROJECT/ROOT/pretty_test.go:218" +
			"\n\ttesting.tRunner" +
			"\n\t\t/PATH/TO/GO/ROOT/src/testing/testing.go:1595" +
			"\n\n\n"

		// Prepare Logger
		var buf bytes.Buffer
		l := newPrettyLogger(&buf, os.Stderr)

		// Execute
		fileName = "./testdata/basic.jsonl"
		l.showErrorReport(fileName, 123)
		str := buf.String()

		// Replace
		goModPath, err := exec.Command("go", "env", "GOMOD").CombinedOutput()
		assert.NoError(t, err)
		pjRoot := strings.ReplaceAll(string(goModPath), "/go.mod\n", "")
		replacedStr := strings.ReplaceAll(str, pjRoot, "/PATH/TO/PROJECT/ROOT")

		// Assert
		assert.Equal(t, expected, replacedStr)
		ResetGlobalLoggerSettings()
	})

	t.Run("error", func(t *testing.T) {
		tests := []struct {
			name     string
			fileName string
			pid      int
			output   io.Writer
			expected string
		}{
			{
				name:     "file not found",
				fileName: "./testdata/not-found.jsonl",
				pid:      123,
				output:   os.Stderr,
				expected: "no such file or directory",
			},
			{
				name:     "faulty log file",
				fileName: "./testdata/faulty.jsonl",
				pid:      123,
				output:   os.Stderr,
				expected: "unexpected end of JSON input",
			},
			{
				name:     "faulty writer",
				fileName: "./testdata/basic.jsonl",
				pid:      123,
				output:   &faultyWriter{},
				expected: "forced writer error",
			},
		}

		// Run test cases
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// Prepare
				var errBuf bytes.Buffer
				l := newPrettyLogger(tt.output, &errBuf)

				// Execute
				l.showErrorReport(tt.fileName, tt.pid)

				// Assert
				assert.Contains(t, errBuf.String(), tt.expected)
				ResetGlobalLoggerSettings()
			})
		}
	})
}

func Test_prettyLogger_consoleMsg(t *testing.T) {
	var buf bytes.Buffer
	pretty = newPrettyLogger(&buf, os.Stderr)
	consoleFields = []string{"name", "id"}

	expected := separator + "\u001B[36mAlice\u001B[0m" + separator + "\u001B[34m1\u001B[0m"
	actual := pretty.consoleMsg([]zap.Field{
		zap.String("name", "Alice"),
		zap.Int("id", 1),
	})
	assert.Equal(t, expected, actual)
}
