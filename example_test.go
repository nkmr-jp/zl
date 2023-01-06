package zl_test

import (
	"encoding/json"
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
	// zl.SetOutput(zl.PrettyOutput)
	zl.SetOmitKeys(zl.TimeKey, zl.CallerKey, zl.VersionKey, zl.HostnameKey, zl.StacktraceKey, zl.PIDKey)
	zl.SetRotateFileName(fileName)

	// Initialize
	zl.Init()
	defer zl.Sync() // flush log buffer

	// Write logs
	zl.Info("USER_INFO", zap.String("user_name", "Alice"), zap.Int("user_age", 20)) // can use zap fields.

	console := "display to console when output type is pretty"
	zl.Info("DISPLAY_TO_CONSOLE", zl.Console(console))
	zl.Info("DISPLAY_TO_CONSOLE", zl.Consolep(&console))
	zl.Info("DISPLAY_TO_CONSOLE", zl.Consolef("message: %s", console))

	// write error to error field.
	_, err := os.ReadFile("test")
	zl.Info("READ_FILE_ERROR", zap.Error(err))
	zl.InfoErr("READ_FILE_ERROR", err) // same to above.
	zl.Debug("READ_FILE_ERROR", zap.Error(err))
	zl.DebugErr("READ_FILE_ERROR", err) // same to above.
	zl.Warn("READ_FILE_ERROR", zap.Error(err))
	zl.WarnErr("READ_FILE_ERROR", err) // same to above.
	zl.Error("READ_FILE_ERROR", zap.Error(err))
	zl.ErrorErr("READ_FILE_ERROR", err) // same to above.
	zl.Err("READ_FILE_ERROR", err)      // same to above.
	zl.ErrRet("READ_FILE_ERROR", err)   // write error to log and return same error.

	bytes, _ := os.ReadFile(fileName)
	fmt.Println(string(bytes))

	// Output to stderr with colored:
	// zl.go:80: DEBUG INIT_LOGGER Severity: DEBUG, Output: Pretty, File: ./log/example.jsonl
	// example_test.go:39: INFO USER_INFO
	// example_test.go:42: INFO DISPLAY_TO_CONSOLE display to console when output type is pretty
	// example_test.go:43: INFO DISPLAY_TO_CONSOLE display to console when output type is pretty
	// example_test.go:44: INFO DISPLAY_TO_CONSOLE message: display to console when output type is pretty
	// example_test.go:48: INFO READ_FILE_ERROR
	// example_test.go:49: INFO READ_FILE_ERROR open test: no such file or directory
	// example_test.go:50: DEBUG READ_FILE_ERROR
	// example_test.go:51: DEBUG READ_FILE_ERROR open test: no such file or directory
	// example_test.go:52: WARN READ_FILE_ERROR
	// example_test.go:53: WARN READ_FILE_ERROR open test: no such file or directory
	// example_test.go:54: ERROR READ_FILE_ERROR
	// example_test.go:55: ERROR READ_FILE_ERROR open test: no such file or directory
	// example_test.go:56: ERROR READ_FILE_ERROR open test: no such file or directory
	// example_test.go:57: ERROR READ_FILE_ERROR open test: no such file or directory

	// Output:
	// {"severity":"DEBUG","function":"github.com/nkmr-jp/zl.Init.func1","message":"INIT_LOGGER","console":"Severity: DEBUG, Output: Pretty, File: ./log/example.jsonl"}
	// {"severity":"INFO","function":"github.com/nkmr-jp/zl_test.Example","message":"USER_INFO","user_name":"Alice","user_age":20}
	// {"severity":"INFO","function":"github.com/nkmr-jp/zl_test.Example","message":"DISPLAY_TO_CONSOLE","console":"display to console when output type is pretty"}
	// {"severity":"INFO","function":"github.com/nkmr-jp/zl_test.Example","message":"DISPLAY_TO_CONSOLE","console":"display to console when output type is pretty"}
	// {"severity":"INFO","function":"github.com/nkmr-jp/zl_test.Example","message":"DISPLAY_TO_CONSOLE","console":"message: display to console when output type is pretty"}
	// {"severity":"INFO","function":"github.com/nkmr-jp/zl_test.Example","message":"READ_FILE_ERROR","error":"open test: no such file or directory"}
	// {"severity":"INFO","function":"github.com/nkmr-jp/zl_test.Example","message":"READ_FILE_ERROR","error":"open test: no such file or directory"}
	// {"severity":"DEBUG","function":"github.com/nkmr-jp/zl_test.Example","message":"READ_FILE_ERROR","error":"open test: no such file or directory"}
	// {"severity":"DEBUG","function":"github.com/nkmr-jp/zl_test.Example","message":"READ_FILE_ERROR","error":"open test: no such file or directory"}
	// {"severity":"WARN","function":"github.com/nkmr-jp/zl_test.Example","message":"READ_FILE_ERROR","error":"open test: no such file or directory"}
	// {"severity":"WARN","function":"github.com/nkmr-jp/zl_test.Example","message":"READ_FILE_ERROR","error":"open test: no such file or directory"}
	// {"severity":"ERROR","function":"github.com/nkmr-jp/zl_test.Example","message":"READ_FILE_ERROR","error":"open test: no such file or directory"}
	// {"severity":"ERROR","function":"github.com/nkmr-jp/zl_test.Example","message":"READ_FILE_ERROR","error":"open test: no such file or directory"}
	// {"severity":"ERROR","function":"github.com/nkmr-jp/zl_test.Example","message":"READ_FILE_ERROR","error":"open test: no such file or directory"}
	// {"severity":"ERROR","function":"github.com/nkmr-jp/zl_test.Example","message":"READ_FILE_ERROR","error":"open test: no such file or directory"}
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
	// {"severity":"DEBUG","caller":"zl/zl.go:77","message":"INIT_LOGGER","version":"v1.0.0","console":"Severity: DEBUG, Output: ConsoleAndFile, File: ./log/example-set-version_v1.0.0.jsonl"}
	// {"severity":"INFO","caller":"https://github.com/nkmr-jp/zl/blob/v1.0.0/example_test.go#L121","message":"INFO_MESSAGE","version":"v1.0.0","detail":"detail info xxxxxxxxxxxxxxxxx"}
	// {"severity":"WARN","caller":"https://github.com/nkmr-jp/zl/blob/v1.0.0/example_test.go#L122","message":"WARN_MESSAGE","version":"v1.0.0","detail":"detail info xxxxxxxxxxxxxxxxx"}

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
	err := fmt.Errorf("error")
	zl.Info("GLOBAL_INFO")
	l1.Info("CONTEXT_SCOPE_INFO", zl.Consolef("some message to console: %s", "test"))
	l1.Err("CONTEXT_SCOPE_ERROR", fmt.Errorf("context scope error message"))
	l2.Info("CONTEXT_SCOPE_INFO2", zl.Consolef("some message to console: %s", "test"))
	l2.Debug("TEST")
	l2.Warn("TEST")
	l2.Error("TEST")
	l2.Err("TEST", err)
	l2.ErrorErr("TEST", err)
	l2.ErrRet("TEST", err) // write error to log and return same error.
	l2.InfoErr("TEST", err)
	l2.DebugErr("TEST", err)
	l2.WarnErr("TEST", err)

	bytes, _ := os.ReadFile(fileName)
	fmt.Println(string(bytes))

	// Output to stderr with colored:
	// zl.go:77: DEBUG INIT_LOGGER Severity: DEBUG, Output: Pretty, File: ./log/example-new.jsonl
	// example_test.go:145: INFO GLOBAL_INFO
	// log1 | example_test.go:146: INFO CONTEXT_SCOPE_INFO some message to console: test c7mg6hnr2g4l6vvuao50
	// log1 | example_test.go:147: ERROR CONTEXT_SCOPE_ERROR context scope error message c7mg6hnr2g4l6vvuao50
	// log2 | example_test.go:148: INFO CONTEXT_SCOPE_INFO2 some message to console: test c7mg6hnr2g4l6vvuao50
	// log2 | example_test.go:149: DEBUG TEST c7mg6hnr2g4l6vvuao50
	// log2 | example_test.go:150: WARN TEST c7mg6hnr2g4l6vvuao50
	// log2 | example_test.go:151: ERROR TEST c7mg6hnr2g4l6vvuao50
	// log2 | example_test.go:152: ERROR TEST error c7mg6hnr2g4l6vvuao50
	// log2 | example_test.go:153: ERROR TEST error c7mg6hnr2g4l6vvuao50
	// log2 | example_test.go:154: ERROR TEST error c7mg6hnr2g4l6vvuao50
	// log2 | example_test.go:155: INFO TEST error c7mg6hnr2g4l6vvuao50
	// log2 | example_test.go:156: DEBUG TEST error c7mg6hnr2g4l6vvuao50
	// log2 | example_test.go:157: WARN TEST error c7mg6hnr2g4l6vvuao50

	// Output:
	// {"severity":"DEBUG","message":"INIT_LOGGER","console":"Severity: DEBUG, Output: Pretty, File: ./log/example-new.jsonl"}
	// {"severity":"INFO","message":"GLOBAL_INFO"}
	// {"severity":"INFO","logger":"log1","message":"CONTEXT_SCOPE_INFO","console":"some message to console: test","user_id":1,"trace":"c7mg6hnr2g4l6vvuao50"}
	// {"severity":"ERROR","logger":"log1","message":"CONTEXT_SCOPE_ERROR","error":"context scope error message","user_id":1,"trace":"c7mg6hnr2g4l6vvuao50"}
	// {"severity":"INFO","logger":"log2","message":"CONTEXT_SCOPE_INFO2","console":"some message to console: test","user_id":1,"trace":"c7mg6hnr2g4l6vvuao50"}
	// {"severity":"DEBUG","logger":"log2","message":"TEST","user_id":1,"trace":"c7mg6hnr2g4l6vvuao50"}
	// {"severity":"WARN","logger":"log2","message":"TEST","user_id":1,"trace":"c7mg6hnr2g4l6vvuao50"}
	// {"severity":"WARN","logger":"log2","message":"TEST","user_id":1,"trace":"c7mg6hnr2g4l6vvuao50"}
	// {"severity":"ERROR","logger":"log2","message":"TEST","error":"error","user_id":1,"trace":"c7mg6hnr2g4l6vvuao50"}
	// {"severity":"ERROR","logger":"log2","message":"TEST","error":"error","user_id":1,"trace":"c7mg6hnr2g4l6vvuao50"}
	// {"severity":"ERROR","logger":"log2","message":"TEST","error":"error","user_id":1,"trace":"c7mg6hnr2g4l6vvuao50"}
	// {"severity":"INFO","logger":"log2","message":"TEST","error":"error","user_id":1,"trace":"c7mg6hnr2g4l6vvuao50"}
	// {"severity":"DEBUG","logger":"log2","message":"TEST","error":"error","user_id":1,"trace":"c7mg6hnr2g4l6vvuao50"}
	// {"severity":"WARN","logger":"log2","message":"TEST","error":"error","user_id":1,"trace":"c7mg6hnr2g4l6vvuao50"}
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

	_, err := os.ReadFile("test")
	zl.Err("READ_FILE_ERROR", err)
	zl.Info("INFO")
	zl.InfoErr("INFO_ERR", fmt.Errorf("error"))
	v := ""
	err = json.Unmarshal([]byte("test"), &v)
	zl.Err("JSON_UNMARSHAL_ERROR", err)

	for i := 0; i < 3; i++ {
		err = fmt.Errorf("if the same error occurs multiple times in the same location, the error report will show them all together")
		zl.Err("ERRORS_IN_LOOPS", err)
	}
	// Output:
}

func ExampleDump() {
	zl.Cleanup() // removes logger and resets settings.
	zl.SetLevel(zl.DebugLevel)
	zl.SetRotateFileName("./log/example-Dump.jsonl")
	zl.Init()
	defer zl.Sync() // flush log buffer
	zl.Dump("test")
	// Output:
}

func ExampleSyncWhenStop() {
	zl.Cleanup() // removes logger and resets settings.
	zl.SetLevel(zl.DebugLevel)
	zl.SetRotateFileName("./log/example-Dump.jsonl")
	zl.Init()
	defer zl.Sync()
	zl.SyncWhenStop()
	zl.Info("TEST")
	// Output:
}
