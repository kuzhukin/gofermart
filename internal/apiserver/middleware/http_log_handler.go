package middleware

import (
	"gophermart/internal/zlog"
	"net/http"
	"time"
)

var _ http.ResponseWriter = &loggingResponseWriter{}

type loggingResponseWriter struct {
	http.ResponseWriter
	size   int
	status int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{
		ResponseWriter: w,
	}
}

func (l *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := l.ResponseWriter.Write(b)
	l.size += size
	return size, err
}

func (l *loggingResponseWriter) WriteHeader(status int) {
	l.ResponseWriter.WriteHeader(status)
	l.status = status
}

func (l *loggingResponseWriter) doRequestWithTimer(h http.Handler, r *http.Request) time.Duration {
	start := time.Now()

	h.ServeHTTP(l.ResponseWriter, r)

	return time.Since(start)
}

func LoggingHTTPHandler(h http.Handler) http.Handler {
	loggingHandler := func(w http.ResponseWriter, r *http.Request) {
		lw := newLoggingResponseWriter(w)
		duration := lw.doRequestWithTimer(h, r)

		zlog.Logger.Infof("uri=%v, method=%v, status=%v, size=%v, duration=%v",
			r.RequestURI, r.Method, lw.status, lw.size, duration)
	}

	return http.HandlerFunc(loggingHandler)
}
