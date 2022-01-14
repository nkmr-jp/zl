package zl_test

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/nkmr-jp/zap-lightning/zl"
	"go.uber.org/zap"
)

var (
	version    string // version git revision or tag. set from go cli.
	srcRootDir string // srcRootDir set from cli.
)

func TestMain(m *testing.M) {
	if err := os.RemoveAll("./log"); err != nil {
		log.Fatal(err)
	}
	m.Run()
}

func Example() {
	// Set Options
	zl.SetLevel(zl.DebugLevel)                                     // default is InfoLevel.
	zl.SetOutput(zl.PrettyOutput)                                  // PrettyOutput is default. recommended for develop environment.
	zl.SetIgnoreKeys(zl.TimeKey, zl.HostnameKey, zl.StacktraceKey) // ignore fields for test.

	// Initialize
	zl.Init()
	defer zl.Sync()   // flush log buffer
	zl.SyncWhenStop() // flush log buffer. when interrupt or terminated.

	// Write Logs
	fmt.Println("Console:")
	zl.Info("USER_INFO", zap.String("user_name", "Alice"), zap.Int("user_age", 20)) // can use zap fields.
	err := fmt.Errorf("error message")
	zl.Error("ERROR_MESSAGE", err) // error level log must with error message.
	zl.Debug("DEBUG_MESSAGE")
	zl.Warn("WARN_MESSAGE", zap.Error(err))    // warn level log with error message.
	zl.WarnErr("WARN_MESSAGE_WITH_ERROR", err) // same to above.
	zl.Info("DISPLAY_TO_CONSOLE", zl.Console("display to console when output type is pretty"))

	fmt.Println("\nFile:")
	app, _ := os.ReadFile("./log/app.jsonl")
	fmt.Println(string(app))

	// Output:
	// Console:
	//
	// File:
	// {"level":"INFO","caller":"zl/zl.go:39","function":"github.com/nkmr-jp/zap-lightning/zl.Init.func1","message":"INIT_LOGGER","version":"7301145","console":"Level: DEBUG, Output: Pretty, FileName: ./log/app.jsonl"}
	// {"level":"INFO","caller":"zl/example_test.go:39","function":"github.com/nkmr-jp/zap-lightning/zl_test.Example","message":"USER_INFO","version":"7301145","user_name":"Alice","user_age":20}
	// {"level":"ERROR","caller":"zl/example_test.go:41","function":"github.com/nkmr-jp/zap-lightning/zl_test.Example","message":"ERROR_MESSAGE","version":"7301145","error":"error message"}
	// {"level":"DEBUG","caller":"zl/example_test.go:42","function":"github.com/nkmr-jp/zap-lightning/zl_test.Example","message":"DEBUG_MESSAGE","version":"7301145"}
	// {"level":"WARN","caller":"zl/example_test.go:43","function":"github.com/nkmr-jp/zap-lightning/zl_test.Example","message":"WARN_MESSAGE","version":"7301145","error":"error message"}
	// {"level":"WARN","caller":"zl/example_test.go:44","function":"github.com/nkmr-jp/zap-lightning/zl_test.Example","message":"WARN_MESSAGE_WITH_ERROR","version":"7301145","error":"error message"}
	// {"level":"INFO","caller":"zl/example_test.go:45","function":"github.com/nkmr-jp/zap-lightning/zl_test.Example","message":"DISPLAY_TO_CONSOLE","version":"7301145","console":"display to console when output type is pretty"}
}

func ExampleSetVersion() {
	// ex.`go run -ldflags "-X main.version=v0.1.1 -X main.srcRootDir=$PWD" main.go`
	version = "v0.1.1"
	srcRootDir, _ = os.Getwd()

	// Set Options
	zl.SetLevel(zl.DebugLevel) // Default is InfoLevel
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

func ExampleSetOutput() {
	// Set options
	zl.SetLevel(zl.DebugLevel) // Default is InfoLevel
	zl.SetOutput(zl.ConsoleOutput)

	// Initialize
	zl.Init()
	defer zl.Sync()   // flush log buffer
	zl.SyncWhenStop() // flush log buffer. when interrupt or terminated.

	// Logs
	zl.Info("USER_INFO", zl.Console("message"))

	// Example output:
	// {"level":"INFO","ts":"2022-01-14T06:43:03.345143+09:00","caller":"zl/zl.go:36","function":"github.com/nkmr-jp/zap-lightning/zl.Init.func1","msg":"INIT_LOGGER","version":"dd90b59","hostname":"nkmrnoMacBook-Pro.local","console":"logLevel: DEBUG, fileName: , outputType: Console"}
	// {"level":"INFO","ts":"2022-01-14T06:43:03.345329+09:00","caller":"zap-lightning/example_test.go:74","function":"github.com/nkmr-jp/zap-lightning_test.ExampleSetOutputType","msg":"USER_INFO","version":"dd90b59","hostname":"nkmrnoMacBook-Pro.local","console":"message"}
	// {"level":"INFO","ts":"2022-01-14T06:43:03.345341+09:00","caller":"zl/logger.go:67","function":"github.com/nkmr-jp/zap-lightning/zl.Sync","msg":"FLUSH_LOG_BUFFER","version":"dd90b59","hostname":"nkmrnoMacBook-Pro.local"}
}

func ExampleNew() {
	// Set options
	traceIDField := "trace_id"
	zl.AddConsoleFields(traceIDField)

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
