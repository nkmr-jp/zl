package zl

import (
	"fmt"
	"log"
	"strings"

	"go.uber.org/zap/zapcore"
)

type ConsoleType int

const (
	ConsoleTypeAll ConsoleType = iota
	ConsoleTypeError
	ConsoleTypeNone
)

func SetConsoleType(option ConsoleType) {
	consoleType = option
}

type OutputType int

const (
	// OutputTypeShortConsoleAndFile output simple console log and detail file log (default)
	OutputTypeShortConsoleAndFile OutputType = iota
	// OutputTypeConsoleAndFile output detail console log and file log
	OutputTypeConsoleAndFile
	// OutputTypeConsole output detail console log
	OutputTypeConsole
	// OutputTypeFile output detail file log
	OutputTypeFile
)

var outputTypeStrings = [4]string{
	"ShortConsoleAndFile",
	"ConsoleAndFile",
	"Console",
	"File",
}

func (o OutputType) String() string {
	return outputTypeStrings[o]
}

func SetOutputType(option OutputType) {
	outputType = option
}

// SetOutputTypeByString outputTypeStr can use (SimpleConsoleAndFile, ConsoleAndFile, Console, File).
func SetOutputTypeByString(outputTypeStr string) {
	var output OutputType
	if outputTypeStr == "" {
		SetOutputType(output)
		return
	}
	for i, i2 := range outputTypeStrings {
		if outputTypeStr == i2 {
			SetOutputType(OutputType(i))
			return
		}
	}
	log.Fatalf(
		"%s is invalid type. can use (SimpleConsoleAndFile, ConsoleAndFile, Console, File)",
		outputTypeStr,
	)
}

func SetLogLevel(option zapcore.Level) {
	logLevel = option
}

// SetLogLevelByString is set log level. levelStr can use (DEBUG,INFO,WARN,ERROR,FATAL).
func SetLogLevelByString(levelStr string) {
	var level zapcore.Level
	err := level.UnmarshalText([]byte(levelStr))
	if err != nil {
		log.Fatalf("%s is invalid level. can use (DEBUG,INFO,WARN,ERROR,FATAL)", levelStr)
	}
	SetLogLevel(level)
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
