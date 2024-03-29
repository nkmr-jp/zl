// Package zl is a logger based on zap.
// It provides advanced logging features.
// It is designed with the developer's experience in mind and allows
// the user to choose the best output format for their purposes.
package zl

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/samber/lo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	consoleFieldDefault = "console"
)

var (
	once           sync.Once
	pretty         *prettyLogger
	zapLogger      *zap.Logger
	encoderConfig  *zapcore.EncoderConfig
	internalLogger *zap.Logger
	outputType     Output
	version        string
	severityLevel  zapcore.Level // Default is InfoLevel
	callerEncoder  zapcore.CallerEncoder
	consoleFields  = []string{consoleFieldDefault}
	omitKeys       []Key
	fieldKeys      = make(map[Key]string)
	isStdOut       bool
	separator      = " "
	pid            int
	isTest         bool
)

type fatalHook struct{}

func (f fatalHook) OnWrite(_ *zapcore.CheckedEntry, _ []zapcore.Field) {
	pretty.showErrorReport(fileName, pid)
	if isTest {
		fmt.Println("os.Exit(1) called.")
	} else {
		os.Exit(1)
	}
}

// Init initializes the logger.
func Init() {
	once.Do(func() {
		encoderConfig = newEncoderConfig()
		zapLogger = newLogger(encoderConfig)
		if outputType == PrettyOutput || isTest {
			pretty = newPrettyLogger(getConsoleOutput(), os.Stderr)
			zapLogger = zapLogger.WithOptions(zap.WithFatalHook(fatalHook{}))
		}

		encInternal := newEncoderConfig()
		encInternal.EncodeCaller = zapcore.ShortCallerEncoder
		internalLogger = newLogger(encInternal)

		var p, f string
		if pid != 0 {
			p = fmt.Sprintf(", PID: %d", pid)
		}
		if outputType == PrettyOutput || outputType == ConsoleAndFileOutput {
			f = fmt.Sprintf(", File: %s", fileName)
		}

		c := fmt.Sprintf(
			"Severity: %s, Output: %s%s%s",
			severityLevel.CapitalString(),
			outputType.String(),
			f,
			p,
		)
		iDebug("INIT_LOGGER", Console(c))
	})
}

func newEncoderConfig() *zapcore.EncoderConfig {
	enc := zapcore.EncoderConfig{
		MessageKey:     fieldKey(MessageKey),
		LevelKey:       fieldKey(LevelKey),
		TimeKey:        fieldKey(TimeKey),
		NameKey:        fieldKey(LoggerKey),
		CallerKey:      fieldKey(CallerKey),
		FunctionKey:    fieldKey(FunctionKey),
		StacktraceKey:  fieldKey(StacktraceKey),
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   getCallerEncoder(),
	}
	setOmitKeys(&enc)
	return &enc
}

func fieldKey(key Key) string {
	value, ok := fieldKeys[key]
	if !ok {
		return string(key)
	}
	return value
}

// See https://pkg.go.dev/go.uber.org/zap
func newLogger(enc *zapcore.EncoderConfig) *zap.Logger {
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(*enc),
		zapcore.NewMultiWriteSyncer(getSyncers()...),
		severityLevel,
	)
	return zap.New(core,
		zap.AddCallerSkip(1),
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	).With(getAdditionalFields()...)
}

func setOmitKeys(enc *zapcore.EncoderConfig) {
	for i := range omitKeys {
		switch omitKeys[i] {
		case MessageKey:
			enc.MessageKey = zapcore.OmitKey
		case LevelKey:
			enc.LevelKey = zapcore.OmitKey
		case TimeKey:
			enc.TimeKey = zapcore.OmitKey
		case LoggerKey:
			enc.NameKey = zapcore.OmitKey
		case CallerKey:
			enc.CallerKey = zapcore.OmitKey
		case FunctionKey:
			enc.FunctionKey = zapcore.OmitKey
		case StacktraceKey:
			enc.StacktraceKey = zapcore.OmitKey
		}
	}
}

func getAdditionalFields() (fields []zapcore.Field) {
	if !lo.Contains(omitKeys, VersionKey) {
		fields = append(fields, zap.String(string(VersionKey), GetVersion()))
	}
	if !lo.Contains(omitKeys, HostnameKey) {
		fields = append(fields, zap.String(string(HostnameKey), *getHost()))
	}
	if !lo.Contains(omitKeys, PIDKey) {
		pid = os.Getpid()
		fields = append(fields, zap.Int(string(PIDKey), pid))
	}
	return fields
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

// Sync is wrapper of Zap's Sync.
//
// Flushes any buffered log entries.(See: https://pkg.go.dev/go.uber.org/zap#Logger.Sync)
// Applications should take care to call Sync before exiting.
//
// Also displays an error report with a formatted stack trace if the outputType is PrettyOutput.
// This is useful for finding the source of errors during development.
//
// An error will occur if zap's Sync is executed when the output destination is console.
// (See: https://github.com/uber-go/zap/issues/880 )
// Therefore, Sync is executed only when console is not included in the zap output destination.
func Sync() {
	if outputType != PrettyOutput && outputType != FileOutput {
		return
	}
	if err := zapLogger.Sync(); err != nil {
		log.Println(err)
	}
	if outputType == PrettyOutput {
		pretty.showErrorReport(fileName, pid)
	}
}

// SyncWhenStop flush log buffer. when interrupt or terminated.
func SyncWhenStop() {
	if outputType != PrettyOutput && outputType != FileOutput {
		return
	}

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

		iDebug(fmt.Sprintf("GOT_SIGNAL_%v", strings.ToUpper(s.String())))
		Sync() // flush log buffer

		if isTest {
			fmt.Printf("os.Exit(%d) called.", 128+sigCode)
		} else {
			os.Exit(128 + sigCode)
		}
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
		syncers = append(syncers, zapcore.AddSync(getConsoleOutput()), zapcore.AddSync(newRotator()))
	case ConsoleOutput:
		syncers = append(syncers, zapcore.AddSync(getConsoleOutput()))
	}
	return
}

func getConsoleOutput() io.Writer {
	if isStdOut {
		return os.Stdout
	} else {
		return os.Stderr
	}
}

// ResetGlobalLoggerSettings resets global logger settings.
// This is convenient for use in tests, etc.
func ResetGlobalLoggerSettings() {
	once = sync.Once{}
	pretty = nil
	zapLogger = nil
	encoderConfig = nil
	internalLogger = nil
	outputType = PrettyOutput
	version = ""
	severityLevel = zapcore.InfoLevel
	callerEncoder = nil
	consoleFields = []string{consoleFieldDefault}
	omitKeys = nil
	fieldKeys = make(map[Key]string)
	isStdOut = false
	separator = " "
	fileName = ""
	maxSize = 0
	maxBackups = 0
	maxAge = 0
	localTime = false
	compress = false
}

// Cleanup
// Deprecated: Use ResetGlobalLoggerSettings instead.
// # codecov ignore
func Cleanup() {
	ResetGlobalLoggerSettings()
}

// SetIsTest sets isTest flag to true.
func SetIsTest() {
	isTest = true
}
