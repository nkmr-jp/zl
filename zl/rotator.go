package zl

import (
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	fileNameDefault   = "./log/app.jsonl"
	maxSizeDefault    = 100 // megabytes
	maxBackupsDefault = 3
	maxAgeDefault     = 7 // days
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
		fileName = fileNameDefault
	}
	if maxSize == 0 {
		maxSize = maxSizeDefault
	}
	if maxBackups == 0 {
		maxBackups = maxBackupsDefault
	}
	if maxAge == 0 {
		maxAge = maxAgeDefault
	}
}
