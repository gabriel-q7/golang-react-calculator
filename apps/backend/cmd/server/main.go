// Command server runs the calculator HTTP server, serving the JSON API and,
// in production builds, the embedded frontend static assets.
package main

import (
	"io/fs"
	"log"
	"net/http"
	"os"

	"calculator-backend/internal/api"
	"calculator-backend/internal/web"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := api.NewMux()

	dist, err := fs.Sub(web.DistFS, "dist")
	if err != nil {
		log.Fatalf("load embedded frontend: %v", err)
	}
	mux.Handle("/", http.FileServer(http.FS(dist)))

	addr := ":" + port
	log.Printf("listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
