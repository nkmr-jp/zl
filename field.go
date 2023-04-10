package zl

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Console is display to console when output type is pretty.
func Console(val string) zap.Field {
	return zap.Field{Key: consoleFieldDefault, Type: zapcore.StringType, String: val}
}

// Consolep is display to console when output type is pretty.
func Consolep(val *string) zap.Field {
	if val == nil {
		return Console("")
	}
	return Console(*val)
}

// Consolef formats according to a format specifier and display to console when output type is pretty.
func Consolef(format string, a ...interface{}) zap.Field {
	return Console(fmt.Sprintf(format, a...))
}
