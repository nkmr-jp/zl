package zl

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
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

func (l *prettyLogger) printTraces() {
	if funk.Contains(omitKeys, StacktraceKey) && funk.Contains(omitKeys, PIDKey) {
		return
	}

	fp, err := os.Open(fileName)
	if err != nil {
		l.Logger.Fatal(err)
	}
	defer func(fp *os.File) {
		err := fp.Close()
		if err != nil {
			l.Logger.Fatal(err)
		}
	}(fp)

	scanner := bufio.NewScanner(fp)
	var traces string
	count := 0
	ln := 1
	for scanner.Scan() {
		trace := l.buildTrace(ln, scanner)
		if trace != "" {
			count++
		}
		traces += trace
		ln++
	}

	if count == 0 {
		return
	}

	head := aurora.BgRed(fmt.Sprintf(
		"                              %v ERROR OCCURRED                              ",
		count,
	))
	output := fmt.Sprintf("\n\n\n%s\n\n\n\n%s", head, traces)
	if isStdOut {
		if _, err := fmt.Fprint(os.Stdout, output); err != nil {
			return
		}
	} else {
		if _, err := fmt.Fprint(os.Stderr, output); err != nil {
			return
		}
	}

	if err = scanner.Err(); err != nil {
		l.Logger.Fatal(err)
	}
}

func (l *prettyLogger) buildTrace(ln int, scanner *bufio.Scanner) string {
	var trace Trace
	var output string
	if err := json.Unmarshal(scanner.Bytes(), &trace); err != nil {
		return ""
	}
	logFile := fmt.Sprintf("%v:%v", filepath.Base(fileName), ln)
	msg := l.coloredLevel(trace.Level) + " " + l.coloredMsg(
		fmt.Sprintf("%s%s%s", trace.Message, separator, aurora.Magenta(trace.Error)),
		trace.Level, nil,
	)
	if trace.Stacktrace != "" && trace.Pid == pid {
		output = fmt.Sprintf("%v ( %s )\n\n%v\n\n\n\n", msg, logFile, trace.Stacktrace)
	}
	return output
}
