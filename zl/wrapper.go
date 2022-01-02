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

type Wrapper struct {
	Fields []zap.Field
}

// NewWrapper can additional fields.
// ex. Use this when you want to add a common value in the scope of a context, such as an API request.
func NewWrapper(fields ...zap.Field) *Wrapper {
	return &Wrapper{Fields: fields}
}

func (w *Wrapper) Debug(msg string, fields ...zap.Field) {
	fields = append(fields, w.Fields...)
	wrapper(msg, "DEBUG", fields).Debug(msg, fields...)
}

func (w *Wrapper) Info(msg string, fields ...zap.Field) {
	fields = append(fields, w.Fields...)
	wrapper(msg, "INFO", fields).Info(msg, fields...)
}

func (w *Wrapper) Warn(msg string, fields ...zap.Field) {
	fields = append(fields, w.Fields...)
	wrapper(msg, "WARN", fields).Warn(msg, fields...)
}

func (w *Wrapper) Error(msg string, fields ...zap.Field) {
	fields = append(fields, w.Fields...)
	wrapper(msg, "ERROR", fields).Error(msg, fields...)
}

func (w *Wrapper) Fatal(msg string, fields ...zap.Field) {
	fields = append(fields, w.Fields...)
	wrapper(msg, "FATAL", fields).Fatal(msg, fields...)
}

func (w *Wrapper) Debugf(msg string, err error, fields ...zap.Field) {
	fields = append(addErrorField(fields, err), w.Fields...)
	wrapperf(msg, "DEBUG", err, fields).Debug(msg, fields...)
}

func (w *Wrapper) Infof(msg string, err error, fields ...zap.Field) {
	fields = append(addErrorField(fields, err), w.Fields...)
	wrapperf(msg, "INFO", err, fields).Info(msg, fields...)
}

func (w *Wrapper) Warnf(msg string, err error, fields ...zap.Field) {
	fields = append(addErrorField(fields, err), w.Fields...)
	wrapperf(msg, "WARN", err, fields).Warn(msg, fields...)
}

func (w *Wrapper) Errorf(msg string, err error, fields ...zap.Field) {
	fields = append(addErrorField(fields, err), w.Fields...)
	wrapperf(msg, "ERROR", err, fields).Error(msg, fields...)
}

func (w *Wrapper) Fatalf(msg string, err error, fields ...zap.Field) {
	fields = append(addErrorField(fields, err), w.Fields...)
	wrapperf(msg, "FATAL", err, fields).Fatal(msg, fields...)
}

// Sync wrapper of Zap's Sync.
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

// Debug is Wrapper of Zap's Debug.
// Outputs a short log to the console. Detailed json log output to log file.
func Debug(msg string, fields ...zap.Field) {
	wrapper(msg, "DEBUG", fields).Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	wrapper(msg, "INFO", fields).Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	wrapper(msg, "WARN", fields).Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	wrapper(msg, "ERROR", fields).Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	wrapper(msg, "FATAL", fields).Fatal(msg, fields...)
}

// Debugf is Outputs a Debug log with formatted error.
func Debugf(msg string, err error, fields ...zap.Field) {
	wrapperf(msg, "DEBUG", err, fields).Debug(msg, addErrorField(fields, err)...)
}

func Infof(msg string, err error, fields ...zap.Field) {
	wrapperf(msg, "INFO", err, fields).Info(msg, addErrorField(fields, err)...)
}

func Warnf(msg string, err error, fields ...zap.Field) {
	wrapperf(msg, "WARN", err, fields).Warn(msg, addErrorField(fields, err)...)
}

func Errorf(msg string, err error, fields ...zap.Field) {
	wrapperf(msg, "ERROR", err, fields).Error(msg, addErrorField(fields, err)...)
}

func Fatalf(msg string, err error, fields ...zap.Field) {
	wrapperf(msg, "FATAL", err, fields).Fatal(msg, addErrorField(fields, err)...)
}

func wrapper(msg, level string, fields []zap.Field) *zap.Logger {
	checkInit()
	shortLog(msg, level, fields)
	return zapLogger.WithOptions(zap.AddCallerSkip(1))
}

func wrapperf(msg, level string, err error, fields []zap.Field) *zap.Logger {
	checkInit()
	shortLogWithError(msg, level, err, fields)
	return zapLogger.WithOptions(zap.AddCallerSkip(1))
}

func addErrorField(fields []zap.Field, err error) []zap.Field {
	return append(fields, zap.String("error", fmt.Sprintf("%+v", err)))
}

func checkInit() {
	if zapLogger == nil {
		log.Fatal("The logger is not initialized. InitLogger() must be called.")
	}
}
