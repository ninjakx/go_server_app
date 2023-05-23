package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type testResponseWriter struct {
	*httptest.ResponseRecorder
}

func (w testResponseWriter) Header() http.Header {
	return w.ResponseRecorder.Header()
}

func (w testResponseWriter) Write(b []byte) (int, error) {
	return w.ResponseRecorder.Write(b)
}

func (w testResponseWriter) WriteHeader(statusCode int) {
	w.ResponseRecorder.WriteHeader(statusCode)
}

func TestRespondJSON_(t *testing.T) {
	w := &testResponseWriter{httptest.NewRecorder()}

	// Create a new Gin context with the custom response writer
	c, _ := gin.CreateTestContext(w)
	status := http.StatusOK
	payload := map[string]string{"message": "dummy message"}

	// test case 1: valid payload
	err := respondJSON(c, status, payload)
	if err != nil {
		t.Fatalf("Error should be nil but got %v", err)
	}

	expectedResponse, _ := json.Marshal(payload)
	if w.Code != status {
		t.Errorf("Expected response status code %v but got %v", status, w.Code)
	}

	if w.Body.String() != string(expectedResponse) {
		t.Errorf("Expected response body %v but got %v", expectedResponse, w.Body.String())
	}

	// test case 2: invalid payload
	err = respondJSON(c, http.StatusOK, func() {})
	if err == nil {
		t.Error("respondJSON should return an error for an invalid payload")
	}
}

func TestRespondError(t *testing.T) {
	w := &testResponseWriter{httptest.NewRecorder()}
	// Create a new Gin context with the custom response writer
	c, _ := gin.CreateTestContext(w)
	code := http.StatusBadRequest
	message := "invalid input"

	respondError(c, code, message)

	expectedResponse, _ := json.Marshal(map[string]string{"error": message})

	if w.Code != code {
		t.Errorf("Expected response status code %v but got %v", code, w.Code)
	}

	if w.Body.String() != string(expectedResponse) {
		t.Errorf("Expected response body %v but got %v", expectedResponse, w.Body.String())
	}
}
