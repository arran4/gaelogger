package gaelogger

import (
	"cloud.google.com/go/logging"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
)

type backLogger interface {
	io.Closer
	logf(severity fmt.Stringer, format string, args ...interface{})
}

type Logger interface {
	io.Closer
	Defaultf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Noticef(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Criticalf(format string, args ...interface{})
	Alertf(format string, args ...interface{})
	Emergencyf(format string, args ...interface{})
}

func NewLogger(r *http.Request) Logger {
	if len(os.Getenv("GOOGLE_CLOUD_PROJECT")) > 0 {
		var ctx context.Context
		if r != nil {
			ctx = r.Context()
		} else {
			ctx = context.Background()
		}
		loggingClient, err := logging.NewClient(ctx, os.Getenv("GOOGLE_CLOUD_PROJECT"))
		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
		}
		logger := loggingClient.Logger("request")
		return &WrapLogger{
			backLogger: &GaeLogger{
				loggingClient: loggingClient,
				logger:        logger,
				r:             r,
			},
		}
	} else {
		return &WrapLogger{
			backLogger: &StdLogger{
				logger: log.New(os.Stderr, "", log.Lshortfile|log.LstdFlags),
			},
		}
	}
}

type StdLogger struct {
	logger *log.Logger
}

func (st *StdLogger) Close() error {
	return nil
}

func (st *StdLogger) logf(severity fmt.Stringer, format string, args ...interface{}) {
	st.logger.SetPrefix(severity.String()+" ")
	st.logger.Output(3, fmt.Sprintf(format, args...))
	st.logger.SetPrefix("")
}

type GaeLogger struct {
	loggingClient *logging.Client
	logger        *logging.Logger
	r             *http.Request
}

func (gl *GaeLogger) Close() error {
	if gl.loggingClient != nil {
		return gl.loggingClient.Close()
	}
	return nil
}

func (gl *GaeLogger) logf(severity fmt.Stringer, format string, args ...interface{}) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	gl.logger.Log(logging.Entry{
		Severity: severity.(logging.Severity),
		Payload:  fmt.Sprintf("%s:%d "+format, append([]interface{}{file, line}, args...)...),
		HTTPRequest: &logging.HTTPRequest{
			Request: gl.r,
		},
	})
}

type WrapLogger struct {
	backLogger
}

func (gl WrapLogger) Defaultf(format string, args ...interface{}) {
	gl.logf(logging.Default, format, args...)
}

func (gl *WrapLogger) Debugf(format string, args ...interface{}) {
	gl.logf(logging.Debug, format, args...)
}

func (gl *WrapLogger) Infof(format string, args ...interface{}) {
	gl.logf(logging.Info, format, args...)
}

func (gl *WrapLogger) Noticef(format string, args ...interface{}) {
	gl.logf(logging.Notice, format, args...)
}

func (gl *WrapLogger) Warningf(format string, args ...interface{}) {
	gl.logf(logging.Warning, format, args...)
}

func (gl *WrapLogger) Errorf(format string, args ...interface{}) {
	gl.logf(logging.Error, format, args...)
}

func (gl *WrapLogger) Criticalf(format string, args ...interface{}) {
	gl.logf(logging.Critical, format, args...)
}

func (gl *WrapLogger) Alertf(format string, args ...interface{}) {
	gl.logf(logging.Alert, format, args...)
}

func (gl *WrapLogger) Emergencyf(format string, args ...interface{}) {
	gl.logf(logging.Emergency, format, args...)
}
