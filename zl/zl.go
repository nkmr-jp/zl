package zl

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	consoleFieldDefault = "console"
)

var (
	once          sync.Once
	zapLogger     *zap.Logger
	consoleType   ConsoleType
	outputType    OutputType
	version       string
	logLevel      zapcore.Level // Default is InfoLevel
	callerEncoder zapcore.CallerEncoder
	consoleFields = []string{consoleFieldDefault}
)

// Initialize the Logger.
// Outputs short logs to the console and Write structured and detailed json logs to the log file.
func Init() *zap.Logger {
	once.Do(func() {
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
		initZapLogger()
		Info("INIT_LOGGER")
	})
	return zapLogger
}

// See https://pkg.go.dev/go.uber.org/zap
func initZapLogger() {
	log.Printf("log level: %v", logLevel.CapitalString())
	log.Printf("output type: %v", outputType.String())
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "name",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		FunctionKey:    "function",
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   getCallerEncoder(),
	}
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(getSyncers()...),
		logLevel,
	)
	zapLogger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel)).With(
		zap.String("version", getVersion()),
		zap.String("hostname", *getHost()),
	)
}

func getVersion() string {
	if version != "" {
		return version
	}
	if out, err := exec.Command("git", "rev-parse", "--short", "HEAD").Output(); err == nil {
		return strings.TrimRight(string(out), "\n")
	}

	return "undefined"
}

func getHost() *string {
	ret, err := os.Hostname()
	if err != nil {
		log.Print(err)
		return nil
	}
	return &ret
}

func getCallerEncoder() zapcore.CallerEncoder {
	if callerEncoder != nil {
		return callerEncoder
	}
	return zapcore.ShortCallerEncoder
}

func getSyncers() (syncers []zapcore.WriteSyncer) {
	switch outputType {
	case OutputTypeShortConsoleAndFile, OutputTypeFile:
		syncers = append(syncers, zapcore.AddSync(newLumberjack()))
	case OutputTypeConsoleAndFile:
		syncers = append(syncers, zapcore.AddSync(os.Stdout), zapcore.AddSync(newLumberjack()))
	case OutputTypeConsole:
		syncers = append(syncers, zapcore.AddSync(os.Stdout))
	}
	return
}
