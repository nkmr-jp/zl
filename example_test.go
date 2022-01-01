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
	consoleField = "console"
	traceIDField = "trace_id"
	urlFormat    = "https://github.com/nkmr-jp/zap-lightning/blob/%s"
)

func Example() {
	// If you use this, you might actually want to pass the value from the go command, like this
	// ex.`go run -ldflags "-X main.version=v1.0.0 -X main.srcRootDir=$PWD" main.go`
	version = "v0.1.1"
	srcRootDir, _ = os.Getwd()

	// Set options
	zl.SetLogFile("./log/app_%Y-%m-%d.jsonl")
	zl.SetVersion(version)
	zl.SetRepositoryCallerEncoder(urlFormat, version, srcRootDir)
	zl.SetConsoleField(consoleField, traceIDField)
	zl.SetLogLevel(zapcore.DebugLevel)
	zl.SetOutputType(zl.OutputTypeShortConsoleAndFile)

	// Initialize
	zl.InitLogger()
	defer zl.Sync()   // flush log buffer
	zl.SyncWhenStop() // flush log buffer. when interrupt or terminated.

	// basic
	zl.Info("USER_INFO", zap.String("name", "Alice"), zap.Int("age", 20))
	zl.Errorf("SOME_ERROR", fmt.Errorf("error message"))
	zl.Debug("DEBUG_MESSAGE")
	zl.Warn("WARN_MESSAGE")
	// display to console log
	zl.Info("DISPLAY_TO_CONSOLE", zap.String(consoleField, "display to console"))
	// Output:
}

func ExampleNewWrapper() {
	// Initialize
	zl.InitLogger()
	defer zl.Sync()
	zl.SyncWhenStop()

	// NewWrapper
	// ex. Use this when you want to add a common value in the scope of a context, such as an API request.
	w := zl.NewWrapper(
		zap.Int("user_id", 1),
		zap.Int64(traceIDField, time.Now().UnixNano()),
	)
	w.Info("CONTEXT_SCOPE_INFO")
	w.Errorf("CONTEXT_SCOPE_ERROR", fmt.Errorf("context scope error message"))
	// Output:
}
