# zl :technologist:
[![Go Reference](https://pkg.go.dev/badge/github.com/nkmr-jp/zl.svg)](https://pkg.go.dev/github.com/nkmr-jp/zl)
[![test](https://github.com/nkmr-jp/zl/actions/workflows/test.yml/badge.svg)](https://github.com/nkmr-jp/zl/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/nkmr-jp/zl)](https://goreportcard.com/report/github.com/nkmr-jp/zl)
[![codecov](https://codecov.io/gh/nkmr-jp/zl/branch/main/graph/badge.svg?token=2Z6M2JYT17)](https://codecov.io/gh/nkmr-jp/zl)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)

zl is a logger based on [zap](https://github.com/uber-go/zap). It provides advanced logging features.

zl is a logging package designed with the developer experience in mind.
You can choose the most suitable output format according to your purpose, such as emphasizing easy-to-read console output during development in a local environment and outputting structured detailed logs in a production environment.
It offers rich functionality but is easy to configure.

This is useful when developing systems that perform complex processing.

# Features
## Selectable output format
### PrettyOutput (Default) :technologist: 
- High Developer Experience.
- The optimal setting for a development environment.
- Output colored simple logs to the console.
- Output detail JSON logs to the logfile.
- Easy-to-read error reports and stack trace.
- It can jumps directly to the line of the file that is output to the console log (when using Goland or VSCode).

![image](https://user-images.githubusercontent.com/8490118/185822142-4667200b-8087-49f0-9e41-68ebb1731985.png)

Log messages do not have to be written in upper snake case, 
but if messages are written in a uniform manner, 
it is easy to extract specific logs using tools such as [Google Cloud Logging](https://cloud.google.com/logging) or the [jq command](https://stedolan.github.io/jq/).
It is recommended that log messages be written in a consistent manner.
By writing log messages succinctly and placing detailed information in separate fields, 
the overall clarity of the logs is improved, making it easier to understand the flow of processes.

### ConsoleOutput :zap:
- High Performance.
- The optimal setting for a production environment.
- Output detail JSON logs to the console.
- Especially suitable for cloud environments such as [Google Cloud Logging](https://cloud.google.com/logging) or [Datadog](https://www.datadoghq.com/).
- It is fast because it uses only the functions provided by [zap](https://github.com/uber-go/zap#performance) (**not sugared**) and does not perform any extra processing.

### FileOutput
- The optimal setting for a production environment.
- It is especially suitable for command line tool or on-premises environments.
- Support logfile rotation.

### ConsoleAndFileOutput
- It is a setting for the development environment.
- Output detail JSON logs to console and logfile.
- It is recommended to use with [jq command](https://stedolan.github.io/jq/) to avoid drowning in a sea of information.
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
	
	// if the same error occurs multiple times in the same location, the error report will show them all together.
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
```

**Console output**. <br>
The console displays minimal information. Displays a stack trace when an error occurs.
![image](https://user-images.githubusercontent.com/8490118/185822142-4667200b-8087-49f0-9e41-68ebb1731985.png)

**File output**. <br>
Detailed information is available in the log file. You can also use jq to extract only the information you need.
```sh
$ cat log/app.jsonl | jq 'select(.message | startswith("USER_")) | select(.pid==921)'
```
```json
{
  "severity": "INFO",
  "timestamp": "2022-08-22T10:30:39.154575+09:00",
  "caller": "basic/main.go:22",
  "function": "main.main",
  "message": "USER_INFO",
  "version": "fc43f68",
  "pid": 921,
  "user_name": "Alice",
  "user_age": 20
}
   
```

# Examples
- [examples](examples)
- [example_test.go](example_test.go)
