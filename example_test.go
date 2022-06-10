package zl_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/nkmr-jp/zl"
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
	fileName := "./log/example.jsonl"
	zl.SetLevel(zl.DebugLevel)
	zl.SetOutput(zl.PrettyOutput)
	zl.SetOmitKeys(zl.TimeKey, zl.CallerKey, zl.VersionKey, zl.HostnameKey, zl.StacktraceKey, zl.PIDKey)
	zl.SetRotateFileName(fileName)

	// Initialize
	zl.Init()
	defer zl.Sync() // flush log buffer

	// Write logs
	zl.Info("USER_INFO", zap.String("user_name", "Alice"), zap.Int("user_age", 20)) // can use zap fields.
	err := fmt.Errorf("error message")
	zl.Error("ERROR_MESSAGE", err) // error level log must with error message.
	zl.Debug("DEBUG_MESSAGE")
	zl.Warn("WARN_MESSAGE", zap.Error(err))    // warn level log with error message.
	zl.WarnErr("WARN_MESSAGE_WITH_ERROR", err) // same to above.
	zl.Info("DISPLAY_TO_CONSOLE", zl.Console("display to console when output type is pretty"))
	zl.DebugErr("DEBUG_MESSAGE_WITH_ERROR_AND_CONSOLE", err, zl.Console("display to console when output type is pretty")) // same to above.

	bytes, _ := os.ReadFile(fileName)
	fmt.Println(string(bytes))

	// Output to stderr with colored:
	// zl.go:44: DEBUG INIT_LOGGER: Severity: DEBUG, Output: Pretty, FileName: ./log/example.jsonl
	// example_test.go:39: INFO USER_INFO
	// example_test.go:41: ERROR ERROR_MESSAGE: error message
	// example_test.go:42: DEBUG DEBUG_MESSAGE
	// example_test.go:43: WARN WARN_MESSAGE
	// example_test.go:44: WARN WARN_MESSAGE_WITH_ERROR: error message
	// example_test.go:45: INFO DISPLAY_TO_CONSOLE: display to console when output type is pretty
	// example_test.go:46: DEBUG DEBUG_MESSAGE_WITH_ERROR_AND_CONSOLE: error message , display to console when output type is pretty
	// zl.go:131: DEBUG FLUSH_LOG_BUFFER

	// Output:
	// {"severity":"DEBUG","function":"github.com/nkmr-jp/zl.Init.func1","message":"INIT_LOGGER","console":"Severity: DEBUG, Output: Pretty, File: ./log/example.jsonl"}
	// {"severity":"INFO","function":"github.com/nkmr-jp/zl_test.Example","message":"USER_INFO","user_name":"Alice","user_age":20}
	// {"severity":"ERROR","function":"github.com/nkmr-jp/zl_test.Example","message":"ERROR_MESSAGE","error":"error message"}
	// {"severity":"DEBUG","function":"github.com/nkmr-jp/zl_test.Example","message":"DEBUG_MESSAGE"}
	// {"severity":"WARN","function":"github.com/nkmr-jp/zl_test.Example","message":"WARN_MESSAGE","error":"error message"}
	// {"severity":"WARN","function":"github.com/nkmr-jp/zl_test.Example","message":"WARN_MESSAGE_WITH_ERROR","error":"error message"}
	// {"severity":"INFO","function":"github.com/nkmr-jp/zl_test.Example","message":"DISPLAY_TO_CONSOLE","console":"display to console when output type is pretty"}
	// {"severity":"DEBUG","function":"github.com/nkmr-jp/zl_test.Example","message":"DEBUG_MESSAGE_WITH_ERROR_AND_CONSOLE","console":"display to console when output type is pretty","error":"error message"}
}

func ExampleSetVersion() {
	zl.Cleanup() // removes logger and resets settings.

	urlFormat := "https://github.com/nkmr-jp/zl/blob/%s"

	// Actually, it is recommended to pass the value from the command line of go.
	// ex. `go run -ldflags "-X main.version=v1.0.0 -X main.srcRootDir=$PWD" main.go`.
	version = "v1.0.0"
	srcRootDir, _ = os.Getwd()

	// Set Options
	zl.SetLevel(zl.DebugLevel)
	zl.SetVersion(version)
	fileName := fmt.Sprintf("./log/example-set-version_%s.jsonl", zl.GetVersion())
	zl.SetRotateFileName(fileName)
	zl.SetRepositoryCallerEncoder(urlFormat, version, srcRootDir)
	zl.SetOmitKeys(zl.TimeKey, zl.FunctionKey, zl.HostnameKey, zl.PIDKey)
	zl.SetOutput(zl.ConsoleAndFileOutput)

	// Initialize
	zl.Init()
	defer zl.Sync() // flush log buffer

	// Write logs
	zl.Info("INFO_MESSAGE", zap.String("detail", "detail info xxxxxxxxxxxxxxxxx"))
	zl.Warn("WARN_MESSAGE", zap.String("detail", "detail info xxxxxxxxxxxxxxxxx"))

	bytes, _ := os.ReadFile(fileName)
	fmt.Println(string(bytes))

	// Output:
	// {"severity":"DEBUG","caller":"zl/zl.go:69","message":"INIT_LOGGER","version":"v1.0.0","console":"Severity: DEBUG, Output: ConsoleAndFile, File: ./log/example-set-version_v1.0.0.jsonl"}
	// {"severity":"INFO","caller":"https://github.com/nkmr-jp/zl/blob/v1.0.0/example_test.go#L96","message":"INFO_MESSAGE","version":"v1.0.0","detail":"detail info xxxxxxxxxxxxxxxxx"}
	// {"severity":"WARN","caller":"https://github.com/nkmr-jp/zl/blob/v1.0.0/example_test.go#L97","message":"WARN_MESSAGE","version":"v1.0.0","detail":"detail info xxxxxxxxxxxxxxxxx"}
}

func ExampleNew() {
	zl.Cleanup() // removes logger and resets settings.

	// Set options
	traceIDField := "trace"
	fileName := "./log/example-new.jsonl"
	zl.SetConsoleFields(traceIDField)
	zl.SetLevel(zl.DebugLevel)
	zl.SetOmitKeys(zl.TimeKey, zl.CallerKey, zl.FunctionKey, zl.VersionKey, zl.HostnameKey, zl.StacktraceKey, zl.PIDKey)
	zl.SetOutput(zl.PrettyOutput)
	zl.SetRotateFileName(fileName)
	traceID := "c7mg6hnr2g4l6vvuao50" // xid.New().String()

	// Initialize
	zl.Init()
	defer zl.Sync() // flush log buffer

	// New
	// e.g. Use this when you want to add a common value in the scope of a context, such as an API request.
	l1 := zl.New(
		zap.Int("user_id", 1),
		zap.String(traceIDField, traceID),
	).Named("log1")

	l2 := zl.New(
		zap.Int("user_id", 1),
		zap.String(traceIDField, traceID),
	).Named("log2")

	// Write logs
	zl.Info("GLOBAL_INFO")
	l1.Info("CONTEXT_SCOPE_INFO", zl.Consolef("some message to console: %s", "test"))
	l1.Error("CONTEXT_SCOPE_ERROR", fmt.Errorf("context scope error message"))
	l2.Info("CONTEXT_SCOPE_INFO2", zl.Consolef("some message to console: %s", "test"))

	bytes, _ := os.ReadFile(fileName)
	fmt.Println(string(bytes))

	// Output to stderr with colored:
	// zl.go:44: DEBUG INIT_LOGGER: Severity: DEBUG, Output: Pretty, FileName: ./log/example-new.jsonl
	// example_test.go:135: INFO GLOBAL_INFO
	// log1 | example_test.go:136: INFO CONTEXT_SCOPE_INFO: some message to console: test, 1642153670000264000
	// log1 | example_test.go:137: ERROR CONTEXT_SCOPE_ERROR: context scope error message : 1642153670000264000
	// log2 | example_test.go:138: INFO CONTEXT_SCOPE_INFO2: some message to console: test, 1642153670000264000
	// zl.go:131: DEBUG FLUSH_LOG_BUFFER

	// Output:
	// {"severity":"DEBUG","message":"INIT_LOGGER","console":"Severity: DEBUG, Output: Pretty, File: ./log/example-new.jsonl"}
	// {"severity":"INFO","message":"GLOBAL_INFO"}
	// {"severity":"INFO","logger":"log1","message":"CONTEXT_SCOPE_INFO","console":"some message to console: test","user_id":1,"trace":"c7mg6hnr2g4l6vvuao50"}
	// {"severity":"ERROR","logger":"log1","message":"CONTEXT_SCOPE_ERROR","error":"context scope error message","user_id":1,"trace":"c7mg6hnr2g4l6vvuao50"}
	// {"severity":"INFO","logger":"log2","message":"CONTEXT_SCOPE_INFO2","console":"some message to console: test","user_id":1,"trace":"c7mg6hnr2g4l6vvuao50"}
}

func ExampleSetLevelByString() {
	zl.Cleanup() // removes logger and resets settings.

	zl.SetLevelByString("DEBUG")
	zl.SetOutputByString("Console")
	zl.SetStdout()
	zl.SetOmitKeys(zl.TimeKey, zl.CallerKey, zl.FunctionKey, zl.VersionKey, zl.HostnameKey, zl.StacktraceKey, zl.PIDKey)

	zl.Init()
	zl.Debug("DEBUG_MESSAGE")
	zl.Info("INFO_MESSAGE")

	// Output:
	// {"severity":"DEBUG","message":"INIT_LOGGER","console":"Severity: DEBUG, Output: Console"}
	// {"severity":"DEBUG","message":"DEBUG_MESSAGE"}
	// {"severity":"INFO","message":"INFO_MESSAGE"}
}

func ExampleError() {
	zl.Cleanup() // removes logger and resets settings.

	zl.SetOmitKeys(zl.TimeKey, zl.VersionKey, zl.HostnameKey)

	// Initialize
	zl.Init()
	defer zl.Sync() // flush log buffer

	zl.Error("ERROR_WITH_STACKTRACE", fmt.Errorf("error occurred"))
	zl.Info("INFO")
	zl.Error("ERROR_WITH_STACKTRACE", fmt.Errorf("error occurred"))
	// Output:
}
