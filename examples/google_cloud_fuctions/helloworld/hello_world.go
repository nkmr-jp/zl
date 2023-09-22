package helloworld

import (
	"fmt"
	"github.com/nkmr-jp/zl"
	"net/http"
	"os"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"go.uber.org/zap"
)

func init() {
	if os.Getenv("ENV") == "local" {
		zl.SetOutput(zl.PrettyOutput)
	} else {
		zl.SetOutput(zl.ConsoleOutput)
	}
	zl.SetLevel(zl.DebugLevel)
	zl.SetOmitKeys(zl.HostnameKey, zl.PIDKey, zl.FunctionKey)
	zl.Init()

	functions.HTTP("HelloWorld", helloWorld)
}

// helloWorld is an HTTP Cloud Function.
func helloWorld(w http.ResponseWriter, r *http.Request) {
	zl.Info("USER_INFO", zap.String("user_name", "Alice"), zap.Int("user_age", 20)) // can use zap fields.
	err := fmt.Errorf("error message")
	zl.Err("ERROR_MESSAGE", err)
	zl.Debug("DEBUG_MESSAGE")
	zl.WarnErr("WARN_MESSAGE", err) // note: WARNING LEVEL in Google Cloud Logging
	//zl.FatalErr("FATAL_MESSAGE", err) // note: CRITICAL LEVEL in Google Cloud Logging

	fmt.Fprint(w, "Hello, World!")
}
