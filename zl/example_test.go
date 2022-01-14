package zl_test

import (
	"fmt"
	"os"
	"time"

	"github.com/nkmr-jp/zap-lightning/zl"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	version    string // version git revision or tag. set from go cli.
	srcRootDir string // srcRootDir set from cli.
)

func Example() {
	// Set Options
	zl.SetLogLevel(zapcore.DebugLevel) // Default is InfoLevel

	// Initialize
	zl.Init()
	defer zl.Sync()   // flush log buffer
	zl.SyncWhenStop() // flush log buffer. when interrupt or terminated.

	// Logs
	zl.Info("USER_INFO", zap.String("name", "Alice"), zap.Int("age", 20))
	err := fmt.Errorf("error message")
	zl.Error("ERROR_MESSAGE", err) // error level log must with error message.
	zl.Debug("DEBUG_MESSAGE")
	zl.Warn("WARN_MESSAGE")
	zl.WarnErr("WARN_MESSAGE_WITH_ERROR", err) // warn level log with error message.
	zl.Info("DISPLAY_TO_CONSOLE", zl.Console("display to console"))

	// Output:
}

func ExampleSetVersion() {
	// ex.`go run -ldflags "-X main.version=v0.1.1 -X main.srcRootDir=$PWD" main.go`
	version = "v0.1.1"
	srcRootDir, _ = os.Getwd()

	// Set Options
	zl.SetLogLevel(zapcore.DebugLevel) // Default is InfoLevel
	zl.SetVersion(version)
	zl.SetFileName(fmt.Sprintf("./log/app_%s.jsonl", zl.GetVersion()))
	zl.SetRepositoryCallerEncoder(
		"https://github.com/nkmr-jp/zap-lightning/blob/%s", version, srcRootDir,
	)

	// Initialize
	zl.Init()
	defer zl.Sync()   // flush log buffer
	zl.SyncWhenStop() // flush log buffer. when interrupt or terminated.

	// Logs
	zl.Info("USER_INFO", zl.Console("message"))

	// Output:
}

func ExampleSetOutputType() {
	// Set options
	zl.SetLogLevel(zapcore.DebugLevel) // Default is InfoLevel
	zl.SetOutputType(zl.OutputTypeConsole)

	// Initialize
	zl.Init()
	defer zl.Sync()   // flush log buffer
	zl.SyncWhenStop() // flush log buffer. when interrupt or terminated.

	// Logs
	zl.Info("USER_INFO", zl.Console("message"))

	// Output:
	// {"level":"INFO","ts":"2022-01-14T06:43:03.345143+09:00","caller":"zl/zl.go:36","function":"github.com/nkmr-jp/zap-lightning/zl.Init.func1","msg":"INIT_LOGGER","version":"dd90b59","hostname":"nkmrnoMacBook-Pro.local","console":"logLevel: DEBUG, fileName: , outputType: Console"}
	// {"level":"INFO","ts":"2022-01-14T06:43:03.345329+09:00","recaller":"zap-lightning/example_test.go:74","function":"github.com/nkmr-jp/zap-lightning_test.ExampleSetOutputType","msg":"USER_INFO","version":"dd90b59","hostname":"nkmrnoMacBook-Pro.local","console":"message"}
	// {"level":"INFO","ts":"2022-01-14T06:43:03.345341+09:00","caller":"zl/logger.go:67","function":"github.com/nkmr-jp/zap-lightning/zl.Sync","msg":"FLUSH_LOG_BUFFER","version":"dd90b59","hostname":"nkmrnoMacBook-Pro.local"}
}

func ExampleNew() {
	// Set options
	traceIDField := "trace_id"
	zl.AddConsoleField(traceIDField)

	// Initialize
	zl.Init()
	defer zl.Sync()
	zl.SyncWhenStop()

	// New
	// ex. Use this when you want to add a common value in the scope of a context, such as an API request.
	w := zl.New(
		zap.Int("user_id", 1),
		zap.Int64(traceIDField, time.Now().UnixNano()),
	)
	err := fmt.Errorf("context scope error message")
	w.Info("CONTEXT_SCOPE_INFO", zl.Consolef("hoge %s", err.Error()))
	w.Error("CONTEXT_SCOPE_ERROR", err)
	// Output:
}
