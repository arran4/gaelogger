# Another google apps engine go logger this one is for the new go111+ world..

Usage is simple:
```go
import (
  "github.com/arran4/gaelogger"
  "net/http"
)

func init() {
    http.HandleFunc("/", handle)
}

func handle(writer http.ResponseWriter, request *http.Request) {
  logger := gaelogger.NewLogger(request)
  defer logger.Close()
  logger.Infof("%s request: %v %v", request.Method, request.URL.RawPath, request.URL.RawQuery)
}
```
