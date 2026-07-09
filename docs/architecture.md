# Architecture

## Overview

The project is a small calculator application split into two independently
buildable pieces that ship as **one** deployable artifact:

```
┌─────────────────────────────────────────────┐
│                 Docker image                 │
│                                               │
│   React/TS SPA (static files)                │
│        │  served by                          │
│        ▼                                     │
│   Go HTTP server  ── /api/*  ──▶  calculator  │
│   (net/http, no framework)         package    │
└─────────────────────────────────────────────┘
```

- `apps/frontend/` — React + TypeScript SPA built with Vite. Talks to the
  backend over `/api/*`.
- `apps/backend/` — Go HTTP server. `internal/calculator` holds the pure
  arithmetic logic, `internal/api` exposes it over JSON, `internal/web`
  embeds the compiled frontend via `go:embed` so the server can serve it
  directly.
- `docs/` — this folder.

## Request flow

1. Browser loads `/` → Go server returns the embedded SPA (`index.html`,
   JS, CSS bundles).
2. SPA calls `POST /api/calculate` with `{a, b, operator}`.
3. `internal/api` decodes the request, calls `internal/calculator.Calculate`,
   and returns `{result}` or `{error}` with an appropriate HTTP status.

## Why a single container in production

Shipping one Go binary with the frontend embedded (`go:embed`) means:

- **One artifact, one process** — no reverse proxy, no nginx, no static
  file host to operate or keep in sync with the API.
- **No CORS** — same-origin frontend and API, so `/api/*` needs no CORS
  headers or credentials handling.
- **Trivial deploy** — a single container image runs anywhere Docker runs.

The tradeoff: frontend and backend release together. That's an acceptable
tradeoff for a project this size; a larger product would likely split them
to scale/deploy independently.

## Why separate directories in development

Locally, `apps/frontend` and `apps/backend` run as separate processes
(`docker-compose.yml`): Vite's dev server for fast HMR, and `go run` for the
API. Vite proxies `/api/*` to the Go server (`apps/frontend/vite.config.ts`),
so the dev experience matches production routing without rebuilding the Go
binary on every frontend change.

## Why `net/http` with no framework

The API surface is one endpoint plus a health check. A router/framework
(gin, echo, chi) would add a dependency and abstraction for no real benefit
at this scale. `net/http.ServeMux` is enough; revisit if the API grows
past a handful of routes or needs middleware chains.

## Testing

- Go: standard `testing` package, table-driven tests for the calculator
  logic, `httptest` for the HTTP handlers.
- Frontend: Vitest + Testing Library, mocking `fetch` to test API
  integration without a running backend.

See [api.md](api.md) for the endpoint contract.
