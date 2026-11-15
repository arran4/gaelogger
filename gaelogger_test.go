package gaelogger

import (
  "bytes"
  "log"
  "strings"
  "testing"
)

func TestNewLoggerWithoutGoogleCloudProject(t *testing.T) {
  t.Setenv("GOOGLE_CLOUD_PROJECT", "")
  logger := NewLogger(nil)
  wrapper, ok := logger.(*WrapLogger)
  if !ok {
    t.Fatalf("NewLogger returned %T, want *WrapLogger", logger)
  }
  if _, ok := wrapper.backLogger.(*StdLogger); !ok {
    t.Fatalf("fallback backend is %T, want *StdLogger", wrapper.backLogger)
  }
  if err := logger.Close(); err != nil {
    t.Fatalf("Close() = %v, want nil", err)
  }
}

func TestStdLoggerFormatsSeverityAndArguments(t *testing.T) {
  var output bytes.Buffer
  logger := &WrapLogger{backLogger: &StdLogger{logger: log.New(&output, "", 0)}}
  cases := []struct {
    name string
    write func()
  }{
    {"DEFAULT", func() { logger.Defaultf("message %d", 42) }},
    {"DEBUG", func() { logger.Debugf("message %d", 42) }},
    {"INFO", func() { logger.Infof("message %d", 42) }},
    {"NOTICE", func() { logger.Noticef("message %d", 42) }},
    {"WARNING", func() { logger.Warningf("message %d", 42) }},
    {"ERROR", func() { logger.Errorf("message %d", 42) }},
    {"CRITICAL", func() { logger.Criticalf("message %d", 42) }},
    {"ALERT", func() { logger.Alertf("message %d", 42) }},
    {"EMERGENCY", func() { logger.Emergencyf("message %d", 42) }},
  }
  for _, tt := range cases {
    t.Run(tt.name, func(t *testing.T) {
      output.Reset()
      tt.write()
      got := output.String()
      want := tt.name[:1] + strings.ToLower(tt.name[1:]) + " message 42"
      if !strings.Contains(got, want) {
        t.Errorf("output %q does not contain expected severity and formatted message", got)
      }
    })
  }
  if err := logger.Close(); err != nil {
    t.Fatalf("Close() = %v, want nil", err)
  }
}

func TestGaeLoggerCloseWithoutClient(t *testing.T) {
  if err := (&GaeLogger{}).Close(); err != nil {
    t.Fatalf("Close() = %v, want nil", err)
  }
}
