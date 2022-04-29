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