package logger

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

// LogFormatter initiates the beginning of a new logEntry per request.
type LogFormatter struct{}

var _ middleware.LogFormatter = (*LogFormatter)(nil)

func (l *LogFormatter) NewLogEntry(r *http.Request) middleware.LogEntry {
	entry := &LogEntry{
		request: r,
	}
	return entry
}

// LogEntry records the final log when a request completes.
type LogEntry struct {
	request *http.Request
}

var _ middleware.LogEntry = (*LogEntry)(nil)

func (l *LogEntry) Write(status, bytes int, header http.Header, duration time.Duration, extra interface{}) {
	reqID := middleware.GetReqID(l.request.Context())
	var headers []string

	for k, v := range header {
		for _, val := range v {
			headers = append(headers, fmt.Sprintf("%s: %s", k, val))
		}
	}

	Log.Info(
		"handle request",
		String("event", "handle request"),
		String("requestID", reqID),                // ID запроса
		Strings("headers", headers),               // Заголовки
		String("path", l.request.RequestURI),      // URI
		String("method", l.request.Method),        // Метод запроса
		Int64("latency", duration.Milliseconds()), // Время, затраченное на выполнение запроса (ms)
		Int("status", status),                     // Код статуса ответа
		Int("bytes", bytes),                       // Размер содержимого ответа (в байтах)
	)
}

func (l *LogEntry) Panic(v interface{}, stack []byte) {
	Log.Panic(fmt.Sprintf("%v", v))
}
