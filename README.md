# zl :technologist:
zl provides [zap-based](https://github.com/uber-go/zap) advanced logging features.
Its design focuses on the developer experience and is easy to use. 

# Features
## Selectable output types
### PrettyOutput (Default) :technologist: 
- High Developer Experience.
- The optimal setting for a development environment.
- Simple colored log to console.
- Detail JSON log to logfile.
- Very easy-to-read error reports and stack trace.

### ConsoleOutput :zap:
- High Performance.
- The optimal setting for a production environment.
- Especially suitable for cloud environments such as [Google Cloud Logging](https://cloud.google.com/logging) or [Datadog](https://www.datadoghq.com/).
- Only uses the features provided by [zap](https://github.com/uber-go/zap#performance) (**not sugared**).

### FileOutput
- The optimal setting for a production environment.
- Especially suitable for on-premises environments.
- Support logfile rotation.


# Installation

```sh
go get -u github.com/nkmr-jp/zl
```

```sh
# If you want to use the latest feature.
go get -u github.com/nkmr-jp/zl@develop
```

# Example

- [examples](examples)
- [example_test.go](example_test.go)