package zl

import (
	"fmt"
	"log"
	"strings"

	"go.uber.org/zap/zapcore"
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

// SetOutputTypeByString can use (SimpleConsoleAndFile, ConsoleAndFile, Console, File).
func SetOutputTypeByString(outputTypeStr string) {
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
		"%s is invalid type. can use (SimpleConsoleAndFile, ConsoleAndFile, Console, File)",
		outputTypeStr,
	)
}

func SetLevel(option zapcore.Level) {
	logLevel = option
}

// SetLogLevelByString is set log level. levelStr can use (DEBUG,INFO,WARN,ERROR,FATAL).
func SetLogLevelByString(levelStr string) {
	var level zapcore.Level
	err := level.UnmarshalText([]byte(levelStr))
	if err != nil {
		log.Fatalf("%s is invalid level. can use (DEBUG,INFO,WARN,ERROR,FATAL)", levelStr)
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

// SetVersion set version.
// note: `revisionOrTag` should be a git revision or a tag. ex. `e86b9a7` or `v1.0.0`.
func SetVersion(revisionOrTag string) {
	version = revisionOrTag
}

// AddConsoleField Add the fields to be displayed in the console.
func AddConsoleField(fieldKey ...string) {
	consoleFields = append(consoleFields, fieldKey...)
}

//
// Log File Rotation Settings
// See: https://github.com/natefinch/lumberjack#type-logger
//

func SetFileName(val string) {
	fileName = val
}

func SetMaxSize(val int) {
	maxSize = val
}

func SetMaxBackups(val int) {
	maxBackups = val
}

func SetMaxAge(val int) {
	maxAge = val
}

func SetLocalTime(val bool) {
	localTime = val
}

func SetCompress(val bool) {
	compress = val
}
