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

	"github.com/davecgh/go-spew/spew"
	au "github.com/logrusorgru/aurora"
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
	if outputType != PrettyOutput || level < severityLevel {
		return
	}
	err := l.Logger.Output(4,
		l.coloredLevel(level).String()+" "+l.coloredMsg(msg, level, fields),
	)
	if err != nil {
		l.Logger.Fatal(err)
	}
}

func (l *prettyLogger) logWithError(msg string, level zapcore.Level, err error, fields []zap.Field) {
	if outputType != PrettyOutput || level < severityLevel {
		return
	}
	err2 := l.Logger.Output(
		4,
		l.coloredLevel(level).String()+" "+l.coloredMsg(
			fmt.Sprintf("%s%s%s", msg, separator, au.Magenta(err.Error())),
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
		msg = au.Faint(msg).String()
		fieldMsg = au.Faint(l.consoleMsg(fields)).String()
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
			consoles = append(consoles, au.Cyan(val).String())
		}
	}
	if consoles != nil {
		ret = separator + fmt.Sprintf("%s", strings.Join(consoles, separator))
	}
	return ret
}

func (l *prettyLogger) coloredLevel(level zapcore.Level) au.Value {
	switch level {
	case FatalLevel:
		return au.Red(level.CapitalString())
	case ErrorLevel:
		return au.Red(level.CapitalString())
	case WarnLevel:
		return au.Yellow(level.CapitalString())
	case InfoLevel:
		return au.BrightBlue(level.CapitalString())
	case DebugLevel:
		return au.BrightBlack(level.CapitalString())
	}
	return au.BrightBlack("")
}

func (l *prettyLogger) printTraces() {
	if funk.Contains(omitKeys, StacktraceKey) && funk.Contains(omitKeys, PIDKey) {
		return
	}

	fp, err := os.Open(fileName)
	if err != nil {
		return
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
		trace := l.buildStackTrace(count, ln, scanner)
		if trace != "" {
			count++
		}
		traces += trace
		ln++
	}

	if count == 0 {
		return
	}

	head := au.Red(fmt.Sprintf(
		"\t\t\t\t%v ERROR OCCURRED\t\t\t\t",
		count,
	)).Bold()
	output := fmt.Sprintf("\n%s\n\n%s", head, traces)
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

func (l *prettyLogger) buildStackTrace(count, ln int, scanner *bufio.Scanner) string {
	var report ErrorReport
	var output string
	if err := json.Unmarshal(scanner.Bytes(), &report); err != nil {
		return ""
	}
	logFile := fmt.Sprintf("%v:%v", filepath.Base(fileName), ln)
	msg := l.coloredLevel(report.Severity).Bold().String() + " " + l.coloredMsg(
		fmt.Sprintf("%s%s%s", au.Bold(report.Message), separator, au.Magenta(report.Error)),
		report.Severity, nil,
	)
	if report.Stacktrace != "" && report.Pid == pid {
		output = fmt.Sprintf(
			"%v %v ( %s )\n\n\t%v\n\n",
			au.Red(fmt.Sprintf("[%d]", count+1)).Bold(),
			msg,
			logFile,
			strings.ReplaceAll(report.Stacktrace, "\n", "\n\t"),
		)
	}
	return output
}

func (l *prettyLogger) dump(a ...interface{}) {
	if outputType != PrettyOutput {
		return
	}
	err := l.Logger.Output(3,
		au.Red("DUMP").Bold().String()+" "+spew.Sdump(a...),
	)
	if err != nil {
		l.Logger.Fatal(err)
	}
}
