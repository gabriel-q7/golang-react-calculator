// Package web embeds the compiled frontend assets into the backend binary.
package web

import "embed"

// DistFS holds the production frontend build (frontend/dist), copied into
// this directory by the root Dockerfile before `go build`. In local dev the
// placeholder in dist/ is served instead; run `make build` for the real app.
//
//go:embed all:dist
var DistFS embed.FS
