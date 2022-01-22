package zl

import (
	"fmt"
	"strings"

	"go.uber.org/zap/zapcore"
)

// Key is used by each log entry.
type Key string

const (
	MessageKey    Key = "message"
	LevelKey      Key = "level"
	TimeKey       Key = "time"
	NameKey       Key = "name"
	CallerKey     Key = "caller"
	FunctionKey   Key = "function"
	StacktraceKey Key = "stacktrace"
	VersionKey    Key = "version"
	HostnameKey   Key = "hostname"
)

type Level zapcore.Level

const (
	DebugLevel = zapcore.DebugLevel
	InfoLevel  = zapcore.InfoLevel
	WarnLevel  = zapcore.WarnLevel
	ErrorLevel = zapcore.ErrorLevel
	FatalLevel = zapcore.FatalLevel
)

type Output int

const (
	// PrettyOutput writes the colored simple log to console,
	// and writes json structured detail log to file.
	// it is Default setting.
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

func (o Output) String() string {
	return outputStrings[o]
}

func SetOutput(option Output) {
	outputType = option
}

func SetLevel(option zapcore.Level) {
	logLevel = option
}

// SetRepositoryCallerEncoder is set CallerEncoder. it set caller's source code's URL of the Repository that called.
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
func SetVersion(revisionOrTag string) {
	version = revisionOrTag
}

// AddConsoleFields add the fields to be displayed in the console.
func AddConsoleFields(fieldKey ...string) {
	consoleFields = append(consoleFields, fieldKey...)
}

// SetIgnoreKeys set ignore fields from default fields that used in each log.
func SetIgnoreKeys(key ...Key) {
	ignoreKeys = key
}

// SetStdout is changes the console log output from stderr to stdout.
func SetStdout() {
	isStdOut = true
}

// SetSeparator is changes the console log output separator.
func SetSeparator(val string) {
	separator = val
}

// SetFileName set the file to write logs to.
// See: https://github.com/natefinch/lumberjack#type-logger
func SetFileName(val string) {
	fileName = val
}

// SetMaxSize set the maximum size in megabytes of the log file before it gets rotated.
// See: https://github.com/natefinch/lumberjack#type-logger
func SetMaxSize(val int) {
	maxSize = val
}

// SetMaxAge set the maximum number of days to retain.
// See: https://github.com/natefinch/lumberjack#type-logger
func SetMaxAge(val int) {
	maxAge = val
}

// SetMaxBackups set the maximum number of old log files to retain.
// See: https://github.com/natefinch/lumberjack#type-logger
func SetMaxBackups(val int) {
	maxBackups = val
}

// SetLocalTime determines if the time used for formatting the timestamps in backup files is the computer's local time.
// See: https://github.com/natefinch/lumberjack#type-logger
func SetLocalTime(val bool) {
	localTime = val
}

// SetCompress determines if the rotated log files should be compressed using gzip.
// See: https://github.com/natefinch/lumberjack#type-logger
func SetCompress(val bool) {
	compress = val
}
