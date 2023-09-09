package zl

import (
	"bytes"
	"errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
	"log"
	"testing"
)

func Test_newPrettyLogger(t *testing.T) {
	outputType = PrettyOutput
	logger := newPrettyLogger()

	assert.NotNil(t, logger)
	assert.NotNil(t, logger.Logger)
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
		})
	}
}

//func Test_prettyLogger_showErrorReport(t *testing.T) {
//	Init()
//	outputType = PrettyOutput
//	var buf bytes.Buffer
//	pretty = &prettyLogger{
//		Logger: log.New(&buf, "", log.Ldate|log.Ltime|log.Lshortfile),
//	}
//
//	Err("SOME_ERROR", errors.New("some error"))
//
//	//logger.logWithError(zapcore.ErrorLevel, errors.New("some error"), nil)
//	pretty.showErrorReport()
//
//	//fmt.Println(buf.String())
//
//	//assert.Contains(t, buf.String(), tt.expected)
//	//assert.Equal(t, "", coloredString)
//
//}
