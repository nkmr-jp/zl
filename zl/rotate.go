package zl

import (
	"log"

	"gopkg.in/natefinch/lumberjack.v2"
)

// See: https://github.com/natefinch/lumberjack
const (
	fileNameDefault   = "./log/app.jsonl"
	maxSizeDefault    = 100 // megabytes
	maxBackupsDefault = 3
	maxAgeDefault     = 7 // days
)

var (
	fileName string
)

// See:
// https://github.com/natefinch/lumberjack
// https://github.com/uber-go/zap/blob/master/FAQ.md#does-zap-support-log-rotation
func newLumberjack() *lumberjack.Logger {
	setRotateDefault()
	res := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    maxSizeDefault,
		MaxBackups: maxBackupsDefault,
		MaxAge:     maxAgeDefault,
	}
	log.Printf("log file path: %v", fileName)
	return res
}

func setRotateDefault() {
	if fileName == "" {
		fileName = fileNameDefault
	}
	// if rotationTime == 0 {
	// 	rotationTime = rotationTimeDefault
	// }
}
