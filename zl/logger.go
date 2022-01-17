package zl

import (
	"log"

	"go.uber.org/zap"
)

type zlLogger struct {
	pretty    *prettyLogger
	zapLogger *zap.Logger
	fields    []zap.Field
}

// New can add additional default fields.
// ex. Use this when you want to add a common value in the scope of a context, such as an API request.
func New(fields ...zap.Field) *zlLogger {
	return &zlLogger{
		pretty:    newPrettyLogger(),
		zapLogger: newZapLogger(),
		fields:    fields,
	}
}

func (l *zlLogger) Named(name string) *zlLogger {
	l.zapLogger = l.zapLogger.Named(name)
	return l
}

func (l *zlLogger) Debug(msg string, fields ...zap.Field) {
	fields = append(fields, l.fields...)
	l.logger(msg, "DEBUG", fields).Debug(msg, fields...)
}

func (l *zlLogger) Info(msg string, fields ...zap.Field) {
	fields = append(fields, l.fields...)
	l.logger(msg, "INFO", fields).Info(msg, fields...)
}

func (l *zlLogger) Warn(msg string, fields ...zap.Field) {
	fields = append(fields, l.fields...)
	l.logger(msg, "WARN", fields).Warn(msg, fields...)
}

func (l *zlLogger) Error(msg string, err error, fields ...zap.Field) {
	fields = append(append(fields, zap.Error(err)), l.fields...)
	l.loggerErr(msg, "ERROR", err, fields).Error(msg, fields...)
}

func (l *zlLogger) Fatal(msg string, err error, fields ...zap.Field) {
	fields = append(append(fields, zap.Error(err)), l.fields...)
	l.loggerErr(msg, "FATAL", err, fields).Fatal(msg, fields...)
}

func (l *zlLogger) DebugErr(msg string, err error, fields ...zap.Field) {
	fields = append(append(fields, zap.Error(err)), l.fields...)
	l.loggerErr(msg, "DEBUG", err, fields).Debug(msg, fields...)
}

func (l *zlLogger) InfoErr(msg string, err error, fields ...zap.Field) {
	fields = append(append(fields, zap.Error(err)), l.fields...)
	l.loggerErr(msg, "INFO", err, fields).Info(msg, fields...)
}

func (l *zlLogger) WarnErr(msg string, err error, fields ...zap.Field) {
	fields = append(append(fields, zap.Error(err)), l.fields...)
	l.loggerErr(msg, "WARN", err, fields).Warn(msg, fields...)
}

func (l *zlLogger) logger(msg, level string, fields []zap.Field) *zap.Logger {
	l.pretty.Log(msg, level, fields)
	return l.zapLogger.WithOptions(zap.AddCallerSkip(1))
}

func (l *zlLogger) loggerErr(msg, level string, err error, fields []zap.Field) *zap.Logger {
	l.pretty.LogWithError(msg, level, err, fields)
	return l.zapLogger.WithOptions(zap.AddCallerSkip(1))
}

// Debug is zlLogger of Zap's Debug.
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
	pretty.Log(msg, level, fields)
	return zapLogger.WithOptions(zap.AddCallerSkip(1))
}

func loggerErr(msg, level string, err error, fields []zap.Field) *zap.Logger {
	checkInit()
	pretty.LogWithError(msg, level, err, fields)
	return zapLogger.WithOptions(zap.AddCallerSkip(1))
}

func checkInit() {
	if zapLogger == nil {
		log.Fatal("The logger is not initialized. Init() must be called.")
	}
}
