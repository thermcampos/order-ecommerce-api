package logging

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func NewFormatter() *slogFormatter {
	return &slogFormatter{}
}

type LogFormatter interface {
	NewLogEntry(r *http.Request) LogEntry
}

type LogEntry interface {
	Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{})
	Panic(v interface{}, stack []byte)
}

type slogFormatter struct{}

type slogLogEntry struct {
	request *http.Request
	start   time.Time
}

func (f *slogFormatter) NewLogEntry(r *http.Request) middleware.LogEntry {
	return &slogLogEntry{
		request: r,
		start:   time.Now(),
	}
}

func (e *slogLogEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	slog.Info("Request",
		"method", e.request.Method,
		"uri", e.request.RequestURI,
		"remote_addr", e.request.RemoteAddr,
		"status", status,
		"bytes", bytes,
		"duration_ms", elapsed.Milliseconds(),
		"request_id", middleware.GetReqID(e.request.Context()),
	)
}

func (e *slogLogEntry) Panic(v interface{}, stack []byte) {
	slog.Error("panic recovered", "error", v, "stack", string(stack))
}
