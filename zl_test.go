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

func TestSync_skipsErrorReportWhenNoErrorWasWritten(t *testing.T) {
	// The report is built by reading the log file back. A process that wrote no ERROR
	// must not read (and possibly choke on) a file shared with other processes.
	ResetGlobalLoggerSettings()
	dir := t.TempDir()
	SetRotateFileName(dir + "/app.jsonl")
	SetOmitKeys(VersionKey, HostnameKey, PIDKey)
	Init()
	defer ResetGlobalLoggerSettings()

	Info("SOME_INFO")
	assert.Equal(t, int64(0), errorCount.Load())

	Error("SOME_ERROR")
	ErrorErr("SOME_ERROR", assert.AnError)
	assert.Equal(t, int64(2), errorCount.Load())

	ResetGlobalLoggerSettings()
	assert.Equal(t, int64(0), errorCount.Load(), "reset clears the counter")
}
