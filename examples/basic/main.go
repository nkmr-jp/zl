package main

import (
	"github.com/nkmr-jp/zl"
	"go.uber.org/zap"
)

func main() {

	// Set Options
	zl.SetLevel(zl.DebugLevel)

	// Initialize
	zl.Init()
	defer zl.Sync() // flush log buffer

	// Write logs
	zl.Info("USER_INFO", zap.String("user_name", "Alice"), zap.Int("user_age", 20)) // can use zap fields.
	zl.Warn("WARN_MESSAGE")
	zl.Debug("DEBUG_MESSAGE")
}
