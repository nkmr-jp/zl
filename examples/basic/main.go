package main

import (
	"fmt"

	"github.com/nkmr-jp/zl"
	"go.uber.org/zap"
)

func main() {
	// Set Options
	zl.SetLevel(zl.DebugLevel)
	zl.SetOmitKeys(zl.HostnameKey)

	// Initialize
	zl.Init()
	defer zl.Sync() // flush log buffer

	// Write logs
	zl.Info("USER_INFO", zap.String("user_name", "Alice"), zap.Int("user_age", 20)) // can use zap fields.
	zl.Info("DISPLAY_TO_CONSOLE", zl.Console("The message you want to display to console"))
	zl.Warn("WARN_MESSAGE")
	zl.Debug("DEBUG_MESSAGE")
	zl.Err("ERROR_MESSAGE1", fmt.Errorf("errors occurring outside the loop"))
	for i := 0; i < 3; i++ {
		zl.Err("ERROR_MESSAGE2", fmt.Errorf("errors occurring in the loop"))
	}

	for i := 0; i < 4; i++ {
		zl.Err("ERROR_MESSAGE3", fmt.Errorf("errors occurring in the loop"))
	}
}
