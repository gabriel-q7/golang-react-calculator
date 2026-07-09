package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleCalculate(t *testing.T) {
	mux := NewMux()

	body, _ := json.Marshal(calculateRequest{A: 2, B: 3, Operator: "add"})
	req := httptest.NewRequest(http.MethodPost, "/api/calculate", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var resp calculateResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Result != 5 {
		t.Fatalf("result = %v, want 5", resp.Result)
	}
}

func TestHandleCalculate_DivisionByZero(t *testing.T) {
	mux := NewMux()

	body, _ := json.Marshal(calculateRequest{A: 1, B: 0, Operator: "divide"})
	req := httptest.NewRequest(http.MethodPost, "/api/calculate", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestHandleHealth(t *testing.T) {
	mux := NewMux()

	req := httptest.NewRequest(http.MethodGet, "/api/health", nil)
	rec := httptest.NewRecorder()

	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}
