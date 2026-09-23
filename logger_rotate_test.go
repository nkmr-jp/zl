package zl

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// Loggers made by New must write through the same file writer as the global logger.
// With a writer per logger, the one that rotates moves on to the new file while the
// others keep writing to the renamed backup, and their lines never reach the new file.
func TestNew_sharesFileWriterAcrossRotation(t *testing.T) {
	ResetGlobalLoggerSettings()
	defer ResetGlobalLoggerSettings()
	dir := t.TempDir()
	file := dir + "/app.jsonl"
	SetOutput(FileOutput)
	SetRotateFileName(file)
	SetRotateMaxSize(1) // megabytes
	SetOmitKeys(VersionKey, HostnameKey, PIDKey)
	Init()

	a := New(zap.String("who", "a"))
	b := New(zap.String("who", "b"))
	Info("GLOBAL_BEFORE")
	b.Info("B_BEFORE")

	// Push a past MaxSize so that a rotation happens while a writes.
	pad := strings.Repeat("x", 1024)
	for range 1100 {
		a.Info("PAD", zap.String("pad", pad))
	}
	backups, err := os.ReadDir(dir)
	require.NoError(t, err)
	require.Len(t, backups, 2, "a rotation must have happened")

	Info("GLOBAL_AFTER")
	a.Info("A_AFTER")
	_, _, line, _ := runtime.Caller(0)
	b.Info("B_AFTER") // the caller must be this line
	Sync()

	data, err := os.ReadFile(file)
	require.NoError(t, err)
	callers := map[string]string{} // message -> caller, for lines in the current file
	for line := range strings.SplitSeq(strings.TrimSpace(string(data)), "\n") {
		var entry struct{ Message, Caller string }
		require.NoError(t, json.Unmarshal([]byte(line), &entry))
		callers[entry.Message] = entry.Caller
	}
	for _, msg := range []string{"GLOBAL_AFTER", "A_AFTER", "B_AFTER"} {
		assert.Contains(t, callers, msg, "%s must be in the current file", msg)
	}
	assert.True(t,
		strings.HasSuffix(callers["B_AFTER"], fmt.Sprintf("/logger_rotate_test.go:%d", line+1)),
		"caller %q must point at the line that called the logger made by New", callers["B_AFTER"])
}
