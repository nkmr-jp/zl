// Package zl provides zap based advanced logging features, and it's easy to use.
package zl

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	consoleFieldDefault = "console"
)

var (
	once          sync.Once
	zapLogger     *zap.Logger
	outputType    Output
	version       string
	logLevel      zapcore.Level // Default is InfoLevel
	callerEncoder zapcore.CallerEncoder
	consoleFields = []string{consoleFieldDefault}
)

// Init initializes the logger.
func Init() *zap.Logger {
	once.Do(func() {
		log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
		initZapLogger()
		Info("INIT_LOGGER", Console(fmt.Sprintf(
			"logLevel: %s, fileName: %s, outputType: %s",
			logLevel.CapitalString(),
			fileName,
			outputType.String(),
		)))
	})
	return zapLogger
}

// See https://pkg.go.dev/go.uber.org/zap
func initZapLogger() {
	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "name",
		CallerKey:      "caller",
		FunctionKey:    "function",
		StacktraceKey:  "stacktrace",
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
		zap.String("version", GetVersion()),
		zap.String("hostname", *getHost()),
	)
}

// GetVersion return version when version is set.
// or return git commit hash when version is not set.
func GetVersion() string {
	if version != "" {
		return version
	}
	if out, err := exec.Command("git", "rev-parse", "--short", "HEAD").Output(); err == nil {
		return strings.TrimRight(string(out), "\n")
	}

	return "undefined"
}

// Sync logger of Zap's Sync.
// Note: If log output to console. error will occur (See: https://github.com/uber-go/zap/issues/880 )
func Sync() {
	Info("FLUSH_LOG_BUFFER")
	if err := zapLogger.Sync(); err != nil {
		log.Println(err)
	}
}

// SyncWhenStop flush log buffer. when interrupt or terminated.
func SyncWhenStop() {
	c := make(chan os.Signal, 1)

	go func() {
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		s := <-c

		sigCode := 0
		switch s.String() {
		case "interrupt":
			sigCode = 2
		case "terminated":
			sigCode = 15
		}

		Info(fmt.Sprintf("GOT_SIGNAL_%v", strings.ToUpper(s.String())))
		Sync() // flush log buffer
		os.Exit(128 + sigCode)
	}()
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
	case PrettyOutput, FileOutput:
		syncers = append(syncers, zapcore.AddSync(newRotator()))
	case ConsoleAndFileOutput:
		syncers = append(syncers, zapcore.AddSync(os.Stderr), zapcore.AddSync(newRotator()))
	case ConsoleOutput:
		syncers = append(syncers, zapcore.AddSync(os.Stderr))
	}
	return
}
