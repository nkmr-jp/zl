package zl

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"go.uber.org/zap"
)

type Logger struct {
	Fields []zap.Field
}

// NewLogger can additional default fields.
// ex. Use this when you want to add a common value in the scope of a context, such as an API request.
func NewLogger(fields ...zap.Field) *Logger {
	return &Logger{Fields: fields}
}

func (w *Logger) Debug(msg string, fields ...zap.Field) {
	fields = append(fields, w.Fields...)
	logger(msg, "DEBUG", fields).Debug(msg, fields...)
}

func (w *Logger) Info(msg string, fields ...zap.Field) {
	fields = append(fields, w.Fields...)
	logger(msg, "INFO", fields).Info(msg, fields...)
}

func (w *Logger) Warn(msg string, fields ...zap.Field) {
	fields = append(fields, w.Fields...)
	logger(msg, "WARN", fields).Warn(msg, fields...)
}

func (w *Logger) Error(msg string, err error, fields ...zap.Field) {
	fields = append(append(fields, zap.Error(err)), w.Fields...)
	loggerErr(msg, "ERROR", err, fields).Error(msg, fields...)
}

func (w *Logger) Fatal(msg string, err error, fields ...zap.Field) {
	fields = append(append(fields, zap.Error(err)), w.Fields...)
	loggerErr(msg, "FATAL", err, fields).Fatal(msg, fields...)
}

func (w *Logger) DebugErr(msg string, err error, fields ...zap.Field) {
	fields = append(append(fields, zap.Error(err)), w.Fields...)
	loggerErr(msg, "DEBUG", err, fields).Debug(msg, fields...)
}

func (w *Logger) InfoErr(msg string, err error, fields ...zap.Field) {
	fields = append(append(fields, zap.Error(err)), w.Fields...)
	loggerErr(msg, "INFO", err, fields).Info(msg, fields...)
}

func (w *Logger) WarnErr(msg string, err error, fields ...zap.Field) {
	fields = append(append(fields, zap.Error(err)), w.Fields...)
	loggerErr(msg, "WARN", err, fields).Warn(msg, fields...)
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

// Debug is Logger of Zap's Debug.
// Outputs a short log to the console. Detailed json log output to log file.
func Debug(msg string, fields ...zap.Field) {
	logger(msg, "DEBUG", fields).Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	logger(msg, "INFO", fields).Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	logger(msg, "WARN", fields).Warn(msg, fields...)
}

func Error(msg string, err error, fields ...zap.Field) {
	loggerErr(msg, "ERROR", err, fields).Error(msg, append(fields, zap.Error(err))...)
}

func Fatal(msg string, err error, fields ...zap.Field) {
	loggerErr(msg, "FATAL", err, fields).Fatal(msg, append(fields, zap.Error(err))...)
}

// DebugErr is Outputs a Debug log with error field.
func DebugErr(msg string, err error, fields ...zap.Field) {
	loggerErr(msg, "DEBUG", err, fields).Debug(msg, append(fields, zap.Error(err))...)
}

// InfoErr is Outputs a Info log with error field.
func InfoErr(msg string, err error, fields ...zap.Field) {
	err.Error()
	loggerErr(msg, "INFO", err, fields).Info(msg, append(fields, zap.Error(err))...)
}

// WarnErr is Outputs a Warn log with error field.
func WarnErr(msg string, err error, fields ...zap.Field) {
	loggerErr(msg, "WARN", err, fields).Warn(msg, append(fields, zap.Error(err))...)
}

func logger(msg, level string, fields []zap.Field) *zap.Logger {
	checkInit()
	shortLog(msg, level, fields)
	return zapLogger.WithOptions(zap.AddCallerSkip(1))
}

func loggerErr(msg, level string, err error, fields []zap.Field) *zap.Logger {
	checkInit()
	shortLogWithError(msg, level, err, fields)
	return zapLogger.WithOptions(zap.AddCallerSkip(1))
}

func checkInit() {
	if zapLogger == nil {
		log.Fatal("The logger is not initialized. Init() must be called.")
	}
}
