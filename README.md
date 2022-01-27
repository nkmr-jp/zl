Refactoring for the release of v1.0

# zl :technologist:
zl provides [zap](https://github.com/uber-go/zap) based advanced logging features, and it's easy to use.

## Install

```sh
go get -u github.com/nkmr-jp/zl
```

```sh
# If you want to use the latest feature.
go get -u github.com/nkmr-jp/zl@develop
```

## Usage

See: [example_test.go](./example_test.go)

```sh
go test -v
```

## Features
- Json structured log to file.
- Simple log to console.
- Stack trace when error.
- Log file rotation.
- Write Code Version and Host to log.
- Write Caller URL to log.
- Context logging.
- etc...
