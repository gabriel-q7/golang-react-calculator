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
- `apps/backend/` — Go HTTP server, split into layers (see below):
  `internal/api` (handlers + rate-limit middleware), `internal/service`
  (validation/orchestration), `internal/operations` (the calculator
  operations themselves), `internal/ratelimit` (the rate limiter).
  `internal/web` embeds the compiled frontend via `go:embed` so the
  server can serve it directly.
- `docs/` — this folder.

## Backend layering

```
internal/api          HTTP handlers + rate-limit middleware: decode/
   │                   encode JSON, route, map errors → status codes
   ▼
internal/service       Validates operands are finite numbers, resolves
   │                   an operation by name, delegates to it — has no
   │                   idea which operations exist (see below)
   ▼
internal/operations    Operation interface + one implementation per
                       operation (add/subtract/multiply/divide/power/
                       sqrt/percentage) + a Registry that resolves them
                       by name. No I/O, no HTTP.
```

`internal/api` and `internal/service` depend only on small interfaces
(`Executor`, `RateLimiter`, `OperationResolver`) — never on each other's
or `internal/operations`' concrete types — so each layer is independently
unit-testable with fakes, and adding an 8th operation touches only
`internal/operations` (implement `Operation`, register it in
`operations.Default()`); the service and HTTP layers need no changes.
See [docs/adr/0003-solid-refactor-and-rate-limiting.md](adr/0003-solid-refactor-and-rate-limiting.md)
for the full rationale, and [api.md](api.md) for the endpoint list.

## Request flow

1. Browser loads `/` → Go server returns the embedded SPA (`index.html`,
   JS, CSS bundles).
2. SPA calls e.g. `POST /api/divide` with `{a, b}`.
3. The rate-limit middleware checks the client IP against its limit
   (`internal/ratelimit`); if exceeded, the request stops here with `429`.
4. The handler decodes the named operands `/api/divide` expects, and
   calls `CalculatorService.Execute("divide", values)`, which validates
   the operands are finite and delegates to the `divide` operation.
5. The handler returns `{result}` on success, or `{error}` with `400`
   (invalid/domain error), `404` (unknown route), `405` (wrong method), or
   `429` (rate limited).

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

## Rate limiting

Every operation route (not the health check) is rate-limited per client
IP — see [api.md](api.md#rate-limiting) for the response shape and
[docs/adr/0003-solid-refactor-and-rate-limiting.md](adr/0003-solid-refactor-and-rate-limiting.md)
for why IP-based limiting (rather than authentication) was chosen, and
how it's designed to evolve if auth is added later.

## Testing

- Go: standard `testing` package, table-driven tests for each operation
  and the registry, `httptest` for the HTTP handlers and rate-limit
  middleware. Test files live in `apps/backend/tests/`, outside the
  production tree — see
  [docs/adr/0002-separate-test-code-from-production-code.md](adr/0002-separate-test-code-from-production-code.md).
- Frontend: Vitest + Testing Library, mocking `fetch` to test API
  integration without a running backend.

See [api.md](api.md) for the endpoint contract.
