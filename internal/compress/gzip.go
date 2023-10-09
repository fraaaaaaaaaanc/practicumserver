package compress

import (
	"compress/gzip"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type compressWriter struct {
	w  http.ResponseWriter
	gz *gzip.Writer
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	return &compressWriter{
		w:  w,
		gz: gzip.NewWriter(w),
	}
}

func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

func (c *compressWriter) Write(b []byte) (int, error) {
	if len(b) < 1 ||
		(c.Header().Get("Content-Type") != "text/plain" &&
			c.Header().Get("Content-Type") != "application/json") {
		return c.w.Write(b)
	}
	return c.gz.Write(b)
}

func (c *compressWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

func (c *compressWriter) Close() error {
	return c.gz.Close()
}

type compressReader struct {
	r  io.ReadCloser
	rz *gzip.Reader
}

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

func (c *compressReader) Read(b []byte) (int, error) {
	return c.rz.Read(b)
}

func (c *compressReader) Close() error {
	return c.rz.Close()
}

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
						fields := []zap.Field{
							zap.Error(err),
						}
						logger.Error("Error:", fields...)
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
						fields := []zap.Field{
							zap.Error(err),
						}
						logger.Error("Error:", fields...)
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
						fields := []zap.Field{
							zap.Error(err),
						}
						logger.Error("Error:", fields...)
						return
					}
				}
			}()
			h.ServeHTTP(ow, r)
		})
	}
}

//func MiddlewareGzipHandleFunc(h http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		var cw *compressWriter
//		var cr *compressReader
//		ow := w
//
//		zipFormAccept := r.Header.Values("Accept-Encoding")
//		for _, elem := range zipFormAccept {
//			if elem == "gzip" {
//				cw = newCompressWriter(w)
//				ow = cw
//			}
//		}
//		defer func() {
//			if cw != nil {
//				err := cw.Close()
//				if err != nil {
//					w.WriteHeader(http.StatusInternalServerError)
//					return
//				}
//			}
//		}()
//
//		zipFormContent := r.Header.Values("Content-Encoding")
//		for _, elem := range zipFormContent {
//			if elem == "gzip" {
//				cr, err := newCompressReader(r.Body)
//				if err != nil {
//					w.WriteHeader(http.StatusInternalServerError)
//					return
//				}
//				r.Body = cr
//			}
//		}
//		defer func() {
//			if cr != nil {
//				err := cr.Close()
//				if err != nil {
//					w.WriteHeader(http.StatusInternalServerError)
//					return
//				}
//			}
//		}()
//		h.ServeHTTP(ow, r)
//	})
//}
