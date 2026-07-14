package api

import (
	"net"
	"net/http"
)

// RateLimiter decides whether a request identified by key may proceed.
// It's defined here, not imported from a concrete limiter package, so
// this HTTP layer only depends on this one-method shape. Swapping the
// underlying algorithm (internal/ratelimit's per-IP token bucket today)
// for something else — or swapping what "key" means, see clientIPKey —
// never requires touching this file.
type RateLimiter interface {
	Allow(key string) bool
}

// rateLimitMiddleware rejects requests with 429 once keyFunc's key has
// exceeded limiter's rate, using the same JSON error shape as every
// other error response.
func rateLimitMiddleware(limiter RateLimiter, keyFunc func(*http.Request) string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow(keyFunc(r)) {
			w.Header().Set("Retry-After", "60")
			writeError(w, http.StatusTooManyRequests, "rate limit exceeded, try again later")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// clientIPKey extracts the connecting peer's IP from r.RemoteAddr, which
// is what net/http sets from the raw TCP connection.
//
// This assumes the server is reached directly (true today: one
// container, no reverse proxy in front of it — see docs/architecture.md).
// If that changes, replace this function with one that reads
// X-Forwarded-For/X-Real-IP from a *trusted* proxy — trusting those
// headers from an untrusted client lets them spoof their rate-limit key.
// Nothing else needs to change: NewMux takes keyFunc as a parameter.
func clientIPKey(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
