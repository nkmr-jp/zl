package zl

import (
	"bytes"
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

	// Point the report at a file that does not exist: opening it is observable as an
	// internal error, so "no internal error" proves Sync did not read the file.
	var out, errBuf bytes.Buffer
	pretty = newPrettyLogger(&out, &errBuf)
	fileName = dir + "/missing.jsonl"

	Info("SOME_INFO")
	Sync()
	assert.Equal(t, int64(0), errorCount.Load())
	assert.Empty(t, errBuf.String(), "Sync must not read the log file when no error was written")

	Error("SOME_ERROR")
	ErrorErr("SOME_ERROR", assert.AnError)
	Sync()
	assert.Equal(t, int64(2), errorCount.Load())
	assert.Contains(t, errBuf.String(), "no such file or directory", "Sync reads the log file once an error was written")

	ResetGlobalLoggerSettings()
	assert.Equal(t, int64(0), errorCount.Load(), "reset clears the counter")
}
