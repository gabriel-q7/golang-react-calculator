package ratelimit_test

import (
	"testing"
	"time"

	"golang.org/x/time/rate"

	"calculator-backend/internal/ratelimit"
)

func TestIPRateLimiter_AllowsUpToBurstThenBlocks(t *testing.T) {
	limiter := ratelimit.New(rate.Limit(1), 3) // 1 req/sec, burst of 3

	for i := 0; i < 3; i++ {
		if !limiter.Allow("1.2.3.4") {
			t.Fatalf("request %d within burst was blocked", i+1)
		}
	}
	if limiter.Allow("1.2.3.4") {
		t.Fatal("request beyond burst was allowed")
	}
}

func TestIPRateLimiter_KeysAreIndependent(t *testing.T) {
	limiter := ratelimit.New(rate.Limit(1), 1)

	if !limiter.Allow("1.2.3.4") {
		t.Fatal("first request for 1.2.3.4 was blocked")
	}
	if limiter.Allow("1.2.3.4") {
		t.Fatal("second immediate request for 1.2.3.4 was allowed (burst is 1)")
	}
	if !limiter.Allow("5.6.7.8") {
		t.Fatal("first request for a different key was blocked by 1.2.3.4's bucket")
	}
}

func TestIPRateLimiter_RefillsOverTime(t *testing.T) {
	// Fast enough to refill within a short, CI-safe sleep, without
	// reaching into the limiter's internals to fake the clock.
	limiter := ratelimit.New(rate.Limit(50), 1)

	if !limiter.Allow("1.2.3.4") {
		t.Fatal("first request was blocked")
	}
	if limiter.Allow("1.2.3.4") {
		t.Fatal("second immediate request was allowed (burst is 1)")
	}

	time.Sleep(100 * time.Millisecond) // 50 req/sec refills well within this
	if !limiter.Allow("1.2.3.4") {
		t.Fatal("request after the refill interval was still blocked")
	}
}
