// Package compress provides middleware for compressing HTTP request and response data using GZIP encoding.
package compress

import (
	"compress/gzip"
	"go.uber.org/zap"
	"io"
	"net/http"
)

// Structure for storing regular and compressed response in a handler.
type compressWriter struct {
	w  http.ResponseWriter
	gz *gzip.Writer
}

// Constructor for the compressWriter structure.
func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		gz: gzip.NewWriter(w),
	}
}

// Overridden methods: Header, Write, WriteHeader, Close.
func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

// Overridden Write method checks the response content type to determine whether it should be compressed.
func (c *compressWriter) Write(b []byte) (int, error) {
	if len(b) < 1 ||
		(c.Header().Get("Content-Type") != "text/plain" &&
			c.Header().Get("Content-Type") != "application/json") {
		return c.w.Write(b)
	}
	return c.gz.Write(b)
}

// Overridden WriteHeader method
func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

// Overridden Close method
func (c *compressWriter) Close() error {
	return c.gz.Close()
}

// Structure for storing regular and compressed request.
type compressReader struct {
	r  io.ReadCloser
	rz *gzip.Reader
}

// Constructor for the compressReader structure.
func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	rz, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		rz: rz,
	}, nil
}

// Overridden Read and Close methods.
func (c *compressReader) Read(b []byte) (int, error) {
	return c.rz.Read(b)
}

func (c *compressReader) Close() error {
	return c.rz.Close()
}

// Middleware function that checks incoming data for compression and client's acceptance of compressed data.
// Based on this, it calls the constructors of the compressReader and compressWriter structures.
func MiddlewareGzipHandleFunc(logger *zap.Logger) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var cw *compressWriter
			var cr *compressReader
			ow := w

			zipFormAccept := r.Header.Values("Accept-Encoding")
			for _, elem := range zipFormAccept {
				if elem == "gzip" {
					cw = newCompressWriter(w)
					ow = cw
				}
			}
			defer func() {
				if cw != nil {
					err := cw.Close()
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						logger.Error("Error:", zap.Error(err))
						return
					}
				}
			}()

			zipFormContent := r.Header.Values("Content-Encoding")
			for _, elem := range zipFormContent {
				if elem == "gzip" {
					cr, err := newCompressReader(r.Body)
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						logger.Error("Error:", zap.Error(err))
						return
					}
					r.Body = cr
				}
			}
			defer func() {
				if cr != nil {
					err := cr.Close()
					if err != nil {
						w.WriteHeader(http.StatusInternalServerError)
						logger.Error("Error:", zap.Error(err))
						return
					}
				}
			}()
			h.ServeHTTP(ow, r)
		})
	}
}
