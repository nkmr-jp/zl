package zl

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Console(val string) zap.Field {
	return zap.Field{Key: "console", Type: zapcore.StringType, String: val}
}

func Consolep(val *string) zap.Field {
	if val == nil {
		return zap.Reflect("console", nil)
	}
	return Console(*val)
}
