package sentry

import (
	"reflect"
	"time"

	"github.com/getsentry/sentry-go"

	"github.com/apex/log"
)

var (
	levelMap = map[log.Level]sentry.Level{
		log.DebugLevel: sentry.LevelDebug,
		log.InfoLevel: sentry.LevelInfo,
		log.WarnLevel: sentry.LevelWarning,
		log.ErrorLevel: sentry.LevelError,
		log.FatalLevel: sentry.LevelFatal,
		log.InvalidLevel: sentry.LevelFatal,
	}
)

type Handler struct {
	client 	*sentry.Client
	hub		*sentry.Hub
}

func NewAsync(dsn string) (*Handler, error) {
	return newHandler(dsn, sentry.NewHTTPTransport())
}

func NewSync(dsn string) (*Handler, error) {
	return newHandler(dsn, sentry.NewHTTPSyncTransport())

}

func newHandler(dsn string, transport sentry.Transport) (*Handler, error) {

	// TODO: Make debug configurable
	c, err := sentry.NewClient(sentry.ClientOptions{Dsn: dsn, Transport: transport, Debug: false})
	if err != nil {
		return nil, err
	}

	hub := sentry.CurrentHub()
	hub.BindClient(c)

	handler := Handler{
		client: c,
		hub: hub,
	}

	return &handler, nil

}

func (h *Handler) HandleLog(e *log.Entry) error {

	event := sentry.NewEvent()
	event.Message = e.Message
	event.Level = levelMap[e.Level]
	event.Timestamp = e.Timestamp

	err := e.Err()

	for k, v := range e.Fields {
		event.Extra[k] = v
	}

	if err != nil {
		stacktrace := sentry.ExtractStacktrace(err)
		event.Exception = []sentry.Exception{{
			Value: err.Error(),
			Type: reflect.TypeOf(err).String(),
			Stacktrace: stacktrace,
		}}
	}

	h.hub.CaptureEvent(event)
	return nil
}

func (h *Handler) Flush(seconds time.Duration) bool {
	return h.hub.Flush(seconds * time.Second)
}
