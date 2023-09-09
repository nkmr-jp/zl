package zl

import (
	"bytes"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
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
		setIsStdOut    bool
		expectedOutput *os.File
		expectedFlags  int
		expectedNil    bool
	}{
		{
			name:          "not PrettyOutput type",
			setOutputType: ConsoleOutput,
			expectedNil:   true,
		},
		{
			name:           "not set options",
			setOutputType:  PrettyOutput,
			expectedOutput: os.Stderr,
			expectedFlags:  log.Ldate | log.Ltime | log.Lshortfile,
		},
		{
			name:           "set omitKeys",
			setOutputType:  PrettyOutput,
			setOmitKeys:    []Key{TimeKey},
			expectedOutput: os.Stderr,
			expectedFlags:  log.Lshortfile,
		},
		{
			name:           "set isStdOut and omitKeys",
			setOutputType:  PrettyOutput,
			setOmitKeys:    []Key{TimeKey},
			setIsStdOut:    true,
			expectedOutput: os.Stdout,
			expectedFlags:  log.Lshortfile,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			outputType = tt.setOutputType
			omitKeys = tt.setOmitKeys
			isStdOut = tt.setIsStdOut

			// Execute
			logger := newPrettyLogger()

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

func Test_prettyLogger_log(t *testing.T) {
	tests := []struct {
		level       zapcore.Level
		message     string
		expectedMsg string
	}{
		{zapcore.DebugLevel, "Debug Message", "\u001B[90mDEBUG\u001B[0m \u001B[2mDebug Message\u001B[0m"},
		{zapcore.InfoLevel, "Info Message", "\u001B[94mINFO\u001B[0m Info Message"},
		{zapcore.WarnLevel, "Warn Message", "\u001B[33mWARN\u001B[0m Warn Message"},
		{zapcore.ErrorLevel, "Error Message", "\u001B[31mERROR\u001B[0m Error Message"},
		{zapcore.FatalLevel, "Fatal Message", "\u001B[31mFATAL\u001B[0m Fatal Message"},
	}

	for _, tt := range tests {
		t.Run(tt.level.String()+"_Level", func(t *testing.T) {
			outputType = PrettyOutput
			severityLevel = tt.level

			var buf bytes.Buffer
			logger := &prettyLogger{
				Logger: log.New(&buf, "", log.Ldate|log.Ltime|log.Lshortfile),
			}

			logger.log(tt.message, tt.level, nil)
			assert.Contains(t, buf.String(), tt.expectedMsg)
			ResetGlobalLoggerSettings()
		})
	}
}

func Test_prettyLogger_logWithError(t *testing.T) {
	tests := []struct {
		level       zapcore.Level
		message     string
		err         error
		expectedMsg string
	}{
		{
			zapcore.DebugLevel,
			"Debug Message",
			errors.New("some error"),
			"\u001B[90mDEBUG\u001B[0m \u001B[2mDebug Message \u001B[35msome error\u001B[0m",
		},
		{
			zapcore.InfoLevel,
			"Info Message",
			errors.New("some error"),
			"\u001B[94mINFO\u001B[0m Info Message \u001B[35msome error\u001B[0m",
		},
		{
			zapcore.WarnLevel,
			"Warn Message",
			errors.New("some error"),
			"\u001B[33mWARN\u001B[0m Warn Message \u001B[35msome error\u001B[0m",
		},
		{
			zapcore.ErrorLevel,
			"Error Message",
			errors.New("some error"),
			"\u001B[31mERROR\u001B[0m Error Message \u001B[35msome error\u001B[0m",
		},
		{
			zapcore.FatalLevel,
			"Fatal Message",
			errors.New("some error"),
			"\u001B[31mFATAL\u001B[0m Fatal Message \u001B[35msome error\u001B[0m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.level.String()+"_Level_With_Error", func(t *testing.T) {
			outputType = PrettyOutput
			severityLevel = tt.level

			var buf bytes.Buffer
			logger := &prettyLogger{
				Logger: log.New(&buf, "", log.Ldate|log.Ltime|log.Lshortfile),
			}

			logger.logWithError(tt.message, tt.level, tt.err, nil)
			assert.Contains(t, buf.String(), tt.expectedMsg)
			ResetGlobalLoggerSettings()
		})
	}
}

func Test_prettyLogger_coloredLevel(t *testing.T) {
	// Initialize the prettyLogger instance
	logger := &prettyLogger{}

	// Create test cases
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
	for _, test := range tests {
		t.Run(test.level.String(), func(t *testing.T) {
			coloredString := logger.coloredLevel(test.level).String()
			assert.Equal(t, test.expected, coloredString)
			ResetGlobalLoggerSettings()
		})
	}
}

func Test_prettyLogger_showErrorReport(t *testing.T) {
	// Prepare expected string
	expected := "" +
		"\n" +
		"\n\u001B[1;31mERROR REPORT" +
		"\n\u001B[0m  \u001B[36mErrorCount\u001B[0m: 1" +
		"\n  \u001B[36mPID\u001B[0m: 16169" +
		"\n\n" +
		"\n\u001B[1m1\u001B[0m. pretty_test.go:218: \u001B[31mERROR\u001B[0m SOME_ERROR \u001B[35msome error\u001B[0m" +
		"\n  \u001B[36mTimestamp\u001B[0m:\t2023-09-09T15:53:17.287179+09:00" +
		"\n  \u001B[36mLogFile\u001B[0m:\t/PATH/TO/PROJECT/ROOT/testdata/pretty-showErrorReport.jsonl:1" +
		"\n  \u001B[36mStackTrace\u001B[0m: " +
		"\n\tgithub.com/nkmr-jp/zl.Test_prettyLogger_showErrorReport" +
		"\n\t\t/PATH/TO/PROJECT/ROOT/pretty_test.go:218" +
		"\n\ttesting.tRunner" +
		"\n\t\t/PATH/TO/GO/ROOT/src/testing/testing.go:1595" +
		"\n\n\n"

	// Prepare Logger
	var buf bytes.Buffer
	pretty = &prettyLogger{
		Logger: log.New(&buf, "", log.Ldate|log.Ltime|log.Lshortfile),
	}

	// Execute
	fileName = "./testdata/pretty-showErrorReport.jsonl"
	pretty.showErrorReport(fileName, 16169)
	str := buf.String()

	// Replace
	goModPath, err := exec.Command("go", "env", "GOMOD").CombinedOutput()
	assert.NoError(t, err)
	pjRoot := strings.ReplaceAll(string(goModPath), "/go.mod\n", "")
	replacedStr := strings.ReplaceAll(str, pjRoot, "/PATH/TO/PROJECT/ROOT")

	// Assert
	assert.Equal(t, expected, replacedStr)
	ResetGlobalLoggerSettings()
}
