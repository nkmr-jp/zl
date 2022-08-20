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
		ret = separator + strings.Join(consoles, separator)
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

func (l *prettyLogger) showErrorReport() {
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

	count, traces := l.scanStackTraces(fp)
	l.printTraces(count, traces)
}

func (l *prettyLogger) scanStackTraces(fp *os.File) (int, string) {
	scanner := bufio.NewScanner(fp)
	var traces, key string
	var groups []*ErrorGroup

	count := 0
	ln := 0
	for scanner.Scan() {
		ln++
		var errorLog *ErrorLog
		var group *ErrorGroup
		flg := false
		if err := json.Unmarshal(scanner.Bytes(), &errorLog); err != nil {
			continue
		}
		if errorLog.Stacktrace == "" || errorLog.Pid != pid {
			continue
		}
		key = fmt.Sprintf("severity:%s,message:%s,caller:%s,error:%s",
			errorLog.Severity, errorLog.Message, errorLog.Error, errorLog.Caller,
		)
		errorLog.Line = ln
		for i := range groups {
			if groups[i].Key == key {
				groups[i].ErrorLogs = append(groups[i].ErrorLogs, errorLog)
				flg = true
			}
		}
		if !flg {
			group = &ErrorGroup{Key: key}
			group.ErrorLogs = append(group.ErrorLogs, errorLog)
			groups = append(groups, group)
		}
		count++
	}

	for i, v := range groups {
		traces += l.fmtStackTrace(i, len(v.ErrorLogs), v.ErrorLogs[len(v.ErrorLogs)-1])
	}

	if err := scanner.Err(); err != nil {
		l.Logger.Fatal(err)
	}
	return len(groups), traces
}

func (l *prettyLogger) printTraces(count int, traces string) {
	var head string
	if count == 0 {
		return
	}
	head += au.Red("ERROR REPORT\n").Bold().String()
	head += fmt.Sprintf("%v: %v\n", l.attr("ErrorCount"), count)
	head += fmt.Sprintf("%v: %v\n", l.attr("PID"), pid)
	output := fmt.Sprintf("\n\n%s\n\n%s", head, traces)
	if isStdOut {
		if _, err := fmt.Fprint(os.Stdout, output); err != nil {
			return
		}
	} else {
		if _, err := fmt.Fprint(os.Stderr, output); err != nil {
			return
		}
	}
}

func (l *prettyLogger) fmtStackTrace(num, count int, el *ErrorLog) string {
	var output, logFileAbsPath, errorCount string
	logFileAbsPath, err := filepath.Abs(fileName)
	if err != nil {
		return ""
	}

	if count > 1 {
		errorCount = au.Faint(fmt.Sprintf("(%v times)", count)).String()
	}
	output += fmt.Sprintf("%v. %s: %s %s%s%s %v\n",
		au.Bold(num+1),
		filepath.Base(el.Caller),
		l.coloredLevel(el.Severity).String(),
		el.Message,
		separator,
		au.Magenta(el.Error),
		errorCount,
	)
	if el.Timestamp != "" {
		output += fmt.Sprintf("%v:\t%v\n", l.attr("Timestamp"), el.Timestamp)
	}
	output += fmt.Sprintf("%v:\t%v:%v\n",
		l.attr("LogFile"),
		logFileAbsPath,
		el.Line,
	)
	output += fmt.Sprintf("%v: \n\t%v\n\n\n",
		l.attr("StackTrace"),
		strings.ReplaceAll(el.Stacktrace, "\n", "\n\t"),
	)

	return output
}

func (l *prettyLogger) attr(str string) string {
	return "  " + au.Cyan(str).String()
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
