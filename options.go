package zl

import (
	"fmt"
	"log"
	"strings"

	"go.uber.org/zap/zapcore"
)

// Key defines a commonly used field name for each log entry.
// Each field defined in Key is output to all logs by default.
// Unnecessary fields can also be excluded using SetOmitKeys.
//
// Field names such as LevelKey and TimeKey are defined with reference to Google Cloud Logging.
// See: https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry
type Key string

const (
	// zapcore.EncoderConfig fields

	// MessageKey is set to zapcore.EncoderConfig.MessageKey
	MessageKey Key = "message"
	// LevelKey is set to zapcore.EncoderConfig.LevelKey
	LevelKey Key = "severity"
	// TimeKey is set to zapcore.EncoderConfig.TimeKey
	TimeKey Key = "timestamp"
	// LoggerKey is set to zapcore.EncoderConfig.NameKey
	LoggerKey Key = "logger"
	// CallerKey is set to zapcore.EncoderConfig.CallerKey
	CallerKey Key = "caller"
	// FunctionKey is set to zapcore.EncoderConfig.FunctionKey
	FunctionKey Key = "function"
	// StacktraceKey is set to zapcore.EncoderConfig.StacktraceKey
	StacktraceKey Key = "stacktrace"

	// Additional fields

	// VersionKey is the name of the field that outputs the version of the application.
	VersionKey Key = "version"
	// HostnameKey is the name of the field that outputs the hostname of the machine.
	HostnameKey Key = "hostname"
	// PIDKey is the name of the field that outputs the process ID of the application.
	PIDKey Key = "pid"
)

// ErrorGroup is a group of ErrorLog.
// It is used in prettyLogger's error report.
type ErrorGroup struct {
	ErrorLogs []*ErrorLog
	Key       string
}

// ErrorLog is a log that contains error information.
// It is used in prettyLogger's error report.
type ErrorLog struct {
	Severity   zapcore.Level `json:"severity"`
	Timestamp  string        `json:"timestamp"`
	Caller     string        `json:"caller"`
	Message    string        `json:"message"`
	Error      string        `json:"error"`
	Stacktrace string        `json:"stacktrace"`
	Pid        int           `json:"pid"`
	Line       int
}

const (
	DebugLevel = zapcore.DebugLevel
	InfoLevel  = zapcore.InfoLevel
	WarnLevel  = zapcore.WarnLevel
	ErrorLevel = zapcore.ErrorLevel
	FatalLevel = zapcore.FatalLevel
)

// Output is log output type.
type Output int

const (
	// PrettyOutput writes the colored simple log to console,
	// and writes json structured detail log to file.
	// It is Default setting.
	// Recommended for Develop Environment.
	PrettyOutput Output = iota

	// ConsoleAndFileOutput writes json structured log to console and file.
	// Recommended for Develop Environment.
	ConsoleAndFileOutput

	// ConsoleOutput writes json structured log to console.
	// Recommended for Develop and Production Environment.
	ConsoleOutput

	// FileOutput writes json structured log to file.
	// Recommended for Develop and Production Environment.
	FileOutput
)

var outputStrings = [4]string{
	"Pretty",
	"ConsoleAndFile",
	"Console",
	"File",
}

// String is return Output type string.
func (o Output) String() string {
	return outputStrings[o]
}

// SetOutput is set Output type.
// option can use (PrettyOutput, ConsoleAndFileOutput, ConsoleOutput, FileOutput).
func SetOutput(option Output) {
	outputType = option
}

// SetOutputByString is set Output type by string.
// outputTypeStr can use (Pretty, ConsoleAndFile, Console, File).
func SetOutputByString(outputTypeStr string) {
	var output Output
	if outputTypeStr == "" {
		SetOutput(output)
		return
	}
	for i, i2 := range outputStrings {
		if outputTypeStr == i2 {
			SetOutput(Output(i))
			return
		}
	}
	log.Fatalf(
		"%s is invalid type. can use (Pretty, ConsoleAndFile, Console, File)",
		outputTypeStr,
	)
}

// SetLevel is set log.
// level can use (DebugLevel, InfoLevel, WarnLevel, ErrorLevel, FatalLevel).
func SetLevel(level zapcore.Level) {
	severityLevel = level
}

// SetLevelByString is set log level.
// levelStr can use (DEBUG, INFO, WARN, ERROR, FATAL).
func SetLevelByString(levelStr string) {
	var level zapcore.Level
	err := level.UnmarshalText([]byte(levelStr))
	if err != nil {
		log.Fatalf("%s is invalid level. can use (DEBUG, INFO, WARN, ERROR, FATAL)", levelStr)
	}
	SetLevel(level)
}

// SetRepositoryCallerEncoder is set CallerEncoder.
// It set caller's source code's URL of the Repository that called.
// It is used in the log output CallerKey field.
func SetRepositoryCallerEncoder(urlFormat, revisionOrTag, srcRootDir string) {
	if revisionOrTag == "" || srcRootDir == "" {
		return
	}
	url := fmt.Sprintf(urlFormat, revisionOrTag)
	callerEncoder = buildRepositoryCallerEncoder(srcRootDir, url)
}

func buildRepositoryCallerEncoder(dir, url string) zapcore.CallerEncoder {
	return func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(
			fmt.Sprintf("%v#L%v", strings.Replace(caller.File, dir, url, 1), caller.Line),
		)
	}
}

// SetVersion `revisionOrTag` should be a git revision or a tag. ex. `e86b9a7` or `v1.0.0`.
// It set version of the application.
// It is used in the log output VersionKey field.
func SetVersion(revisionOrTag string) {
	version = revisionOrTag
}

// SetConsoleFields add the fields to be displayed in the console when PrettyOutput is used.
func SetConsoleFields(fieldKey ...string) {
	consoleFields = append(consoleFields, fieldKey...)
}

// SetOmitKeys set fields to omit from default fields that used in each log.
func SetOmitKeys(key ...Key) {
	omitKeys = key
}

// SetFieldKey is changes the key of the default field.
func SetFieldKey(key Key, val string) {
	if key == "" || val == "" {
		return
	}
	fieldKeys[key] = val
}

// SetStdout is changes the console log output from stderr to stdout.
func SetStdout() {
	isStdOut = true
}

// SetSeparator is changes the console log output separator when PrettyOutput is used.
func SetSeparator(val string) {
	separator = val
}
