package zl

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"testing"
)

func TestConsole(t *testing.T) {
	field := Console("test")
	expected := zap.Field{Key: consoleFieldDefault, Type: zapcore.StringType, String: "test"}
	assert.Equal(t, expected, field)
}

func TestConsolep(t *testing.T) {
	val := "pointerTest"
	field := Consolep(&val)
	expected := zap.Field{Key: consoleFieldDefault, Type: zapcore.StringType, String: "pointerTest"}
	assert.Equal(t, expected, field)

	nilField := Consolep(nil)
	expectedNil := zap.Field{Key: consoleFieldDefault, Type: zapcore.StringType, String: ""}
	assert.Equal(t, expectedNil, nilField)
}

func TestConsolef(t *testing.T) {
	format := "Hello %s"
	name := "World"
	field := Consolef(format, name)
	expected := zap.Field{Key: consoleFieldDefault, Type: zapcore.StringType, String: "Hello World"}
	assert.Equal(t, expected, field)
}
