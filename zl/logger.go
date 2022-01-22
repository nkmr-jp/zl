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
		zapLogger: newZapLogger(),
		fields:    fields,
	}
}

func (l *zlLogger) Named(name string) *zlLogger {
	l.pretty.Logger.SetPrefix(fmt.Sprintf("%s | ", name))
	l.zapLogger = l.zapLogger.Named(name)
	return l
}

func (l *zlLogger) Debug(msg string, fields ...zap.Field) {
	fields = append(fields, l.fields...)
	l.logger(msg, DebugLevel, fields).Debug(msg, fields...)
}

func (l *zlLogger) Info(msg string, fields ...zap.Field) {
	fields = append(fields, l.fields...)
	l.logger(msg, InfoLevel, fields).Info(msg, fields...)
}

func (l *zlLogger) Warn(msg string, fields ...zap.Field) {
	fields = append(fields, l.fields...)
	l.logger(msg, WarnLevel, fields).Warn(msg, fields...)
}

func (l *zlLogger) Error(msg string, err error, fields ...zap.Field) {
	fields = append(append(fields, zap.Error(err)), l.fields...)
	l.loggerErr(msg, ErrorLevel, err, fields).Error(msg, fields...)
}

func (l *zlLogger) Fatal(msg string, err error, fields ...zap.Field) {
	fields = append(append(fields, zap.Error(err)), l.fields...)
	l.loggerErr(msg, FatalLevel, err, fields).Fatal(msg, fields...)
}

func (l *zlLogger) DebugErr(msg string, err error, fields ...zap.Field) {
	fields = append(append(fields, zap.Error(err)), l.fields...)
	l.loggerErr(msg, DebugLevel, err, fields).Debug(msg, fields...)
}

func (l *zlLogger) InfoErr(msg string, err error, fields ...zap.Field) {
	fields = append(append(fields, zap.Error(err)), l.fields...)
	l.loggerErr(msg, InfoLevel, err, fields).Info(msg, fields...)
}

func (l *zlLogger) WarnErr(msg string, err error, fields ...zap.Field) {
	fields = append(append(fields, zap.Error(err)), l.fields...)
	l.loggerErr(msg, WarnLevel, err, fields).Warn(msg, fields...)
}

func (l *zlLogger) logger(msg string, level zapcore.Level, fields []zap.Field) *zap.Logger {
	l.pretty.log(msg, level, fields)
	return l.zapLogger
}

func (l *zlLogger) loggerErr(msg string, level zapcore.Level, err error, fields []zap.Field) *zap.Logger {
	l.pretty.logWithError(msg, level, err, fields)
	return l.zapLogger
}

// Debug is wrapper of Zap's Debug.
func Debug(msg string, fields ...zap.Field) {
	logger(msg, DebugLevel, fields).Debug(msg, fields...)
}

// Info is wrapper of Zap's Info.
func Info(msg string, fields ...zap.Field) {
	logger(msg, InfoLevel, fields).Info(msg, fields...)
}

// Warn is wrapper of Zap's Warn.
func Warn(msg string, fields ...zap.Field) {
	logger(msg, WarnLevel, fields).Warn(msg, fields...)
}

// Error is wrapper of Zap's Error with error field.
func Error(msg string, err error, fields ...zap.Field) {
	loggerErr(msg, ErrorLevel, err, fields).Error(msg, append(fields, zap.Error(err))...)
}

// Fatal is wrapper of Zap's Fatal.
func Fatal(msg string, err error, fields ...zap.Field) {
	loggerErr(msg, FatalLevel, err, fields).Fatal(msg, append(fields, zap.Error(err))...)
}

// DebugErr is Outputs a Debug log with error field.
func DebugErr(msg string, err error, fields ...zap.Field) {
	loggerErr(msg, DebugLevel, err, fields).Debug(msg, append(fields, zap.Error(err))...)
}

// InfoErr is Outputs a Info log with error field.
func InfoErr(msg string, err error, fields ...zap.Field) {
	err.Error()
	loggerErr(msg, InfoLevel, err, fields).Info(msg, append(fields, zap.Error(err))...)
}

// WarnErr is Outputs a Warn log with error field.
func WarnErr(msg string, err error, fields ...zap.Field) {
	loggerErr(msg, WarnLevel, err, fields).Warn(msg, append(fields, zap.Error(err))...)
}

func logger(msg string, level zapcore.Level, fields []zap.Field) *zap.Logger {
	checkInit()
	pretty.log(msg, level, fields)
	return zapLogger
}

func loggerErr(msg string, level zapcore.Level, err error, fields []zap.Field) *zap.Logger {
	checkInit()
	pretty.logWithError(msg, level, err, fields)
	return zapLogger
}

func checkInit() {
	if zapLogger == nil {
		log.Fatal("The logger is not initialized. Init() must be called.")
	}
}
