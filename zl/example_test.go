package zl_test

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/nkmr-jp/zap-lightning/zl"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
	zl.SetLevel(zapcore.DebugLevel) // Default is InfoLevel
	zl.SetOutput(zl.PrettyOutput)   // Default. it's recommended for develop environment.

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

	// Example output:
	// Colored Simple Log
	//
	// 2022/01/14 12:01:52 zl.go:37: INFO INIT_LOGGER: logLevel: DEBUG, fileName: ./log/app.jsonl, outputType: Pretty
	// 2022/01/14 12:01:52 example_test.go:28: INFO USER_INFO
	// 2022/01/14 12:01:52 example_test.go:30: ERROR ERROR_MESSAGE: error message
	// 2022/01/14 12:01:52 example_test.go:31: DEBUG DEBUG_MESSAGE
	// 2022/01/14 12:01:52 example_test.go:32: WARN WARN_MESSAGE
	// 2022/01/14 12:01:52 example_test.go:33: WARN WARN_MESSAGE_WITH_ERROR: error message
	// 2022/01/14 12:01:52 example_test.go:34: INFO DISPLAY_TO_CONSOLE: display to console
	// 2022/01/14 12:01:52 zl.go:89: INFO FLUSH_LOG_BUFFER
	//
	// {"level":"INFO","time":"2022-01-14T12:16:03.274814+09:00","caller":"zl/zl.go:37","function":"github.com/nkmr-jp/zap-lightning/zl.Init.func1","message":"INIT_LOGGER","version":"f0ebbca","hostname":"nkmrnoMacBook-Pro.local","console":"logLevel: DEBUG, fileName: ./log/app.jsonl, outputType: Pretty"}
	// {"level":"INFO","time":"2022-01-14T12:16:03.27503+09:00","caller":"zl/example_test.go:28","function":"github.com/nkmr-jp/zap-lightning/zl_test.Example","message":"USER_INFO","version":"f0ebbca","hostname":"nkmrnoMacBook-Pro.local","name":"Alice","age":20}
	// {"level":"ERROR","time":"2022-01-14T12:16:03.275073+09:00","caller":"zl/example_test.go:30","function":"github.com/nkmr-jp/zap-lightning/zl_test.Example","message":"ERROR_MESSAGE","version":"f0ebbca","hostname":"nkmrnoMacBook-Pro.local","error":"error message","stacktrace":"github.com/nkmr-jp/zap-lightning/zl_test.Example\n\t/Users/nkmr/ghq/github.com/nkmr-jp/zap-lightning/zl/example_test.go:30\ntesting.runExample\n\t/Users/nkmr/.anyenv/envs/goenv/versions/1.17.5/src/testing/run_example.go:64\ntesting.runExamples\n\t/Users/nkmr/.anyenv/envs/goenv/versions/1.17.5/src/testing/example.go:44\ntesting.(*M).Run\n\t/Users/nkmr/.anyenv/envs/goenv/versions/1.17.5/src/testing/testing.go:1505\nmain.main\n\t_testmain.go:49\nruntime.main\n\t/Users/nkmr/.anyenv/envs/goenv/versions/1.17.5/src/runtime/proc.go:255"}
	// {"level":"DEBUG","time":"2022-01-14T12:16:03.275159+09:00","caller":"zl/example_test.go:31","function":"github.com/nkmr-jp/zap-lightning/zl_test.Example","message":"DEBUG_MESSAGE","version":"f0ebbca","hostname":"nkmrnoMacBook-Pro.local"}
	// {"level":"WARN","time":"2022-01-14T12:16:03.275186+09:00","caller":"zl/example_test.go:32","function":"github.com/nkmr-jp/zap-lightning/zl_test.Example","message":"WARN_MESSAGE","version":"f0ebbca","hostname":"nkmrnoMacBook-Pro.local"}
	// {"level":"WARN","time":"2022-01-14T12:16:03.275208+09:00","caller":"zl/example_test.go:33","function":"github.com/nkmr-jp/zap-lightning/zl_test.Example","message":"WARN_MESSAGE_WITH_ERROR","version":"f0ebbca","hostname":"nkmrnoMacBook-Pro.local","error":"error message"}
	// {"level":"INFO","time":"2022-01-14T12:16:03.275233+09:00","caller":"zl/example_test.go:34","function":"github.com/nkmr-jp/zap-lightning/zl_test.Example","message":"DISPLAY_TO_CONSOLE","version":"f0ebbca","hostname":"nkmrnoMacBook-Pro.local","console":"display to console"}

	// Output:
}

func ExampleSetVersion() {
	// ex.`go run -ldflags "-X main.version=v0.1.1 -X main.srcRootDir=$PWD" main.go`
	version = "v0.1.1"
	srcRootDir, _ = os.Getwd()

	// Set Options
	zl.SetLevel(zapcore.DebugLevel) // Default is InfoLevel
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
	zl.SetLevel(zapcore.DebugLevel) // Default is InfoLevel
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
