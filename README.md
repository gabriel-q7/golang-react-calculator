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
│   │   └── internal/web/  embeds the built frontend (go:embed)
│   └── frontend/         React + TypeScript SPA (Vite, Tailwind, shadcn/ui)
│       └── src/
│           ├── components/ui/       shadcn/ui primitives (Button, Card, Tabs, ...)
│           ├── features/calculator/  feature module (see below)
│           └── lib/                  shared helpers (`cn` class merger)
├── docs/
│   ├── architecture.md  system design and rationale
│   └── api.md           endpoint reference
├── docker-compose.yml    local dev stack (hot-reloading, two containers)
├── Dockerfile            production image (multi-stage, single container)
└── Makefile
```

See [docs/architecture.md](docs/architecture.md) for the full design
rationale and [docs/api.md](docs/api.md) for the API contract.

## API

`POST` JSON to any of: `/api/add`, `/api/subtract`, `/api/multiply`,
`/api/divide`, `/api/power`, `/api/sqrt`, `/api/percentage`. Every endpoint
returns `{ "result": <number> }` or `{ "error": "<message>" }`. Invalid
input (malformed JSON, non-finite numbers) and domain errors (division by
zero, negative square root, undefined exponentiation) are reported as
`400`. Full request/response shapes: [docs/api.md](docs/api.md).

## Frontend

React + TypeScript, styled with Tailwind CSS v4 and shadcn/ui components.
One tabbed card per operation (Add, Subtract, Multiply, Divide, Power,
Square Root, Percentage); each tab is a self-contained form with its own
inputs, validation state, and result/error display. Layout is responsive
— the tab bar scrolls horizontally and labels collapse to symbols on
narrow screens.

`apps/frontend/src/features/calculator/` is organized by concern:

| Path             | Responsibility                                                     |
| ----------------- | --------------------------------------------------------------------- |
| `config.ts`        | Declarative list of operations (id, label, endpoint, fields)          |
| `validation/`      | Pure functions: parse a field to a number, run client-side domain checks (e.g. reject division by zero) before any request is sent |
| `api/`             | `postCalculate` — the only place that calls `fetch`                   |
| `hooks/`           | `useOperationForm` — wires validation + API calls to form state       |
| `components/`      | `OperationCard` (one operation's form) and `CalculatorPage` (tab shell) |

Client-side validation mirrors the backend's rules (division by zero,
negative square root, etc.) so obviously-invalid requests never round-trip
to the API — but the backend remains the source of truth and re-validates
independently. `components/ui/` holds the shadcn/ui primitives (Button,
Input, Card, Tabs, Alert), hand-authored against this project's Tailwind
theme rather than pulled in via the shadcn CLI, since only a handful of
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

## Running tests without Docker

```bash
cd apps/backend && go test ./...
cd apps/frontend && npm install && npm test
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
  `api/`, `hooks/`, and `components/` inside `features/calculator/` are
  each independently testable and reusable across the 7 operation forms.

Full detail in [docs/architecture.md](docs/architecture.md).
