package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"calculator-backend/internal/api"
	"calculator-backend/internal/service"
)

// resultBody and errorBody mirror the JSON contract documented in
// docs/api.md. They're defined here, not imported from internal/api,
// because these tests exercise the HTTP/JSON interface only — the
// package's own (unexported) request/response types are an implementation
// detail the tests have no business depending on.
type resultBody struct {
	Result float64 `json:"result"`
}

type errorBody struct {
	Error string `json:"error"`
}

func newTestMux() *http.ServeMux {
	return api.NewMux(service.New())
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

func decodeResult(t *testing.T, rec *httptest.ResponseRecorder) resultBody {
	t.Helper()
	var resp resultBody
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	return resp
}

func decodeError(t *testing.T, rec *httptest.ResponseRecorder) errorBody {
	t.Helper()
	var resp errorBody
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
	rec := doRequest(t, newTestMux(), http.MethodPost, "/api/add", map[string]float64{"a": 2, "b": 3})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got := decodeResult(t, rec).Result; got != 5 {
		t.Fatalf("result = %v, want 5", got)
	}
}

func TestHandleSubtract(t *testing.T) {
	rec := doRequest(t, newTestMux(), http.MethodPost, "/api/subtract", map[string]float64{"a": 5, "b": 3})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got := decodeResult(t, rec).Result; got != 2 {
		t.Fatalf("result = %v, want 2", got)
	}
}

func TestHandleMultiply(t *testing.T) {
	rec := doRequest(t, newTestMux(), http.MethodPost, "/api/multiply", map[string]float64{"a": 4, "b": 3})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got := decodeResult(t, rec).Result; got != 12 {
		t.Fatalf("result = %v, want 12", got)
	}
}

func TestHandleDivide(t *testing.T) {
	rec := doRequest(t, newTestMux(), http.MethodPost, "/api/divide", map[string]float64{"a": 9, "b": 3})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got := decodeResult(t, rec).Result; got != 3 {
		t.Fatalf("result = %v, want 3", got)
	}
}

func TestHandleDivide_ByZero(t *testing.T) {
	rec := doRequest(t, newTestMux(), http.MethodPost, "/api/divide", map[string]float64{"a": 1, "b": 0})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
	if got := decodeError(t, rec).Error; got != "division by zero" {
		t.Fatalf("error = %q, want %q", got, "division by zero")
	}
}

func TestHandlePower(t *testing.T) {
	rec := doRequest(t, newTestMux(), http.MethodPost, "/api/power", map[string]float64{"base": 2, "exponent": 10})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got := decodeResult(t, rec).Result; got != 1024 {
		t.Fatalf("result = %v, want 1024", got)
	}
}

func TestHandlePower_UndefinedResult(t *testing.T) {
	rec := doRequest(t, newTestMux(), http.MethodPost, "/api/power", map[string]float64{"base": -4, "exponent": 0.5})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestHandleSqrt(t *testing.T) {
	rec := doRequest(t, newTestMux(), http.MethodPost, "/api/sqrt", map[string]float64{"value": 9})
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if got := decodeResult(t, rec).Result; got != 3 {
		t.Fatalf("result = %v, want 3", got)
	}
}

func TestHandleSqrt_Negative(t *testing.T) {
	rec := doRequest(t, newTestMux(), http.MethodPost, "/api/sqrt", map[string]float64{"value": -4})
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
	if got := decodeError(t, rec).Error; got != "cannot take the square root of a negative number" {
		t.Fatalf("error = %q, want %q", got, "cannot take the square root of a negative number")
	}
}

func TestHandlePercentage(t *testing.T) {
	rec := doRequest(t, newTestMux(), http.MethodPost, "/api/percentage", map[string]float64{"value": 200, "percent": 10})
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
	rec := doRequest(t, newTestMux(), http.MethodPost, "/api/modulo", map[string]float64{"a": 1, "b": 2})
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusNotFound)
	}
}
