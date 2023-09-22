package zl_test

import (
	"encoding/json"
	"fmt"
	"github.com/nkmr-jp/zl"
	"go.uber.org/zap"
	"log"
	"os"
	"syscall"
	"time"
)

var (
	version    string // version git revision or tag. set from go cli.
	srcRootDir string // srcRootDir set from cli.
)

func setupForExampleTest() {
	if err := os.RemoveAll("./log"); err != nil {
		log.Fatal(err)
	}
	zl.ResetGlobalLoggerSettings() // removes logger and resets settings.
	zl.SetIsTest()
}

func Example() {
	setupForExampleTest()

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
	zl.Fatal("READ_FILE_ERROR", zap.Error(err))
	zl.FatalErr("READ_FILE_ERROR", err) // same to above.

	fmt.Println("\nlog file output:")
	bytes, _ := os.ReadFile(fileName)
	fmt.Println(string(bytes))

	// Output to stderr with colored:
	// zl.go:82: DEBUG INIT_LOGGER:Severity: DEBUG, Output: Pretty, File: ./log/example.jsonl
	// example_test.go:40: INFO USER_INFO
	// example_test.go:43: INFO DISPLAY_TO_CONSOLE:display to console when output type is pretty
	// example_test.go:44: INFO DISPLAY_TO_CONSOLE:display to console when output type is pretty
	// example_test.go:45: INFO DISPLAY_TO_CONSOLE:message: display to console when output type is pretty
	// example_test.go:49: INFO READ_FILE_ERROR
	// example_test.go:50: INFO READ_FILE_ERROR:open test: no such file or directory
	// example_test.go:51: DEBUG READ_FILE_ERROR
	// example_test.go:52: DEBUG READ_FILE_ERROR:open test: no such file or directory
	// example_test.go:53: WARN READ_FILE_ERROR
	// example_test.go:54: WARN READ_FILE_ERROR:open test: no such file or directory
	// example_test.go:55: ERROR READ_FILE_ERROR
	// example_test.go:56: ERROR READ_FILE_ERROR:open test: no such file or directory
	// example_test.go:57: ERROR READ_FILE_ERROR:open test: no such file or directory
	// example_test.go:58: ERROR READ_FILE_ERROR:open test: no such file or directory
	// example_test.go:59: FATAL READ_FILE_ERROR
	// example_test.go:60: FATAL READ_FILE_ERROR:open test: no such file or directory

	// Output:
	// os.Exit(1) called.
	// os.Exit(1) called.
	//
	// log file output:
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
	// {"severity":"FATAL","function":"github.com/nkmr-jp/zl_test.Example","message":"READ_FILE_ERROR","error":"open test: no such file or directory"}
	// {"severity":"FATAL","function":"github.com/nkmr-jp/zl_test.Example","message":"READ_FILE_ERROR","error":"open test: no such file or directory"}
}

func ExampleSetVersion() {
	setupForExampleTest()

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
	// {"severity":"DEBUG","caller":"zl/zl.go:82","message":"INIT_LOGGER","version":"v1.0.0","console":"Severity: DEBUG, Output: ConsoleAndFile, File: ./log/example-set-version_v1.0.0.jsonl"}
	// {"severity":"INFO","caller":"https://github.com/nkmr-jp/zl/blob/v1.0.0/example_test.go#L135","message":"INFO_MESSAGE","version":"v1.0.0","detail":"detail info xxxxxxxxxxxxxxxxx"}
	// {"severity":"WARN","caller":"https://github.com/nkmr-jp/zl/blob/v1.0.0/example_test.go#L136","message":"WARN_MESSAGE","version":"v1.0.0","detail":"detail info xxxxxxxxxxxxxxxxx"}

}

func ExampleNew() {
	setupForExampleTest()

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
	log := zl.New(
		zap.Int("user_id", 1),
		zap.String(traceIDField, traceID),
	)
	emptyName := log.Named("")
	named1 := log.Named("named1")
	named2 := log.Named("named2")
	named3 := named1.Named("named3") // named1.named3

	// Write logs
	err := fmt.Errorf("error")
	zl.Info("GLOBAL_INFO")
	log.Info("CONTEXT_SCOPE_INFO", zl.Consolef("some message to console: %s", "test"))
	emptyName.Err("CONTEXT_SCOPE_ERROR", fmt.Errorf("context scope error message"))
	named1.Info("CONTEXT_SCOPE_INFO2", zl.Consolef("some message to console: %s", "test"))
	named2.Debug("TEST")
	named3.Warn("TEST")
	log.Error("TEST")
	named1.Err("TEST", err)
	named2.ErrorErr("TEST", err)
	named3.ErrRet("TEST", err) // write error to log and return same error.
	log.InfoErr("TEST", err)
	named1.DebugErr("TEST", err)
	named2.WarnErr("TEST", err)
	named3.Fatal("TEST")
	log.FatalErr("TEST", err)

	fmt.Println("\nlog file output:")
	bytes, _ := os.ReadFile(fileName)
	fmt.Println(string(bytes))

	// Output to stderr with colored:
	// zl.go:82: DEBUG INIT_LOGGER Severity: DEBUG, Output: Pretty, File: ./log/example-new.jsonl
	// example_test.go:176: INFO GLOBAL_INFO
	// example_test.go:177: INFO CONTEXT_SCOPE_INFO some message to console: test c7mg6hnr2g4l6vvuao50
	// example_test.go:178: ERROR CONTEXT_SCOPE_ERROR context scope error message c7mg6hnr2g4l6vvuao50
	// named1 | example_test.go:179: INFO CONTEXT_SCOPE_INFO2 some message to console: test c7mg6hnr2g4l6vvuao50
	// named2 | example_test.go:180: DEBUG TEST c7mg6hnr2g4l6vvuao50
	// named1.named3 | example_test.go:181: WARN TEST c7mg6hnr2g4l6vvuao50
	// example_test.go:182: ERROR TEST c7mg6hnr2g4l6vvuao50
	// named1 | example_test.go:183: ERROR TEST error c7mg6hnr2g4l6vvuao50
	// named2 | example_test.go:184: ERROR TEST error c7mg6hnr2g4l6vvuao50
	// named1.named3 | example_test.go:185: ERROR TEST error c7mg6hnr2g4l6vvuao50
	// example_test.go:186: INFO TEST error c7mg6hnr2g4l6vvuao50
	// named1 | example_test.go:187: DEBUG TEST error c7mg6hnr2g4l6vvuao50
	// named2 | example_test.go:188: WARN TEST error c7mg6hnr2g4l6vvuao50
	// named1.named3 | example_test.go:189: FATAL TEST c7mg6hnr2g4l6vvuao50
	// example_test.go:190: FATAL TEST error c7mg6hnr2g4l6vvuao50

	// Output:
	// os.Exit(1) called.
	// os.Exit(1) called.
	//
	// log file output:
	// {"severity":"DEBUG","message":"INIT_LOGGER","console":"Severity: DEBUG, Output: Pretty, File: ./log/example-new.jsonl"}
	// {"severity":"INFO","message":"GLOBAL_INFO"}
	// {"severity":"INFO","message":"CONTEXT_SCOPE_INFO","console":"some message to console: test","user_id":1,"trace":"c7mg6hnr2g4l6vvuao50"}
	// {"severity":"ERROR","message":"CONTEXT_SCOPE_ERROR","error":"context scope error message","user_id":1,"trace":"c7mg6hnr2g4l6vvuao50"}
	// {"severity":"INFO","logger":"named1","message":"CONTEXT_SCOPE_INFO2","console":"some message to console: test","user_id":1,"trace":"c7mg6hnr2g4l6vvuao50"}
	// {"severity":"DEBUG","logger":"named2","message":"TEST","user_id":1,"trace":"c7mg6hnr2g4l6vvuao50"}
	// {"severity":"WARN","logger":"named1.named3","message":"TEST","user_id":1,"trace":"c7mg6hnr2g4l6vvuao50"}
	// {"severity":"ERROR","message":"TEST","user_id":1,"trace":"c7mg6hnr2g4l6vvuao50"}
	// {"severity":"ERROR","logger":"named1","message":"TEST","error":"error","user_id":1,"trace":"c7mg6hnr2g4l6vvuao50"}
	// {"severity":"ERROR","logger":"named2","message":"TEST","error":"error","user_id":1,"trace":"c7mg6hnr2g4l6vvuao50"}
	// {"severity":"ERROR","logger":"named1.named3","message":"TEST","error":"error","user_id":1,"trace":"c7mg6hnr2g4l6vvuao50"}
	// {"severity":"INFO","message":"TEST","error":"error","user_id":1,"trace":"c7mg6hnr2g4l6vvuao50"}
	// {"severity":"DEBUG","logger":"named1","message":"TEST","error":"error","user_id":1,"trace":"c7mg6hnr2g4l6vvuao50"}
	// {"severity":"WARN","logger":"named2","message":"TEST","error":"error","user_id":1,"trace":"c7mg6hnr2g4l6vvuao50"}
	// {"severity":"FATAL","logger":"named1.named3","message":"TEST","user_id":1,"trace":"c7mg6hnr2g4l6vvuao50"}
	// {"severity":"FATAL","message":"TEST","error":"error","user_id":1,"trace":"c7mg6hnr2g4l6vvuao50"}
}

func ExampleSetLevelByString() {
	setupForExampleTest()

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

func ExampleSetOmitKeys() {
	setupForExampleTest()

	zl.SetOutputByString("Console")
	zl.SetStdout()
	zl.SetOmitKeys(
		zl.MessageKey, zl.LevelKey, zl.LoggerKey, zl.TimeKey,
		zl.CallerKey, zl.VersionKey, zl.HostnameKey, zl.StacktraceKey, zl.PIDKey,
	)
	zl.Init()
	zl.Info("INFO_MESSAGE")

	// Output:
	// {"function":"github.com/nkmr-jp/zl_test.ExampleSetOmitKeys"}
}

func ExampleError() {
	setupForExampleTest()

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
	setupForExampleTest()

	zl.SetLevel(zl.DebugLevel)
	zl.SetRotateFileName("./log/example-Dump.jsonl")
	zl.Init()
	defer zl.Sync() // flush log buffer
	zl.Dump("test")
	// Output:
}

func ExampleSyncWhenStop() {
	// syscall.SIGINT
	setupForExampleTest()
	zl.SetLevel(zl.DebugLevel)
	zl.SetRotateFileName("./log/example-SyncWhenStop.jsonl")
	zl.Init()
	zl.SyncWhenStop()

	go func() {
		time.Sleep(time.Millisecond * 50)
		syscall.Kill(os.Getpid(), syscall.SIGINT)
	}()
	time.Sleep(time.Millisecond * 100)

	// syscall.SIGTERM
	fmt.Println()
	setupForExampleTest()
	zl.SetLevel(zl.DebugLevel)
	zl.SetRotateFileName("./log/example-SyncWhenStop.jsonl")
	zl.Init()
	zl.SyncWhenStop()

	go func() {
		time.Sleep(time.Millisecond * 50)
		syscall.Kill(os.Getpid(), syscall.SIGTERM)
	}()
	time.Sleep(time.Millisecond * 100)

	// Output:
	// os.Exit(130) called.
	// os.Exit(143) called.
}
