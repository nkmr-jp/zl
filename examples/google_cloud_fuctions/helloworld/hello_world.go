package helloworld

import (
	"fmt"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/nkmr-jp/zl/examples/google_cloud_fuctions/helloworld/pkg_test/zl"
	"go.uber.org/zap"
)

func init() {
	zl.SetLevel(zl.DebugLevel)
	zl.SetOutput(zl.ConsoleOutput)
	zl.SetOmitKeys(zl.HostnameKey, zl.PIDKey, zl.FunctionKey)

	zl.Init()
	defer zl.SyncWhenStop() // flush log buffer

	functions.HTTP("HelloGet", helloGet)
}

// helloGet is an HTTP Cloud Function.
func helloGet(w http.ResponseWriter, r *http.Request) {
	zl.Info("USER_INFO", zap.String("user_name", "Alice"), zap.Int("user_age", 20)) // can use zap fields.
	err := fmt.Errorf("error message")
	zl.Error("ERROR_MESSAGE", err) // error level log must with error message.
	zl.Debug("DEBUG_MESSAGE")
	zl.Warning("WARN_MESSAGE", zap.Error(err)) // warn level log with error message.

	fmt.Fprint(w, "Hello, World!")
}
