package zl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_newRotator(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		ret := newRotator()
		assert.Equal(t, "./log/app.jsonl", ret.Filename)
		assert.Equal(t, 100, ret.MaxSize)
		assert.Equal(t, 3, ret.MaxBackups)
		assert.Equal(t, 7, ret.MaxAge)
		assert.Equal(t, false, ret.LocalTime)
		assert.Equal(t, false, ret.Compress)

	})
	t.Run("set options", func(t *testing.T) {
		SetRotateFileName("./log/Test_newRotator.jsonl")
		SetRotateMaxSize(1000)
		SetRotateMaxBackups(5)
		SetRotateMaxAge(14)
		SetRotateLocalTime(true)
		SetRotateCompress(true)
		ret := newRotator()
		assert.Equal(t, "./log/Test_newRotator.jsonl", ret.Filename)
		assert.Equal(t, 1000, ret.MaxSize)
		assert.Equal(t, 5, ret.MaxBackups)
		assert.Equal(t, 14, ret.MaxAge)
		assert.Equal(t, true, ret.LocalTime)
		assert.Equal(t, true, ret.Compress)
	})

}
