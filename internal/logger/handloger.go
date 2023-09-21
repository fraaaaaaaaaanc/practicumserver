package logger

import (
	"net/http"
	"time"
)

type (
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (l *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := l.ResponseWriter.Write(b)
	l.responseData.size += size
	return size, err
}

func (l *loggingResponseWriter) WriteHeader(statusCode int) {
	l.ResponseWriter.WriteHeader(statusCode)
	l.responseData.status = statusCode
}

func WithLogging(h http.Handler) http.Handler {
	logfn := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		responseData := &responseData{
			size:   0,
			status: 0,
		}

		lw := loggingResponseWriter{
			ResponseWriter: w,
			responseData:   responseData,
		}

		h.ServeHTTP(&lw, r)

		duration := time.Since(start)

		log.Infoln(
			"request info",
			"uri", r.RequestURI,
			"method", r.Method,
			"status", responseData.status,
			"duration", duration,
			"size", responseData.size,
		)
	}
	return http.HandlerFunc(logfn)
}
