package gaelogger

import (
	"cloud.google.com/go/logging"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
)

type GaeLogger struct {
	loggingClient *logging.Client
	logger        *logging.Logger
	r             *http.Request
}

func NewGaeLogger(r *http.Request) *GaeLogger {
	if len(os.Getenv("GOOGLE_CLOUD_PROJECT")) > 0 {
		loggingClient, err := logging.NewClient(r.Context(), os.Getenv("GOOGLE_CLOUD_PROJECT"))
		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
		}
		logger := loggingClient.Logger("request")
		return &GaeLogger{
			loggingClient: loggingClient,
			logger:        logger,
			r:             r,
		}
	} else {
		return &GaeLogger{}
	}
}

func (gl *GaeLogger) Close() error {
	if gl.loggingClient != nil {
		return gl.loggingClient.Close()
	}
	return nil
}

func (gl *GaeLogger) logf(severity logging.Severity, format string, args ...interface{}) {
	if gl.logger != nil {
		_, file, line, ok := runtime.Caller(2)
		if !ok {
			file = "???"
			line = 0
		}
		gl.logger.Log(logging.Entry{
			Severity: severity,
			Payload:  fmt.Sprintf("%s:%d "+format, append([]interface{}{file, line}, args...)...),
			HTTPRequest: &logging.HTTPRequest{
				Request: gl.r,
			},
		})
	} else {
		log.Printf("Error: "+format, args...)
	}
}

func (gl *GaeLogger) Defaultf(format string, args ...interface{}) {
	gl.logf(logging.Default, format, args...)
}

func (gl *GaeLogger) Debugf(format string, args ...interface{}) {
	gl.logf(logging.Debug, format, args...)
}

func (gl *GaeLogger) Infof(format string, args ...interface{}) {
	gl.logf(logging.Info, format, args...)
}

func (gl *GaeLogger) Noticef(format string, args ...interface{}) {
	gl.logf(logging.Notice, format, args...)
}

func (gl *GaeLogger) Warningf(format string, args ...interface{}) {
	gl.logf(logging.Warning, format, args...)
}

func (gl *GaeLogger) Errorf(format string, args ...interface{}) {
	gl.logf(logging.Error, format, args...)
}

func (gl *GaeLogger) Criticalf(format string, args ...interface{}) {
	gl.logf(logging.Critical, format, args...)
}

func (gl *GaeLogger) Alertf(format string, args ...interface{}) {
	gl.logf(logging.Alert, format, args...)
}

func (gl *GaeLogger) Emergencyf(format string, args ...interface{}) {
	gl.logf(logging.Emergency, format, args...)
}
