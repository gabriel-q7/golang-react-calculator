// Package ratelimit implements a per-key (in practice, per-IP) request
// rate limiter, independent of HTTP — see internal/api's middleware for
// how it's wired into a request pipeline.
package ratelimit

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// staleAfter is how long a key's limiter is kept after its last request
// before IPRateLimiter evicts it, so a long-running server doesn't
// accumulate one entry per client IP forever.
const staleAfter = 3 * time.Minute

type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// IPRateLimiter is a token-bucket rate limiter (golang.org/x/time/rate)
// per key, guarded by a mutex so it's safe for concurrent use across
// request goroutines.
type IPRateLimiter struct {
	mu        sync.Mutex
	visitors  map[string]*visitor
	rate      rate.Limit
	burst     int
	lastSweep time.Time
	now       func() time.Time
}

// New creates an IPRateLimiter allowing r requests per second, per key,
// with burst as the maximum number of requests admitted in a single
// instant (before the steady-state rate applies). Use rate.Limit(rpm/60)
// to configure it in requests-per-minute terms.
func New(r rate.Limit, burst int) *IPRateLimiter {
	return &IPRateLimiter{
		visitors: make(map[string]*visitor),
		rate:     r,
		burst:    burst,
		now:      time.Now,
	}
}

// Allow reports whether a request for key is within the limit. The
// first call for a given key creates its bucket; unseen keys always
// start with a full burst allowance.
func (l *IPRateLimiter) Allow(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := l.now()
	l.sweep(now)

	v, ok := l.visitors[key]
	if !ok {
		v = &visitor{limiter: rate.NewLimiter(l.rate, l.burst)}
		l.visitors[key] = v
	}
	v.lastSeen = now

	return v.limiter.AllowN(now, 1)
}

// sweep evicts keys idle for longer than staleAfter. It runs at most
// once per staleAfter interval (checked under the same lock Allow already
// holds), so it adds no per-request cost beyond a timestamp comparison.
func (l *IPRateLimiter) sweep(now time.Time) {
	if now.Sub(l.lastSweep) < staleAfter {
		return
	}
	for key, v := range l.visitors {
		if now.Sub(v.lastSeen) > staleAfter {
			delete(l.visitors, key)
		}
	}
	l.lastSweep = now
}
