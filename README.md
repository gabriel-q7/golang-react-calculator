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
│   │   ├── internal/calculator/  pure math + domain errors
│   │   ├── internal/service/     input validation + orchestration
│   │   ├── internal/api/  HTTP handlers
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

## API

`POST` JSON to any of: `/api/add`, `/api/subtract`, `/api/multiply`,
`/api/divide`, `/api/power`, `/api/sqrt`, `/api/percentage`. Every endpoint
returns `{ "result": <number> }` or `{ "error": "<message>" }`. Invalid
input (malformed JSON, non-finite numbers) and domain errors (division by
zero, negative square root, undefined exponentiation) are reported as
`400`. Full request/response shapes: [docs/api.md](docs/api.md).

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
(`tests/calculator/`, `tests/service/`, `tests/api/`). Every test file is
an external test package (`package calculator_test`, importing
`calculator-backend/internal/calculator`), so tests only exercise each
package's exported API — the same constraint an outside consumer of the
package would have.

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
- **Three backend layers** — handlers (`internal/api`) only know HTTP;
  the service (`internal/service`) validates operands and orchestrates;
  the calculator (`internal/calculator`) is pure, dependency-free math.
  Each layer is independently unit tested.
- **Frontend organized by feature, not by file type** — `validation/`,
  `api/`, `engine/`, `hooks/`, and `components/` inside
  `features/calculator/` each have one job and are independently testable,
  so the keypad UI, the calculator state machine, and the API client can
  all change without touching each other.
- **Test code lives outside the production tree** (`apps/*/tests/`, not
  colocated) — see
  [docs/adr/0002-separate-test-code-from-production-code.md](docs/adr/0002-separate-test-code-from-production-code.md).

Full detail in [docs/architecture.md](docs/architecture.md).
