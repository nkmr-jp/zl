package zl

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Console(val string) zap.Field {
	return zap.Field{Key: consoleFieldDefault, Type: zapcore.StringType, String: val}
}

func Consolef(format string, a ...interface{}) zap.Field {
	return Console(fmt.Sprintf(format, a...))
}

func Consolep(val *string) zap.Field {
	if val == nil {
		return zap.Reflect(consoleFieldDefault, nil)
	}
	return Console(*val)
}
