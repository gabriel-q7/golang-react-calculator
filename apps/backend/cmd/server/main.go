// Command server runs the calculator HTTP server, serving the JSON API and,
// in production builds, the embedded frontend static assets.
package main

import (
	"io/fs"
	"log"
	"net/http"
	"os"

	"calculator-backend/internal/api"
	"calculator-backend/internal/service"
	"calculator-backend/internal/web"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dist, err := fs.Sub(web.DistFS, "dist")
	if err != nil {
		log.Fatalf("load embedded frontend: %v", err)
	}

	// Mounted on its own "/api/" subtree so that unmatched or wrong-method
	// API requests get the API's 404/405 handling, instead of falling
	// through to the SPA file server below (which would 404 on its own
	// terms for e.g. GET /api/add).
	mux := http.NewServeMux()
	mux.Handle("/api/", api.NewMux(service.New()))
	mux.Handle("/", http.FileServer(http.FS(dist)))

	addr := ":" + port
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
