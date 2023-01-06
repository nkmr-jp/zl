package zl

import (
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	FileNameDefault   = "./log/app.jsonl"
	MaxSizeDefault    = 100 // megabytes
	MaxBackupsDefault = 3
	MaxAgeDefault     = 7 // days
)

var (
	fileName   string
	maxSize    int
	maxBackups int
	maxAge     int
	localTime  bool
	compress   bool
)

// newRotator
// See: https://github.com/natefinch/lumberjack
// See: https://github.com/uber-go/zap/blob/master/FAQ.md#does-zap-support-log-rotation
func newRotator() *lumberjack.Logger {
	setRotateDefault()
	res := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		LocalTime:  localTime,
		Compress:   compress,
	}
	return res
}

func setRotateDefault() {
	if fileName == "" {
		fileName = FileNameDefault
	}
	if maxSize == 0 {
		maxSize = MaxSizeDefault
	}
	if maxBackups == 0 {
		maxBackups = MaxBackupsDefault
	}
	if maxAge == 0 {
		maxAge = MaxAgeDefault
	}
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
