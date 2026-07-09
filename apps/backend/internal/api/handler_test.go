package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"calculator-backend/internal/service"
)

func newTestMux() *http.ServeMux {
	return NewMux(service.New())
}

func doRequest(t *testing.T, mux *http.ServeMux, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()

	var reader *bytes.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal request body: %v", err)
		}
		reader = bytes.NewReader(b)
	} else {
		reader = bytes.NewReader(nil)
	}

	req := httptest.NewRequest(method, path, reader)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec
}

func decodeResult(t *testing.T, rec *httptest.ResponseRecorder) calculateResponse {
	t.Helper()
	var resp calculateResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return resp
}

func decodeError(t *testing.T, rec *httptest.ResponseRecorder) errorResponse {
	t.Helper()
	var resp errorResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error response: %v", err)
	}
	return resp
}

func TestHandleHealth(t *testing.T) {
	rec := doRequest(t, newTestMux(), http.MethodGet, "/api/health", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestHandleAdd(t *testing.T) {
	rec := doRequest(t, newTestMux(), http.MethodPost, "/api/add", binaryRequest{A: 2, B: 3})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got := decodeResult(t, rec).Result; got != 5 {
		t.Fatalf("result = %v, want 5", got)
	}
}

func TestHandleSubtract(t *testing.T) {
	rec := doRequest(t, newTestMux(), http.MethodPost, "/api/subtract", binaryRequest{A: 5, B: 3})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got := decodeResult(t, rec).Result; got != 2 {
		t.Fatalf("result = %v, want 2", got)
	}
}

func TestHandleMultiply(t *testing.T) {
	rec := doRequest(t, newTestMux(), http.MethodPost, "/api/multiply", binaryRequest{A: 4, B: 3})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got := decodeResult(t, rec).Result; got != 12 {
		t.Fatalf("result = %v, want 12", got)
	}
}

func TestHandleDivide(t *testing.T) {
	rec := doRequest(t, newTestMux(), http.MethodPost, "/api/divide", binaryRequest{A: 9, B: 3})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got := decodeResult(t, rec).Result; got != 3 {
		t.Fatalf("result = %v, want 3", got)
	}
}

func TestHandleDivide_ByZero(t *testing.T) {
	rec := doRequest(t, newTestMux(), http.MethodPost, "/api/divide", binaryRequest{A: 1, B: 0})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
	if got := decodeError(t, rec).Error; got != "division by zero" {
		t.Fatalf("error = %q, want %q", got, "division by zero")
	}
}

func TestHandlePower(t *testing.T) {
	rec := doRequest(t, newTestMux(), http.MethodPost, "/api/power", powerRequest{Base: 2, Exponent: 10})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got := decodeResult(t, rec).Result; got != 1024 {
		t.Fatalf("result = %v, want 1024", got)
	}
}

func TestHandlePower_UndefinedResult(t *testing.T) {
	rec := doRequest(t, newTestMux(), http.MethodPost, "/api/power", powerRequest{Base: -4, Exponent: 0.5})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestHandleSqrt(t *testing.T) {
	rec := doRequest(t, newTestMux(), http.MethodPost, "/api/sqrt", sqrtRequest{Value: 9})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got := decodeResult(t, rec).Result; got != 3 {
		t.Fatalf("result = %v, want 3", got)
	}
}

func TestHandleSqrt_Negative(t *testing.T) {
	rec := doRequest(t, newTestMux(), http.MethodPost, "/api/sqrt", sqrtRequest{Value: -4})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
	if got := decodeError(t, rec).Error; got != "cannot take the square root of a negative number" {
		t.Fatalf("error = %q, want %q", got, "cannot take the square root of a negative number")
	}
}

func TestHandlePercentage(t *testing.T) {
	rec := doRequest(t, newTestMux(), http.MethodPost, "/api/percentage", percentageRequest{Value: 200, Percent: 10})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got := decodeResult(t, rec).Result; got != 20 {
		t.Fatalf("result = %v, want 20", got)
	}
}

func TestHandle_InvalidJSON(t *testing.T) {
	mux := newTestMux()
	req := httptest.NewRequest(http.MethodPost, "/api/add", bytes.NewReader([]byte("{not json")))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestHandle_UnknownFields(t *testing.T) {
	mux := newTestMux()
	req := httptest.NewRequest(http.MethodPost, "/api/add", bytes.NewReader([]byte(`{"a":1,"b":2,"operator":"add"}`)))
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestHandle_MethodNotAllowed(t *testing.T) {
	rec := doRequest(t, newTestMux(), http.MethodGet, "/api/add", nil)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
	}
}

func TestHandle_UnknownRoute(t *testing.T) {
	rec := doRequest(t, newTestMux(), http.MethodPost, "/api/modulo", binaryRequest{A: 1, B: 2})
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}
