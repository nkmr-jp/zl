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
	if outputType != PrettyOutput {
		return nil
	}
	l := log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile)
	if funk.Contains(omitKeys, TimeKey) {
		l.SetFlags(log.Lshortfile)
	}
	if isStdOut {
		l.SetOutput(os.Stdout)
	}
	return &prettyLogger{
		Logger: l,
	}
}

func (l *prettyLogger) log(msg string, level zapcore.Level, fields []zap.Field) {
	if outputType != PrettyOutput {
		return
	}
	err := l.Logger.Output(4,
		l.coloredLevel(level)+" "+l.coloredMsg(msg, level, fields),
	)
	if err != nil {
		l.Logger.Fatal(err)
	}
}

func (l *prettyLogger) logWithError(msg string, level zapcore.Level, err error, fields []zap.Field) {
	if outputType != PrettyOutput {
		return
	}
	err2 := l.Logger.Output(
		4,
		l.coloredLevel(level)+" "+l.coloredMsg(
			fmt.Sprintf("%s%s%s", msg, separator, aurora.Magenta(err.Error())),
			level, fields,
		),
	)
	if err2 != nil {
		l.Logger.Fatal(err2)
	}
}

func (l *prettyLogger) coloredMsg(msg string, level zapcore.Level, fields []zap.Field) string {
	var fieldMsg string
	if level == DebugLevel {
		msg = aurora.Faint(msg).String()
		fieldMsg = aurora.Faint(l.consoleMsg(fields)).String()
	} else {
		fieldMsg = l.consoleMsg(fields)
	}
	return fmt.Sprintf("%s%s", msg, fieldMsg)
}

func (l *prettyLogger) consoleMsg(fields []zap.Field) string {
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
			consoles = append(consoles, aurora.Cyan(val).String())
		}
	}
	if consoles != nil {
		ret = separator + fmt.Sprintf("%s", strings.Join(consoles, separator))
	}
	return ret
}

func (l *prettyLogger) coloredLevel(level zapcore.Level) string {
	switch level {
	case FatalLevel:
		return aurora.Red(level.CapitalString()).String()
	case ErrorLevel:
		return aurora.Red(level.CapitalString()).String()
	case WarnLevel:
		return aurora.Yellow(level.CapitalString()).String()
	case InfoLevel:
		return aurora.BrightBlue(level.CapitalString()).String()
	case DebugLevel:
		return aurora.BrightBlack(level.CapitalString()).String()
	}
	return ""
}
