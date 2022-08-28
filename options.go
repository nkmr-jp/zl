package zl

import (
	"fmt"
	"log"
	"strings"

	"go.uber.org/zap/zapcore"
)

// Key defines a commonly used field name for each log entry
type Key string

// See: https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry
const (
	MessageKey    Key = "message"
	LevelKey      Key = "severity"
	TimeKey       Key = "timestamp"
	LoggerKey     Key = "logger"
	CallerKey     Key = "caller"
	FunctionKey   Key = "function"
	StacktraceKey Key = "stacktrace"
	VersionKey    Key = "version"
	HostnameKey   Key = "hostname"
	PIDKey        Key = "pid"
)

type ErrorGroup struct {
	ErrorLogs []*ErrorLog
	Key       string
}

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

//
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

// SetOutputByString outputTypeStr can use (Pretty, ConsoleAndFile, Console, File).
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

func SetLevel(option zapcore.Level) {
	severityLevel = option
}

// SetLevelByString is set log level. levelStr can use (DEBUG,INFO,WARN,ERROR,FATAL).
func SetLevelByString(levelStr string) {
	var level zapcore.Level
	err := level.UnmarshalText([]byte(levelStr))
	if err != nil {
		log.Fatalf("%s is invalid level. can use (DEBUG, INFO, WARN, ERROR, FATAL)", levelStr)
	}
	SetLevel(level)
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

// SetConsoleFields add the fields to be displayed in the console.
func SetConsoleFields(fieldKey ...string) {
	consoleFields = append(consoleFields, fieldKey...)
}

// SetOmitKeys set fields to omit from default fields that used in each log.
func SetOmitKeys(key ...Key) {
	omitKeys = key
}

// SetStdout is changes the console log output from stderr to stdout.
func SetStdout() {
	isStdOut = true
}

// SetSeparator is changes the console log output separator.
func SetSeparator(val string) {
	separator = val
}

// SetRotateFileName set the file to write logs to.
// See: https://github.com/natefinch/lumberjack#type-logger
func SetRotateFileName(val string) {
	fileName = val
}

// SetRotateMaxSize set the maximum size in megabytes of the log file before it gets rotated.
// See: https://github.com/natefinch/lumberjack#type-logger
func SetRotateMaxSize(val int) {
	maxSize = val
}

// SetRotateMaxAge set the maximum number of days to retain.
// See: https://github.com/natefinch/lumberjack#type-logger
func SetRotateMaxAge(val int) {
	maxAge = val
}

// SetRotateMaxBackups set the maximum number of old log files to retain.
// See: https://github.com/natefinch/lumberjack#type-logger
func SetRotateMaxBackups(val int) {
	maxBackups = val
}

// SetRotateLocalTime determines if the time used for formatting the timestamps in backup files is the computer's local time.
// See: https://github.com/natefinch/lumberjack#type-logger
func SetRotateLocalTime(val bool) {
	localTime = val
}

// SetRotateCompress determines if the rotated log files should be compressed using gzip.
// See: https://github.com/natefinch/lumberjack#type-logger
func SetRotateCompress(val bool) {
	compress = val
}
