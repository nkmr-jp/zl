Refactoring for the release of v1.0

# Zap Lightning

Zap Lightning is a lightweight wrapper for [zap](https://github.com/uber-go/zap).<br>
It provides presets for easy implementation of advanced logging features.

## How it works

```sh
go test ./zl -v
```

### Write Colored Simple Log to console
![img_1.png](img_1.png)

### Write JSON Structured Log to file

```shell
cat ./zl/log/example.jsonl
```

```json lines
{"level":"DEBUG","caller":"zl/zl.go:44","function":"github.com/nkmr-jp/zap-lightning/zl.Init.func1","message":"INIT_LOGGER","console":"Level: DEBUG, Output: Pretty, FileName: ./log/example.jsonl"}
{"level":"INFO","caller":"zl/example_test.go:39","function":"github.com/nkmr-jp/zap-lightning/zl_test.Example","message":"USER_INFO","user_name":"Alice","user_age":20}
{"level":"ERROR","caller":"zl/example_test.go:41","function":"github.com/nkmr-jp/zap-lightning/zl_test.Example","message":"ERROR_MESSAGE","error":"error message"}
{"level":"DEBUG","caller":"zl/example_test.go:42","function":"github.com/nkmr-jp/zap-lightning/zl_test.Example","message":"DEBUG_MESSAGE"}
{"level":"WARN","caller":"zl/example_test.go:43","function":"github.com/nkmr-jp/zap-lightning/zl_test.Example","message":"WARN_MESSAGE","error":"error message"}
{"level":"WARN","caller":"zl/example_test.go:44","function":"github.com/nkmr-jp/zap-lightning/zl_test.Example","message":"WARN_MESSAGE_WITH_ERROR","error":"error message"}
{"level":"INFO","caller":"zl/example_test.go:45","function":"github.com/nkmr-jp/zap-lightning/zl_test.Example","message":"DISPLAY_TO_CONSOLE","console":"display to console when output type is pretty"}
{"level":"DEBUG","caller":"zl/zl.go:131","function":"github.com/nkmr-jp/zap-lightning/zl.Sync","message":"FLUSH_LOG_BUFFER"}
```

## Install

```sh
go get -u github.com/nkmr-jp/zap-lightning/zl
```

```sh
# If you want to use the latest feature.
go get -u github.com/nkmr-jp/zap-lightning/zl@develop
```

## Usage

See: [example_test.go](./zl/example_test.go)

## Features
- Json structured log to file.
- Simple log to console.
- Stack trace when error.
- Log file rotation.
- Write Code Version and Host to log.
- Write Caller URL to log.
- Context logging.
- etc...
