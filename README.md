# zl :technologist:
[![Go Reference](https://pkg.go.dev/badge/github.com/nkmr-jp/zl.svg)](https://pkg.go.dev/github.com/nkmr-jp/zl)
[![test](https://github.com/nkmr-jp/zl/actions/workflows/test.yml/badge.svg)](https://github.com/nkmr-jp/zl/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/nkmr-jp/zl)](https://goreportcard.com/report/github.com/nkmr-jp/zl)
[![codecov](https://codecov.io/gh/nkmr-jp/zl/branch/develop/graph/badge.svg?token=2Z6M2JYT17)](https://codecov.io/gh/nkmr-jp/zl)

zl provides [zap-based](https://github.com/uber-go/zap) advanced logging features.

Its design focuses on the developer experience and is easy to use. 

# Features
## Selectable output types
### PrettyOutput (Default) :technologist: 
- High Developer Experience.
- The optimal setting for a development environment.
- Output colored simple logs to the console.
- Output detail JSON logs to logfile.
- Very easy-to-read error reports and stack trace.
- It can jumps directly to the line of the file that is output to the console log (when using Goland or VSCode).

### ConsoleOutput :zap:
- High Performance.
- The optimal setting for a production environment.
- Especially suitable for cloud environments such as [Google Cloud Logging](https://cloud.google.com/logging) or [Datadog](https://www.datadoghq.com/).
- Only uses the features provided by [zap](https://github.com/uber-go/zap#performance) (**not sugared**).

### FileOutput
- The optimal setting for a production environment.
- Especially suitable for on-premises environments.
- Support logfile rotation.

### ConsoleAndFileOutput
- It is a setting for the development environment.
- Output detail JSON logs to console and logfile.
- It is recommended to use with [jq](https://stedolan.github.io/jq/) to avoid drowning in a sea of information.
- It is recommended to set PrettyOutput instead.


# Installation

```sh
go get -u github.com/nkmr-jp/zl
```

# Quick Start

code: [examples/basic/main.go](examples/basic/main.go)
```go
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
	zl.Error("ERROR_MESSAGE", fmt.Errorf("some error occurred"))
}
```

console output: <br>
<img width="100%" src="https://user-images.githubusercontent.com/8490118/173165186-74b001e4-80f1-4573-b99c-a94445760360.png" />


file output:
```sh
$ cat log/app.jsonl | jq 'select(.message | startswith("USER_")) | select(.pid==17925)'
{
  "severity": "INFO",
  "timestamp": "2022-06-11T09:24:01.14941+09:00",
  "caller": "basic/main.go:20",
  "function": "main.main",
  "message": "USER_INFO",
  "version": "5bb45d7",
  "pid": 17925,
  "user_name": "Alice",
  "user_age": 20
}   
```

# Examples
- [examples](examples)
- [example_test.go](example_test.go)
