package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/time/rate"

	"calculator-backend/internal/api"
	"calculator-backend/internal/operations"
	"calculator-backend/internal/ratelimit"
	"calculator-backend/internal/service"
)

// spyLimiter records every key it's asked about and returns a
// caller-controlled verdict, decoupling these tests from any real rate
// limiting algorithm.
type spyLimiter struct {
	allow bool
	keys  []string
}

func (s *spyLimiter) Allow(key string) bool {
	s.keys = append(s.keys, key)
	return s.allow
}

func muxWithLimiter(t *testing.T, limiter api.RateLimiter) *http.ServeMux {
	t.Helper()
	registry, err := operations.Default()
	if err != nil {
		t.Fatalf("operations.Default() error = %v", err)
	}
	return api.NewMux(service.New(registry), registry.All(), limiter)
}

func TestRateLimit_BlocksWhenLimiterDenies(t *testing.T) {
	spy := &spyLimiter{allow: false}
	mux := muxWithLimiter(t, spy)

	rec := doRequest(t, mux, http.MethodPost, "/api/add", map[string]float64{"a": 1, "b": 2})

	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusTooManyRequests)
	}
	if got := decodeError(t, rec).Error; got != "rate limit exceeded, try again later" {
		t.Fatalf("error = %q, want the rate-limit message", got)
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("Content-Type = %q, want application/json", got)
	}
	if got := rec.Header().Get("Retry-After"); got == "" {
		t.Fatal("Retry-After header missing on a 429 response")
	}
}

func TestRateLimit_AllowsWhenLimiterPermits(t *testing.T) {
	spy := &spyLimiter{allow: true}
	mux := muxWithLimiter(t, spy)

	rec := doRequest(t, mux, http.MethodPost, "/api/add", map[string]float64{"a": 1, "b": 2})

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if len(spy.keys) != 1 {
		t.Fatalf("limiter was consulted %d times, want 1", len(spy.keys))
	}
}

func TestRateLimit_HealthCheckIsExempt(t *testing.T) {
	spy := &spyLimiter{allow: false} // deny everything
	mux := muxWithLimiter(t, spy)

	rec := doRequest(t, mux, http.MethodGet, "/api/health", nil)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d (health check should bypass rate limiting)", rec.Code, http.StatusOK)
	}
	if len(spy.keys) != 0 {
		t.Fatalf("limiter was consulted for /api/health, want it untouched")
	}
}

func TestRateLimit_KeyedByRemoteAddr(t *testing.T) {
	spy := &spyLimiter{allow: true}
	mux := muxWithLimiter(t, spy)

	for _, addr := range []string{"203.0.113.1:1111", "203.0.113.2:2222"} {
		req := httptest.NewRequest(http.MethodPost, "/api/add", nil)
		req.RemoteAddr = addr
		mux.ServeHTTP(httptest.NewRecorder(), req)
	}

	if len(spy.keys) != 2 || spy.keys[0] == spy.keys[1] {
		t.Fatalf("expected two distinct keys for two distinct remote addrs, got %v", spy.keys)
	}
	if spy.keys[0] != "203.0.113.1" || spy.keys[1] != "203.0.113.2" {
		t.Fatalf("keys = %v, want the RemoteAddr host without its port", spy.keys)
	}
}

func TestRateLimit_NilLimiterDisablesRateLimiting(t *testing.T) {
	mux := muxWithLimiter(t, nil)

	for i := 0; i < 20; i++ {
		rec := doRequest(t, mux, http.MethodPost, "/api/add", map[string]float64{"a": 1, "b": 2})
		if rec.Code != http.StatusOK {
			t.Fatalf("request %d: status = %d, want %d (nil limiter must never block)", i, rec.Code, http.StatusOK)
		}
	}
}

// TestRateLimit_Integration wires the real ratelimit.IPRateLimiter (the
// one main.go uses in production) through api.NewMux, to prove the
// pieces work together end to end, not just against the spy above.
func TestRateLimit_Integration(t *testing.T) {
	limiter := ratelimit.New(rate.Limit(1), 2) // burst of 2, refilling slowly
	mux := muxWithLimiter(t, limiter)

	req := func() *httptest.ResponseRecorder {
		r := httptest.NewRequest(http.MethodPost, "/api/add", nil)
		r.RemoteAddr = "198.51.100.7:5555"
		rec := httptest.NewRecorder()
		mux.ServeHTTP(rec, r)
		return rec
	}

	if got := req().Code; got != http.StatusBadRequest {
		// no body on this bare request, so it's a 400 (invalid request
		// body) rather than 200 — the point here is only that it got
		// PAST the limiter, not that the calculation itself succeeded.
		t.Fatalf("request 1: status = %d, want %d", got, http.StatusBadRequest)
	}
	if got := req().Code; got != http.StatusBadRequest {
		t.Fatalf("request 2: status = %d, want %d", got, http.StatusBadRequest)
	}
	if got := req().Code; got != http.StatusTooManyRequests {
		t.Fatalf("request 3: status = %d, want %d (burst of 2 exhausted)", got, http.StatusTooManyRequests)
	}
}
