# Note on `slog` and Ecosystem Changes

As of Go 1.21, the Go standard library includes [log/slog](https://pkg.go.dev/log/slog), a structured logging package. This is now the recommended way to perform structured logging in Go applications.

For Google Cloud (App Engine, Cloud Run, Cloud Functions), writing structured JSON logs to `stdout`/`stderr` is the recommended approach. `slog` with `JSONHandler` supports this natively.

## Resources

*   [Structured Logging with slog](https://go.dev/blog/slog) (Official Go Blog)
*   [log/slog documentation](https://pkg.go.dev/log/slog)
*   [Resources for slog](https://go.dev/wiki/Resources-for-slog) (Community Wiki)
*   [Structured logging in Google Cloud](https://cloud.google.com/logging/docs/structured-logging)

## How to Migrate to `slog`

Instead of using this package, you can use `slog` directly. To make `slog` compatible with Google Cloud Logging's expectations (e.g., using `severity` instead of `level`), you can configure the handler like this:

```go
import (
	"log/slog"
	"os"
)

func init() {
	// Configure slog to use JSON handler and map LevelKey to "severity"
	opts := &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Google Cloud Logging uses "severity" instead of "level"
			if a.Key == slog.LevelKey {
				return slog.Attr{Key: "severity", Value: a.Value}
			}
			return a
		},
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	slog.SetDefault(logger)
}

func handle(w http.ResponseWriter, r *http.Request) {
	// Use slog.Info, slog.Error, etc.
	slog.Info("handling request",
		"method", r.Method,
		"path", r.URL.Path,
        "query", r.URL.RawQuery,
	)
}
```

---

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
