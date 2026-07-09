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
- `apps/backend/` — Go HTTP server, split into three layers (see below):
  `internal/api` (handlers), `internal/service` (validation/orchestration),
  `internal/calculator` (pure math). `internal/web` embeds the compiled
  frontend via `go:embed` so the server can serve it directly.
- `docs/` — this folder.

## Backend layering

```
internal/api          HTTP handlers: decode/encode JSON, route,
   │                   map errors → status codes
   ▼
internal/service       Validates operands are finite numbers,
   │                   delegates to calculator, no HTTP knowledge
   ▼
internal/calculator    Pure math + domain errors (division by zero,
                       negative sqrt, undefined results). No I/O.
```

Each layer is independently unit tested: `calculator` tests are pure
table-driven math tests, `service` tests cover validation and error
propagation, `api` tests exercise full HTTP request/response cycles via
`httptest`. See [api.md](api.md) for the endpoint list.

## Request flow

1. Browser loads `/` → Go server returns the embedded SPA (`index.html`,
   JS, CSS bundles).
2. SPA calls e.g. `POST /api/divide` with `{a, b}`.
3. The handler decodes the request and calls the matching
   `CalculatorService` method, which validates the operands and delegates
   to `internal/calculator`.
4. The handler returns `{result}` on success, or `{error}` with `400`
   (invalid/domain error), `404` (unknown route), or `405` (wrong method).

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
