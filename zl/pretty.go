package zl

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/logrusorgru/aurora"
	"github.com/thoas/go-funk"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type prettyLogger struct {
	Logger *log.Logger
}

func newPrettyLogger() *prettyLogger {
	l := log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile)
	if funk.Contains(ignoreKeys, TimeKey) {
		l.SetFlags(log.Lshortfile)
	}
	if isStdOut {
		l.SetOutput(os.Stdout)
	}
	return &prettyLogger{
		Logger: l,
	}
}

func (l *prettyLogger) Log(msg, levelStr string, fields []zap.Field) {
	if outputType != PrettyOutput {
		return
	}
	if !checkLevel(levelStr) {
		return
	}

	var fieldMsg string
	if levelStr == "DEBUG" {
		msg = aurora.Faint(msg).String()
		fieldMsg = aurora.Faint(getConsoleMsg(fields)).String()
	} else {
		fieldMsg = getConsoleMsg(fields)
	}

	err := l.Logger.Output(4, fmt.Sprintf("%v %v%v", color(levelStr), msg, fieldMsg))
	if err != nil {
		l.Logger.Fatal(err)
	}
}

func (l *prettyLogger) LogWithError(msg string, levelStr string, err error, fields []zap.Field) {
	if outputType != PrettyOutput {
		return
	}
	if !checkLevel(levelStr) {
		return
	}
	err2 := l.Logger.Output(
		4,
		fmt.Sprintf(
			"%v %v: %v %v",
			color(levelStr), msg, aurora.Magenta(err.Error()), getConsoleMsg(fields),
		),
	)
	if err2 != nil {
		l.Logger.Fatal(err2)
	}
}

func checkLevel(levelStr string) bool {
	var level zapcore.Level
	if err := level.UnmarshalText([]byte(levelStr)); err != nil {
		return false
	}
	if logLevel > level {
		return false
	}
	return true
}

func getConsoleMsg(fields []zap.Field) string {
	var ret string
	var consoles []string
	for i := range fields {
		if funk.ContainsString(consoleFields, fields[i].Key) {
			var val string
			if fields[i].String != "" {
				val = fields[i].String
			} else {
				val = strconv.Itoa(int(fields[i].Integer))
			}
			// consoles = append(consoles, fmt.Sprintf("%s=%s", v.Key, val))
			consoles = append(consoles, val)
		}
	}
	if consoles != nil {
		ret = ": " + fmt.Sprintf("%v", aurora.Cyan(strings.Join(consoles, ", ")))
	}
	return ret
}

func color(level string) string {
	switch level {
	case "FATAL":
		level = aurora.Red(level).String()
	case "ERROR":
		level = aurora.Red(level).String()
	case "WARN":
		level = aurora.Yellow(level).String()
	case "INFO":
		level = aurora.BrightBlue(level).String()
	case "DEBUG":
		level = aurora.BrightBlack(level).String()
	}
	return level
}
