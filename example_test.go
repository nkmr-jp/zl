package lightning_test

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

const (
	urlFormat = "https://github.com/nkmr-jp/zap-lightning/blob/%s"
)

func Example() {
	// If you use this, you might actually want to pass the value from the go command, like this
	// ex.`go run -ldflags "-X main.version=v1.0.0 -X main.srcRootDir=$PWD" main.go`
	version = "v0.1.1"
	srcRootDir, _ = os.Getwd()

	// Set options
	zl.SetLogFile("./log/app.jsonl")
	zl.SetVersion(version)
	zl.SetRepositoryCallerEncoder(urlFormat, version, srcRootDir)
	zl.SetLogLevel(zapcore.DebugLevel)
	zl.SetOutputType(zl.OutputTypeShortConsoleAndFile)

	// Initialize
	zl.Init()
	defer zl.Sync()   // flush log buffer
	zl.SyncWhenStop() // flush log buffer. when interrupt or terminated.

	// basic
	zl.Info("USER_INFO", zap.String("name", "Alice"), zap.Int("age", 20))
	err := fmt.Errorf("error message")
	zl.Error("ERROR_MESSAGE", err) // error level log must with error message.
	zl.Debug("DEBUG_MESSAGE")
	zl.Warn("WARN_MESSAGE")
	zl.WarnErr("WARN_MESSAGE_WITH_ERROR", err) // warn level log with error message.
	zl.Info("DISPLAY_TO_CONSOLE", zl.Console("display to console"))
	// Output:
}

func ExampleNew() {
	// Initialize
	traceIDField := "trace_id"
	zl.AddConsoleField(traceIDField)
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
