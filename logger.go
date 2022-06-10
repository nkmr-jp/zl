package zl

import (
	"fmt"
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zlLogger struct {
	pretty    *prettyLogger
	zapLogger *zap.Logger
	fields    []zap.Field
}

// New can add additional default fields.
// e.g. Use this when you want to add a common value in the scope of a context, such as an API request.
func New(fields ...zap.Field) *zlLogger {
	return &zlLogger{
		pretty:    newPrettyLogger(),
		zapLogger: newLogger(encoderConfig),
		fields:    fields,
	}
}

func (l *zlLogger) Named(loggerName string) *zlLogger {
	l.pretty.Logger.SetPrefix(fmt.Sprintf("%s | ", loggerName))
	l.zapLogger = l.zapLogger.Named(loggerName)
	return l
}

func (l *zlLogger) Debug(message string, fields ...zap.Field) {
	fields = append(fields, l.fields...)
	l.logger(message, DebugLevel, fields).Debug(message, fields...)
}

func (l *zlLogger) Info(message string, fields ...zap.Field) {
	fields = append(fields, l.fields...)
	l.logger(message, InfoLevel, fields).Info(message, fields...)
}

func (l *zlLogger) Warn(message string, fields ...zap.Field) {
	fields = append(fields, l.fields...)
	l.logger(message, WarnLevel, fields).Warn(message, fields...)
}

func (l *zlLogger) Error(message string, err error, fields ...zap.Field) {
	fields = append(append(fields, zap.Error(err)), l.fields...)
	l.loggerErr(message, ErrorLevel, err, fields).Error(message, fields...)
}

func (l *zlLogger) Fatal(message string, err error, fields ...zap.Field) {
	fields = append(append(fields, zap.Error(err)), l.fields...)
	l.loggerErr(message, FatalLevel, err, fields).Fatal(message, fields...)
}

func (l *zlLogger) DebugErr(message string, err error, fields ...zap.Field) {
	fields = append(append(fields, zap.Error(err)), l.fields...)
	l.loggerErr(message, DebugLevel, err, fields).Debug(message, fields...)
}

func (l *zlLogger) InfoErr(message string, err error, fields ...zap.Field) {
	fields = append(append(fields, zap.Error(err)), l.fields...)
	l.loggerErr(message, InfoLevel, err, fields).Info(message, fields...)
}

func (l *zlLogger) WarnErr(message string, err error, fields ...zap.Field) {
	fields = append(append(fields, zap.Error(err)), l.fields...)
	l.loggerErr(message, WarnLevel, err, fields).Warn(message, fields...)
}

func (l *zlLogger) logger(message string, level zapcore.Level, fields []zap.Field) *zap.Logger {
	l.pretty.log(message, level, fields)
	return l.zapLogger
}

func (l *zlLogger) loggerErr(message string, level zapcore.Level, err error, fields []zap.Field) *zap.Logger {
	l.pretty.logWithError(message, level, err, fields)
	return l.zapLogger
}

// Debug is wrapper of Zap's Debug.
func Debug(message string, fields ...zap.Field) {
	logger(message, DebugLevel, fields).Debug(message, fields...)
}

// Info is wrapper of Zap's Info.
func Info(message string, fields ...zap.Field) {
	logger(message, InfoLevel, fields).Info(message, fields...)
}

// Warn is wrapper of Zap's Warn.
func Warn(message string, fields ...zap.Field) {
	logger(message, WarnLevel, fields).Warn(message, fields...)
}

// Error is wrapper of Zap's Error with error field.
func Error(message string, err error, fields ...zap.Field) {
	loggerErr(message, ErrorLevel, err, fields).Error(message, append(fields, zap.Error(err))...)
}

// Fatal is wrapper of Zap's Fatal.
func Fatal(message string, err error, fields ...zap.Field) {
	loggerErr(message, FatalLevel, err, fields).Fatal(message, append(fields, zap.Error(err))...)
}

// DebugErr is Outputs a DEBUG log with error field.
func DebugErr(message string, err error, fields ...zap.Field) {
	loggerErr(message, DebugLevel, err, fields).Debug(message, append(fields, zap.Error(err))...)
}

// InfoErr is Outputs INFO log with error field.
func InfoErr(message string, err error, fields ...zap.Field) {
	loggerErr(message, InfoLevel, err, fields).Info(message, append(fields, zap.Error(err))...)
}

// WarnErr is Outputs WARN log with error field.
func WarnErr(message string, err error, fields ...zap.Field) {
	loggerErr(message, WarnLevel, err, fields).Warn(message, append(fields, zap.Error(err))...)
}

// Dump is a deep pretty printer for Go data structures to aid in debugging.
// It is only works with PrettyOutput settings.
//
// It is wrapper of go-spew.
// See: https://github.com/davecgh/go-spew
func Dump(a ...interface{}) {
	checkInit()
	pretty.dump(a...)
}

func logger(message string, level zapcore.Level, fields []zap.Field) *zap.Logger {
	checkInit()
	pretty.log(message, level, fields)
	return zapLogger
}

func iDebug(message string, fields ...zap.Field) {
	iLogger(message, DebugLevel, fields).Debug(message, fields...)
}

func iLogger(message string, level zapcore.Level, fields []zap.Field) *zap.Logger {
	checkInit()
	pretty.log(message, level, fields)
	return internalLogger
}

func loggerErr(message string, level zapcore.Level, err error, fields []zap.Field) *zap.Logger {
	checkInit()
	pretty.logWithError(message, level, err, fields)
	return zapLogger
}

func checkInit() {
	if zapLogger == nil {
		log.Fatal("The logger is not initialized. Init() must be called.")
	}
}
