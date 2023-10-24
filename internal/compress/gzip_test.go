package compress

//
//import (
//	"bytes"
//	"compress/gzip"
//	"context"
//	"fmt"
//	"github.com/stretchr/testify/assert"
//	"io"
//	"net/http"
//	"net/http/httptest"
//	"practicumserver/internal/handlers/allhandlers"
//	"practicumserver/internal/logger"
//	"practicumserver/internal/models"
//	"practicumserver/internal/storage/pg"
//	"testing"
//)
//
//type HandlerFuncAdapter func(http.ResponseWriter, *http.Request)
//
//func (h HandlerFuncAdapter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	h(w, r)
//}
//
//func TestMiddlewareGzipHandleFunc(t *testing.T) {
//	log, _ := logger.NewZapLogger(false)
//	strg, _ := storage.NewStorage(log.Logger, "", "")
//	hndlrs := handlers.NewHandlers(strg, log.Logger, "http://localhost:8080")
//
//	adapter := HandlerFuncAdapter(hndlrs.PostRequestAPIShorten)
//	newHandler := MiddlewareGzipHandleFunc(nil)
//	handler := newHandler(adapter)
//
//	srv := httptest.NewServer(handler)
//	defer srv.Close()
//
//	requestBody := `{"url": "http://test.com"}`
//
//	t.Run("send_gzip", func(t *testing.T) {
//		buf := bytes.NewBuffer(nil)
//		zb := gzip.NewWriter(buf)
//		_, err := zb.Write([]byte(requestBody))
//		assert.NoError(t, err)
//		err = zb.Close()
//		assert.NoError(t, err)
//
//		request := httptest.NewRequest(http.MethodPost, srv.URL+"/api/shorten", buf)
//		assert.NoError(t, err)
//		request.Header.Set("Content-Type", "application/json; charset=utf-8")
//		request.Header.Set("Content-Encoding", "gzip")
//		request.RequestURI = ""
//		ctx := context.WithValue(context.Background(), models.UserIDKey, "1234")
//		request = request.WithContext(ctx)
//
//		resp, err := http.DefaultClient.Do(request)
//		assert.NoError(t, err)
//		assert.Equal(t, http.StatusCreated, resp.StatusCode)
//
//		defer resp.Body.Close()
//
//		b, err := io.ReadAll(resp.Body)
//		assert.NoError(t, err)
//		assert.NotNil(t, string(b))
//	})
//
//	t.Run("accept_gzip", func(t *testing.T) {
//		buf := bytes.NewBufferString(requestBody)
//		request := httptest.NewRequest("POST", srv.URL+"/api/shorten", buf)
//		request.RequestURI = ""
//		request.Header.Set("Accept-Encoding", "gzip")
//		request.Header.Set("Content-Type", "application/json; charset=utf-8")
//		ctx := context.WithValue(context.Background(), models.UserIDKey, "1234")
//		request = request.WithContext(ctx)
//		fmt.Println(request.Context().Value(models.UserIDKey))
//
//		resp, err := http.DefaultClient.Do(request)
//		assert.NoError(t, err)
//
//		defer resp.Body.Close()
//
//		zr, err := gzip.NewReader(resp.Body)
//		assert.NoError(t, err)
//
//		b, err := io.ReadAll(zr)
//		assert.NoError(t, err)
//		assert.NotNil(t, string(b))
//	})
//}
