package zl_test

import (
	"fmt"
	"log"
	"os"
	"testing"

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
	// Set options
	zl.SetLevel(zl.DebugLevel)
	zl.SetOutput(zl.PrettyOutput)
	zl.SetIgnoreKeys(zl.TimeKey, zl.VersionKey, zl.HostnameKey, zl.StacktraceKey)

	// Initialize
	zl.Init()
	defer zl.Sync()   // flush log buffer
	zl.SyncWhenStop() // flush log buffer. when interrupt or terminated.

	// Write logs
	zl.Info("USER_INFO", zap.String("user_name", "Alice"), zap.Int("user_age", 20)) // can use zap fields.
	err := fmt.Errorf("error message")
	zl.Error("ERROR_MESSAGE", err) // error level log must with error message.
	zl.Debug("DEBUG_MESSAGE")
	zl.Warn("WARN_MESSAGE", zap.Error(err))    // warn level log with error message.
	zl.WarnErr("WARN_MESSAGE_WITH_ERROR", err) // same to above.
	zl.Info("DISPLAY_TO_CONSOLE", zl.Console("display to console when output type is pretty"))

	bytes, _ := os.ReadFile("./log/app.jsonl")
	fmt.Println(string(bytes))

	// Output to stderr with colored:
	// zl.go:41: DEBUG INIT_LOGGER: Level: DEBUG, Output: Pretty, FileName: ./log/app.jsonl
	// example_test.go:37: INFO USER_INFO
	// example_test.go:39: ERROR ERROR_MESSAGE: error message
	// example_test.go:40: DEBUG DEBUG_MESSAGE
	// example_test.go:41: WARN WARN_MESSAGE
	// example_test.go:42: WARN WARN_MESSAGE_WITH_ERROR: error message
	// example_test.go:43: INFO DISPLAY_TO_CONSOLE: display to console when output type is pretty
	// zl.go:140: DEBUG FLUSH_LOG_BUFFER

	// Output:
	// {"level":"DEBUG","caller":"zl/zl.go:41","function":"github.com/nkmr-jp/zap-lightning/zl.Init.func1","message":"INIT_LOGGER","console":"Level: DEBUG, Output: Pretty, FileName: ./log/app.jsonl"}
	// {"level":"INFO","caller":"zl/example_test.go:37","function":"github.com/nkmr-jp/zap-lightning/zl_test.Example","message":"USER_INFO","user_name":"Alice","user_age":20}
	// {"level":"ERROR","caller":"zl/example_test.go:39","function":"github.com/nkmr-jp/zap-lightning/zl_test.Example","message":"ERROR_MESSAGE","error":"error message"}
	// {"level":"DEBUG","caller":"zl/example_test.go:40","function":"github.com/nkmr-jp/zap-lightning/zl_test.Example","message":"DEBUG_MESSAGE"}
	// {"level":"WARN","caller":"zl/example_test.go:41","function":"github.com/nkmr-jp/zap-lightning/zl_test.Example","message":"WARN_MESSAGE","error":"error message"}
	// {"level":"WARN","caller":"zl/example_test.go:42","function":"github.com/nkmr-jp/zap-lightning/zl_test.Example","message":"WARN_MESSAGE_WITH_ERROR","error":"error message"}
	// {"level":"INFO","caller":"zl/example_test.go:43","function":"github.com/nkmr-jp/zap-lightning/zl_test.Example","message":"DISPLAY_TO_CONSOLE","console":"display to console when output type is pretty"}
}

func ExampleSetVersion() {
	zl.Cleanup() // removes logger and resets settings.

	urlFormat := "https://github.com/nkmr-jp/zap-lightning/blob/%s"

	// Actually, it is recommended to pass the value from the command line of go.
	// ex. `go run -ldflags "-X main.version=v0.1.1 -X main.srcRootDir=$PWD" main.go`.
	version = "v1.0.0"
	srcRootDir, _ = os.Getwd()

	// Set Options
	zl.SetVersion(version)
	zl.SetFileName(fmt.Sprintf("./log/app_%s.jsonl", zl.GetVersion()))
	zl.SetRepositoryCallerEncoder(urlFormat, version, srcRootDir)
	zl.SetIgnoreKeys(zl.TimeKey, zl.FunctionKey, zl.HostnameKey)
	zl.SetOutput(zl.ConsoleOutput)
	zl.SetStdout()

	// Initialize
	zl.Init()
	defer zl.Sync()   // flush log buffer
	zl.SyncWhenStop() // flush log buffer. when interrupt or terminated.

	// Write logs
	zl.Warn("WARN_MESSAGE", zap.String("detail", "detail info xxxxxxxxxxxxxxxxx"))

	// Output:
	// {"level":"WARN","caller":"https://github.com/nkmr-jp/zap-lightning/blob/v1.0.0/example_test.go#L92","message":"WARN_MESSAGE","version":"v1.0.0","detail":"detail info xxxxxxxxxxxxxxxxx"}
}

func ExampleNew() {
	zl.Cleanup() // removes logger and resets settings.

	// Set options
	traceIDField := "trace_id"
	zl.AddConsoleFields(traceIDField)
	zl.SetIgnoreKeys(zl.TimeKey, zl.FunctionKey, zl.VersionKey, zl.HostnameKey, zl.StacktraceKey)
	zl.SetOutput(zl.PrettyOutput)

	// New
	// ex. Use this when you want to add a common value in the scope of a context, such as an API request.
	l1 := zl.New(
		zap.Int("user_id", 1),
		zap.Int64(traceIDField, 1642153670000264000),
	).Named("log1")

	l2 := zl.New(
		zap.Int("user_id", 1),
		zap.Int64(traceIDField, 1642153670000264000),
	).Named("log2")

	// Write logs
	l1.Info("CONTEXT_SCOPE_INFO", zl.Consolef("some message to console: %s", "test"))
	l1.Error("CONTEXT_SCOPE_ERROR", fmt.Errorf("context scope error message"))
	l2.Info("CONTEXT_SCOPE_INFO_LOG2", zl.Consolef("some message to console: %s", "test"))

	bytes, _ := os.ReadFile("./log/app.jsonl")
	fmt.Println(string(bytes))

	// Output to stderr with colored:
	// example_test.go:117: INFO CONTEXT_SCOPE_INFO: some message to console: test, 1642153670000264000
	// example_test.go:118: ERROR CONTEXT_SCOPE_ERROR: context scope error message : 1642153670000264000

	// Output:
	// {"level":"INFO","name":"log1","caller":"zl/example_test.go:120","message":"CONTEXT_SCOPE_INFO","console":"some message to console: test","user_id":1,"trace_id":1642153670000264000}
	// {"level":"ERROR","name":"log1","caller":"zl/example_test.go:121","message":"CONTEXT_SCOPE_ERROR","error":"context scope error message","user_id":1,"trace_id":1642153670000264000}
	// {"level":"INFO","name":"log2","caller":"zl/example_test.go:122","message":"CONTEXT_SCOPE_INFO_LOG2","console":"some message to console: test","user_id":1,"trace_id":1642153670000264000}
}
