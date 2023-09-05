# Google Cloud Functions Example

## Test in mac local

start functions framework server.
```sh
make start
```

send request.
```sh
curl http://localhost:8080
# > Hello, World!
```

## Deploy to Google Cloud

### Before you begin

See: https://cloud.google.com/functions/docs/2nd-gen/getting-started#before-you-begin

### Deploy

deploy to google cloud.
```sh
make deploy
```

describe functions detail.
```sh
make show
```

show functions log. 
```sh
make log
```

open functions detail in google cloud console.
```sh
make open
```

## Severity Level

In google cloud logging, warn is converted to warning and fatal is converted to critical.

See: https://cloud.google.com/logging/docs/reference/v2/rest/v2/LogEntry#logseverity


| zl    | google cloud logging | description                                                 |
|-------|----------------------|-------------------------------------------------------------|
| DEBUG | DEBUG                | Debug or trace information.                                 |
| INFO  | INFO                 | Routine information, such as ongoing status or performance. |
| WARN  | WARNING              | Warning events might cause problems.                        |
| ERROR | ERROR                | Error events are likely to cause problems.                  |
| FATAL | CRITICAL             | Critical events cause more severe problems or outages.      |



<img width="1015" alt="image" src="https://github.com/nkmr-jp/zl/assets/8490118/8a80e40e-c572-4f40-be3f-82145b96c387">