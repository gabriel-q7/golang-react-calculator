// Command server runs the calculator HTTP server, serving the JSON API and,
// in production builds, the embedded frontend static assets.
package main

import (
	"io/fs"
	"log"
	"net/http"
	"os"
	"strconv"

	"golang.org/x/time/rate"

	"calculator-backend/internal/api"
	"calculator-backend/internal/operations"
	"calculator-backend/internal/ratelimit"
	"calculator-backend/internal/service"
	"calculator-backend/internal/web"
)

// Rate-limit defaults, used whenever the corresponding env var is unset
// or invalid. 60 req/min with a burst of 10 comfortably covers a single
// interactive user clicking through calculations, while still bounding
// how much load one IP can put on the server.
const (
	defaultRateLimitRPM   = 60
	defaultRateLimitBurst = 10
)

func main() {
	port := getEnv("PORT", "8080")

	registry, err := operations.Default()
	if err != nil {
		log.Fatalf("build operation registry: %v", err)
	}
	svc := service.New(registry)

	rpm := getEnvInt("RATE_LIMIT_RPM", defaultRateLimitRPM)
	burst := getEnvInt("RATE_LIMIT_BURST", defaultRateLimitBurst)
	limiter := ratelimit.New(rate.Limit(float64(rpm)/60), burst)

	dist, err := fs.Sub(web.DistFS, "dist")
	if err != nil {
		log.Fatalf("load embedded frontend: %v", err)
	}

	// Mounted on its own "/api/" subtree so that unmatched or wrong-method
	// API requests get the API's 404/405 handling, instead of falling
	// through to the SPA file server below (which would 404 on its own
	// terms for e.g. GET /api/add).
	mux := http.NewServeMux()
	mux.Handle("/api/", api.NewMux(svc, registry.All(), limiter))
	mux.Handle("/", http.FileServer(http.FS(dist)))

	addr := ":" + port
	log.Printf("listening on %s (rate limit: %d req/min, burst %d)", addr, rpm, burst)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		log.Printf("invalid %s=%q, using default %d", key, v, fallback)
		return fallback
	}
	return n
}
