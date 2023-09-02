# example/httpserver

This is a sample for configuring logs on http server. A trace ID is issued for each request so that the processing of each request can be traced.
Also measure the processing time for each method and add it to the log.


Console Output (Colored simple text)
<img width="908" alt="image" src="https://github.com/nkmr-jp/zl/assets/8490118/bdcc41c1-c08e-49aa-bd7e-eda566af8c8e">

File Output (JSON)
<img width="2021" alt="image" src="https://github.com/nkmr-jp/zl/assets/8490118/0521b531-0f1b-44a1-b19b-c3ada4d64339">


# Usage
## start http server 
```sh
ENV=local PORT=8080 go run *.go
```

## call http server
```sh
curl http://localhost:8080
```



