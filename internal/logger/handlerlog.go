package logger

import (
	"bytes"
	"go.uber.org/zap"
	"io"
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

func MiddlewareLogHandleFunc(logger *zap.Logger) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			rd := &responseData{
				size:   0,
				status: 0,
			}

			lrw := loggingResponseWriter{
				ResponseWriter: w,
				responseData:   rd,
			}

			body, err := io.ReadAll(r.Body)
			if err != nil {
				// Обработка ошибки чтения тела запроса
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			r.Body = io.NopCloser(bytes.NewReader(body))
			h.ServeHTTP(&lrw, r)

			duration := time.Since(start)

			fields := []zap.Field{
				zap.String("uri", r.RequestURI),
				zap.String("method", r.Method),
				zap.Int("status", rd.status),
				zap.Duration("duration", duration),
				zap.Int("size", rd.size),
			}
			logger.Info("Received request:", fields...)
		})
	}
}
