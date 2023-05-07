package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRespondJSON_(t *testing.T) {
	w := httptest.NewRecorder()
	status := http.StatusOK
	payload := map[string]string{"message": "dummy message"}

	// test case 1: valid payload
	err := respondJSON(w, status, payload)
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
	err = respondJSON(w, http.StatusOK, func() {})
	if err == nil {
		t.Error("respondJSON should return an error for an invalid payload")
	}
}

func TestRespondError(t *testing.T) {
	w := httptest.NewRecorder()
	code := http.StatusBadRequest
	message := "invalid input"

	respondError(w, code, message)

	expectedResponse, _ := json.Marshal(map[string]string{"error": message})

	if w.Code != code {
		t.Errorf("Expected response status code %v but got %v", code, w.Code)
	}

	if w.Body.String() != string(expectedResponse) {
		t.Errorf("Expected response body %v but got %v", expectedResponse, w.Body.String())
	}
}
