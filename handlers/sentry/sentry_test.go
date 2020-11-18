package sentry

import (
	"errors"
	"testing"

	"github.com/apex/log"
)

func TestLogging(t *testing.T) {

	// Without a dsn handler is no-op
	handler, err := NewAsync("")
	if err != nil {
		t.Fatal(err)
	}
	log.SetHandler(handler)
	log.SetLevel(log.DebugLevel)

	ctx := log.WithField("test", "testing")

	ctx.Debug("testing debug message")
	ctx.Info("testing info message")
	ctx.Warn("testing warning message")
	ctx.WithError(errors.New("test error")).Error("testing error message")
	ctx.Errorf("testing formatted error: %s", "error object")

	// Async Handler must be flushed before app termination
	defer handler.Flush(20)
}
