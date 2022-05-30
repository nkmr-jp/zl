# zl :technologist:
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

```sh
# If you want to use the latest feature.
go get -u github.com/nkmr-jp/zl@develop
```

# Example

- [examples](examples)
- [example_test.go](example_test.go)
