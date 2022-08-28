package main

import (
	"encoding/json"
	"os"
	"strconv"

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
	_, err := os.ReadFile("test")
	zl.Err("READ_FILE_ERROR", err)
	for i := 0; i < 2; i++ {
		_, err = strconv.Atoi("one")
		zl.Err("A_TO_I_ERROR", err)
	}
	for i := 0; i < 3; i++ {
		v := ""
		err = json.Unmarshal([]byte("test"), &v)
		zl.Err("JSON_UNMARSHAL_ERROR", err)
	}
}
