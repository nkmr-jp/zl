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
