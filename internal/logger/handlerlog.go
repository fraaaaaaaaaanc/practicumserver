// Package logger provides a logging interface and an implementation using the Uber Zap logging library.
// The ZapLogger struct allows configuring loggers with the option to write logs to both the console and a file.
// The NewZapLogger function initializes a logger based on the FileLog parameter, which determines whether logs
// should be written to a file in addition to the console. The Info and Error methods are used for logging
// informational and error messages. This package is designed to offer flexibility in logging configurations.
package logger

import (
	"bytes"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
)

type (
	// The responseData structure is used to store data about the server's response.
	responseData struct {
		status int
		size   int
	}
	// The loggingResponseWriter is a wrapper around http.ResponseWriter that allows recording response data.
	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

// The Write method is overridden to record response size.
func (l *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := l.ResponseWriter.Write(b)
	l.responseData.size += size
	return size, err
}

// The WriteHeader method is overridden to record response status.
func (l *loggingResponseWriter) WriteHeader(statusCode int) {
	l.ResponseWriter.WriteHeader(statusCode)
	l.responseData.status = statusCode
}

// MiddlewareLogHandleFunc is a middleware for logging incoming requests.
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

			// Read the request body for further processing.
			body, err := io.ReadAll(r.Body)
			if err != nil {
				// Handling the error of reading the request body.
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// Restore the request body after reading.
			r.Body = io.NopCloser(bytes.NewReader(body))
			h.ServeHTTP(&lrw, r)

			duration := time.Since(start)

			// Create an array of fields to be logged.
			fields := []zap.Field{
				zap.String("uri", r.RequestURI),
				zap.String("method", r.Method),
				zap.Int("status", rd.status),
				zap.Duration("duration", duration),
				zap.Int("size", rd.size),
			}
			// Record request and response information in the log.
			logger.Info("Received request:", fields...)
		})
	}
}
