package handlers

//import (
//	"net/http"
//	"practicumserver/internal/logger"
//	"practicumserver/internal/storage"
//	"testing"
//)
//
//func TestPostRequestAPIShortenBatch(t *testing.T) {
//	log, _ := logger.NewZapLogger(false)
//	strg, _ := storage.NewStorage(log.Logger, "", "")
//	hndlrs := NewHandlers(strg, log.Logger, "http://localhost:8080")
//
//	type wantPost struct {
//		expectedCode int
//		expectedBody string
//		expectCt     string
//	}
//	type request struct {
//		method      string
//		body        string
//		contentType string
//		url         string
//	}
//	tests := []struct {
//		name string
//		request
//		wantPost
//	}{
//		{
//			name: "post request without body, result: status code 500",
//			request: request{
//				method: http.MethodPost,
//				body: "",
//				contentType: "application/json",
//				url: "http://localhost:8080/api/shorten/batch",
//			},
//			wantPost: wantPost{
//				expectedCode: http.StatusInternalServerError,
//				expectedBody: "",
//				expectCt: "",
//			},
//		},
//	}
//}
