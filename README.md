# golang-react-calculator

A small calculator app: a Go HTTP API backing a React + TypeScript
frontend, shipped in production as a single Docker image (the Go binary
serves the compiled frontend via `go:embed`).

## Project layout

```
.
├── apps/
│   ├── backend/          Go API (net/http, no framework)
│   │   ├── cmd/server/    entrypoint
│   │   ├── internal/operations/  Operation interface + 7 implementations + registry
│   │   ├── internal/service/     resolves + validates + executes an operation
│   │   ├── internal/api/  HTTP handlers + rate-limit middleware
│   │   ├── internal/ratelimit/  per-IP token-bucket limiter
│   │   ├── internal/web/  embeds the built frontend (go:embed)
│   │   └── tests/         all *_test.go files (see "Testing" below)
│   └── frontend/         React + TypeScript SPA (Vite, Tailwind, shadcn/ui)
│       ├── src/
│       │   ├── components/ui/       shadcn/ui primitives (Button, Card, Alert, ...)
│       │   ├── features/calculator/  feature module (see below)
│       │   └── lib/                  shared helpers (`cn` class merger)
│       └── tests/        all *.test.ts(x) files (see "Testing" below)
├── docs/
│   ├── architecture.md  system design and rationale
│   ├── api.md           endpoint reference
│   └── adr/              Architecture Decision Records
├── docker-compose.yml    local dev stack (hot-reloading, two containers)
├── Dockerfile            production image (multi-stage, single container)
└── Makefile
```

See [docs/architecture.md](docs/architecture.md) for the full design
rationale, [docs/api.md](docs/api.md) for the API contract, and
[docs/adr/](docs/adr/) for the reasoning behind specific decisions (the
hex-keypad/vaporwave UI, separating tests from production code, etc).

## Backend

`POST` JSON to any of: `/api/add`, `/api/subtract`, `/api/multiply`,
`/api/divide`, `/api/power`, `/api/sqrt`, `/api/percentage`. Every endpoint
returns `{ "result": <number> }` or `{ "error": "<message>" }`. Invalid
input (malformed JSON, non-finite numbers), domain errors (division by
zero, negative square root, undefined exponentiation), and rate limiting
are all reported with that same `{ "error": ... }` shape (`400`/`429`).
Full request/response shapes: [docs/api.md](docs/api.md).

Each operation is an independent `operations.Operation` implementation
(`internal/operations/`) — a `Registry` resolves them by name, so the
service and HTTP layers never enumerate operations themselves. Adding an
8th operation is: implement `Operation`, register it in
`operations.Default()`, done. See
[docs/adr/0003-solid-refactor-and-rate-limiting.md](docs/adr/0003-solid-refactor-and-rate-limiting.md)
for the full rationale.

Every operation route is also rate-limited per client IP (not the health
check), configurable via `RATE_LIMIT_RPM` (default `60`) and
`RATE_LIMIT_BURST` (default `10`) environment variables:

```bash
docker run -p 8080:8080 -e RATE_LIMIT_RPM=120 -e RATE_LIMIT_BURST=20 calculator:latest
```

## Frontend

React + TypeScript, styled with Tailwind CSS v4 and shadcn/ui, themed as a
dark "vaporwave" physical calculator: a display (running expression +
current value) above a honeycomb keypad of hexagonal keys. See
[docs/adr/0001-hex-keypad-vaporwave-theme.md](docs/adr/0001-hex-keypad-vaporwave-theme.md)
for why.

`apps/frontend/src/features/calculator/` is organized by concern:

| Path              | Responsibility                                                                   |
| ----------------- | --------------------------------------------------------------------------------- |
| `config.ts`        | Declarative list of the 7 backend operations (id, label, endpoint, fields)        |
| `validation/`      | Pure functions: parse a field to a number, run client-side domain checks (e.g. reject division by zero) before any request is sent |
| `api/`             | `postCalculate` — the only place that calls `fetch`                              |
| `engine/`          | `reducer.ts` (calculator state machine), `operations.ts` (operator → endpoint), `format.ts` (display number formatting) |
| `hooks/`           | `useCalculatorEngine` — wires the reducer + API calls together                   |
| `components/`      | `HexButton`/`HexKeypad` (the keypad), `CalculatorDisplay`, `CalculatorPage`       |

Client-side validation mirrors the backend's rules (division by zero,
negative square root, etc.) so obviously-invalid requests never round-trip
to the API — but the backend remains the source of truth and re-validates
independently. `components/ui/` holds the shadcn/ui primitives (Button,
Input, Card, Alert), hand-authored against this project's Tailwind theme
rather than pulled in via the shadcn CLI, since only a handful of
components are needed.

## Requirements

- Docker + Docker Compose (for `make up` / `make build`)
- Go 1.23+ and Node 22+ (only if you want to run things outside Docker)

## Quick start (local development)

```bash
make up
```

This builds and starts two containers via `docker-compose.yml`:

- **backend** — `go run ./cmd/server` on `localhost:8080`, source
  bind-mounted so restarting the container picks up changes.
- **frontend** — Vite dev server with HMR on `localhost:5173`, proxying
  `/api/*` to the backend container.

Open `http://localhost:5173`. Stop everything with:

```bash
make down
```

## Production build

```bash
make build          # docker build -t calculator:latest .
docker run -p 8080:8080 calculator:latest
```

Open `http://localhost:8080` — one container serves both the API and the
compiled frontend.

## Makefile targets

| Target           | Description                                          |
| ---------------- | ----------------------------------------------------- |
| `make build`      | Build the production Docker image                     |
| `make up`         | Start the local dev stack (frontend + backend)         |
| `make down`       | Stop the local dev stack                               |
| `make test`       | Run backend (`go test`) and frontend (`vitest`) suites |
| `make coverage`   | Generate coverage reports for backend and frontend     |
| `make clean`      | Remove containers, volumes, `node_modules`, build output |

## Testing

Test code is kept out of the production packages/source tree in both
apps — see
[docs/adr/0002-separate-test-code-from-production-code.md](docs/adr/0002-separate-test-code-from-production-code.md)
for the rationale.

**Backend** — `apps/backend/tests/` mirrors `internal/`'s package layout
(`tests/operations/`, `tests/service/`, `tests/api/`, `tests/ratelimit/`).
Every test file is an external test package (`package operations_test`,
importing `calculator-backend/internal/operations`), so tests only
exercise each package's exported API — the same constraint an outside
consumer of the package would have.

**Frontend** — `apps/frontend/tests/` mirrors `src/features/calculator/`.
Tests import production code via the `@/` alias (e.g.
`@/features/calculator/engine/format`) rather than relative paths, so a
test's location never has to match the file it's testing, and production
code never imports from `tests/` (no risk of a circular dependency).
`tests/setup.ts` (jest-dom matchers) is the one non-test file in that
tree.

Run them:

```bash
# from the repo root
make test              # backend + frontend
make coverage          # backend + frontend, with coverage reports

# directly
cd apps/backend && go test ./...                                   # tests/... has the actual tests
cd apps/backend && go test ./tests/... -coverpkg=./internal/... -coverprofile=coverage.out
cd apps/frontend && npm test                                       # runs tests/**/*.test.{ts,tsx}
cd apps/frontend && npm run coverage                                # coverage measured over src/, not tests/
```

## Design rationale (short version)

- **One Go process serves everything in production** — the frontend is
  embedded into the binary at build time, so there's no nginx, no static
  host, no CORS to configure, and deployment is "run one container."
- **Two containers in development** — Vite's dev server gives fast HMR for
  frontend work; the Go server runs separately so backend changes don't
  require a frontend rebuild. `docker-compose.yml` wires them together with
  a proxy so routing matches production.
- **No web framework on the backend** — `net/http.ServeMux`'s method-aware
  routing (Go 1.22+) is enough for a handful of routes plus a health check.
- **Backend built around a small `Operation` abstraction** — handlers
  (`internal/api`) only know HTTP; the service (`internal/service`)
  validates operands and resolves an operation by name through a
  1-method `OperationResolver` interface; each operation
  (`internal/operations`) is its own type knowing only its own math. New
  operations are additive (Open/Closed), and each layer only depends on
  interfaces, not the layer below's concrete types (Dependency
  Inversion) — see
  [docs/adr/0003-solid-refactor-and-rate-limiting.md](docs/adr/0003-solid-refactor-and-rate-limiting.md).
- **IP-based rate limiting, no auth** — the frontend has no login, so a
  per-IP token bucket (`internal/ratelimit`, `golang.org/x/time/rate`) is
  the cheapest meaningful abuse guard available; see the same ADR for why,
  and how it's designed to swap for auth-based limiting later.
- **Frontend organized by feature, not by file type** — `validation/`,
  `api/`, `engine/`, `hooks/`, and `components/` inside
  `features/calculator/` each have one job and are independently testable,
  so the keypad UI, the calculator state machine, and the API client can
  all change without touching each other.
- **Test code lives outside the production tree** (`apps/*/tests/`, not
  colocated) — see
  [docs/adr/0002-separate-test-code-from-production-code.md](docs/adr/0002-separate-test-code-from-production-code.md).

Full detail in [docs/architecture.md](docs/architecture.md).
