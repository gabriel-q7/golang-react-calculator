# 0003. SOLID backend refactor and IP-based rate limiting

## Status

Accepted

## Context

The backend worked and was already reasonably layered (`internal/api` →
`internal/service` → `internal/calculator`), but adding an 8th operation
would have meant touching four different places every time:

- A new request DTO in `internal/api/types.go`.
- A new `handle<Op>` method plus a new route registration in
  `internal/api/handler.go`.
- A new method on `CalculatorService` in `internal/service/service.go`
  that re-implemented the same finite-input validation as the other six.
- A new pure function (plus, sometimes, a new sentinel error) in
  `internal/calculator/calculator.go`.

None of that is broken, exactly, but it's the textbook Open/Closed
violation: "add a feature" meant "edit N existing files," and the
`CalculatorService` had grown one method per operation
(`Add`/`Subtract`/`Multiply`/`Divide`/`Power`/`Sqrt`/`Percentage`) that
were identical except for which pure function they called — a `switch`-
shaped design in everything but literal syntax. Separately, the frontend
has no login and no API keys, so nothing stood between the public
internet and the calculation endpoints except the endpoints themselves —
worth a minimum abuse guard before this goes further.

## Decision

### An `Operation` abstraction, one type per operation, a registry

```go
// internal/operations/operation.go
type Operation interface {
    Name() string                                    // "add" — the route name
    Operands() []string                              // e.g. []string{"a", "b"}
    Apply(values map[string]float64) (float64, error) // the actual math
}
```

Each operation (`add.go`, `subtract.go`, ..., `percentage.go`) is an
unexported struct implementing this 3-method interface and a `New*()`
constructor — e.g. `divideOperation` owns `ErrDivisionByZero` and the
`b == 0` check; `sqrtOperation` owns `ErrNegativeSqrt`. A `Registry`
(`internal/operations/registry.go`) resolves operations by name and lists
them all; `operations.Default()` is the **one** function that has to
enumerate all seven — everything downstream works off the interface.

This directly targets every SOLID letter the task asked for:

- **SRP** — each operation file has exactly one reason to change (its own
  math); the registry's only job is name → operation lookup; the service
  only validates + delegates; the HTTP layer only does HTTP.
- **OCP** — adding operation #8 is: write `newthing.go` implementing
  `Operation`, add one line to `Default()`. `internal/service` and
  `internal/api` do not change.
- **LSP** — every `Operation` is used only through the interface (`Name`,
  `Operands`, `Apply`); any implementation is a drop-in substitute for
  any other, by construction — there's no operation-specific branching
  anywhere that would break if one were swapped for another.
- **ISP** — `Operation` is 3 small methods, not a fat interface with
  every operation's concerns mixed in. The two interfaces built to
  connect the layers are even smaller: `service.OperationResolver` is
  one method (`Get`), `api.Executor` is one method (`Execute`),
  `api.RateLimiter` is one method (`Allow`).
- **DIP** — `CalculatorService` depends on `OperationResolver`
  (satisfied structurally by `*operations.Registry`, but the service
  never imports or names that type); `internal/api` depends on its own
  `Executor` and `RateLimiter` interfaces and doesn't import
  `internal/service` or `internal/ratelimit` **at all** — main.go is the
  only place concrete types from every package meet.

### The service collapses to one method

```go
func (s *CalculatorService) Execute(name string, values map[string]float64) (float64, error) {
    op, ok := s.resolver.Get(name)
    if !ok {
        return 0, ErrUnknownOperation
    }
    for _, v := range values {
        if math.IsNaN(v) || math.IsInf(v, 0) {
            return 0, ErrInvalidInput
        }
    }
    return op.Apply(values)
}
```

Finite-input validation is a cross-cutting rule that applies identically
to any operation, so it lives here, once — not duplicated per operation
the way the old `CalculatorService.Add`/`.Subtract`/etc. methods each
repeated it.

### The HTTP layer builds routes from the registry, generically

```go
func NewMux(exec Executor, ops []operations.Operation, limiter RateLimiter) *http.ServeMux {
    mux := http.NewServeMux()
    mux.HandleFunc("GET /api/health", handleHealth)
    for _, op := range ops {
        handler := http.Handler(makeOperationHandler(exec, op))
        if limiter != nil {
            handler = rateLimitMiddleware(limiter, clientIPKey, handler)
        }
        mux.Handle("POST /api/"+op.Name(), handler)
    }
    return mux
}
```

One handler function (`makeOperationHandler`) serves every operation: it
decodes exactly the JSON keys `op.Operands()` names (rejecting anything
more or fewer — this is what used to be per-operation `DisallowUnknownFields`
structs), calls `exec.Execute(op.Name(), values)`, and writes the result.
Nothing here changes when an operation is added, because nothing here
knows what any specific operation does.

**Preserving the exact API contract** was the constraint throughout: the
`docs/api.md` request/response shapes, field names (`a`/`b`, `base`/
`exponent`, `value`, `value`/`percent`), error message strings, and
status codes (`400`/`404`/`405`, now also `429`) are byte-for-byte
identical to before the refactor — verified with the full existing test
suite (rewritten against the new internals, but asserting the same
request/response behavior) plus a manual pass against a running
container. This was a pure internal restructuring; no frontend change
was needed or made.

### Error handling

Every error the API can produce is still a plain Go sentinel error
(`operations.ErrDivisionByZero`, `operations.ErrNegativeSqrt`,
`operations.ErrUndefinedResult`, `service.ErrInvalidInput`,
`service.ErrUnknownOperation`, plus the HTTP layer's own
"invalid request body" and "rate limit exceeded..."), each now owned by
the file/layer that raises it rather than centralized in one
`calculator` package — which is itself a small SRP win (a package no
longer has to know about error cases that belong to operations it
doesn't otherwise touch). They all funnel through the same
`writeError`/`writeJSON` helpers into the one JSON shape,
`{"error": "<message>"}`, so "consistent JSON error responses" holds for
domain errors, validation errors, and rate-limit errors alike. A
heavier error taxonomy (error codes, structured `Kind` enums, an
error → status-code registry) was considered and left out — see
"Alternatives" below.

Business logic (`internal/operations`, `internal/service`) imports
nothing from `net/http` or `encoding/json`; it's fully testable (and
tested — see `apps/backend/tests/operations`, `tests/service`) without an
HTTP server in the loop.

### Rate limiting

`internal/ratelimit.IPRateLimiter` wraps one
[`golang.org/x/time/rate.Limiter`](https://pkg.go.dev/golang.org/x/time/rate)
(a token bucket) per key, with lazy eviction of keys idle for over 3
minutes so a long-running server doesn't accumulate one bucket per
client IP forever. `internal/api` never imports this concrete type — it
depends on its own 1-method `RateLimiter` interface, and a middleware
function wraps every operation route (not `/api/health`, which has no
abuse potential and may be polled by infrastructure) with it:

```go
func rateLimitMiddleware(limiter RateLimiter, keyFunc func(*http.Request) string, next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if !limiter.Allow(keyFunc(r)) {
            w.Header().Set("Retry-After", "60")
            writeError(w, http.StatusTooManyRequests, "rate limit exceeded, try again later")
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

Configured via `RATE_LIMIT_RPM` (default 60) and `RATE_LIMIT_BURST`
(default 10) environment variables, read in `main.go` and converted to
the token bucket's requests-per-second rate.

## The rate-limiting decision specifically: why IP, not auth

**There is no authentication anywhere in this app.** The frontend is a
static SPA with no login, no session, no API key — every request to
`/api/*` is anonymous. That rules out the strongest, most correct form of
abuse protection (per-user or per-API-key quotas) as simply not
implementable today: there is no user or key to key a quota on. The
choices actually available are:

1. **No rate limiting at all.** Simplest, but leaves the calculation
   endpoints (and, transitively, the server they share a process with)
   open to unlimited request volume from a single source — the task
   explicitly asked for "a minimum level of protection," and doing
   nothing doesn't meet that bar.
2. **IP-based rate limiting** (chosen). The connecting IP is the only
   piece of client identity available without adding auth. It's not a
   strong identity — see limitations below — but it's a real, immediately
   available signal that meaningfully raises the cost of casual abuse
   (a script hammering the API from one machine/IP) without requiring
   any change to the frontend or a login flow that doesn't otherwise
   exist in this app.
3. **Global rate limiting** (one shared bucket for the whole server,
   no per-client key). Simpler than per-IP, but one aggressive client
   would exhaust the budget for every other client sharing the server —
   worse fairness than per-IP for no implementation savings big enough
   to justify it.
4. **Add authentication now, rate-limit per user/API key.** The
   strongest option, and where this should eventually go, but it's a
   materially bigger change (login flow, credential storage, frontend
   changes) than "add a minimum abuse guard" calls for right now.
   Rejected for *this* change, not rejected forever — see below.

IP-based limiting was picked as the best available option given the
actual constraint (no auth exists), not as the ideal option in the
abstract.

### Limitations (worth being explicit about)

- **IP is a weak, shared, and spoofable-adjacent identity.** Multiple
  users behind the same NAT/corporate network/VPN share one IP and share
  one bucket; a single motivated abuser can rotate IPs (residential
  proxies, cloud instances) to get a fresh bucket per request. This is a
  **deterrent against casual abuse**, not a hard security boundary.
- **`clientIPKey` trusts `r.RemoteAddr` directly** (see
  `internal/api/middleware.go`), which is correct only because this app
  is reached directly — one container, no reverse proxy in front of it
  (see `docs/architecture.md`). If a reverse proxy or load balancer is
  introduced later, `r.RemoteAddr` would become the proxy's IP for every
  request (collapsing all clients into one bucket) unless `clientIPKey`
  is changed to read `X-Forwarded-For`/`X-Real-IP` — and that header must
  only be trusted when it's known to come from that trusted proxy
  (an untrusted client can set arbitrary `X-Forwarded-For` values to
  spoof their rate-limit key otherwise).
- **In-memory only.** `IPRateLimiter`'s buckets live in one process's
  memory; running multiple backend replicas behind a load balancer would
  give each replica its own independent limit (an attacker spread across
  replicas effectively gets `limit × replica count`). Fine for this
  project's single-container deployment; would need a shared store
  (Redis, etc.) to hold under horizontal scaling.

### How this evolves toward authenticated rate limiting

The design's whole shape is aimed at making that upgrade additive, not a
rewrite:

- `api.RateLimiter` is an interface (`Allow(key string) bool`). Swapping
  `IPRateLimiter` for a limiter with different backing storage (e.g.
  Redis-backed, for the multi-replica case) is a `main.go` wiring change
  — `internal/api` doesn't need to know.
- The **key** and the **limiter** are independent parameters to
  `rateLimitMiddleware`/`NewMux`. Moving from "one bucket per IP" to "one
  bucket per authenticated user/API key" is a matter of writing a new
  `keyFunc` (e.g. extract a validated API key or user ID from an
  `Authorization` header) and passing it in alongside — the middleware
  itself, and every operation handler it wraps, are unchanged.
  Authenticated requests could reasonably get a materially higher limit
  than anonymous ones by choosing different `RateLimiter` instances (or
  a limiter whose `Allow` branches on whether the key looks
  authenticated) per class of request.
- Nothing about the `Operation`/registry/service refactor in this ADR
  assumes anonymity — adding auth is orthogonal to it.

## Alternatives considered (SOLID refactor)

1. **Keep one `CalculatorService` method per operation, just extract the
   repeated validation into a shared helper.** Removes the duplicated
   `math.IsNaN`/`IsInf` checks but does nothing for the Open/Closed
   problem — adding operation #8 would still mean a new DTO, a new
   handler, a new service method, and a new route registration.
   Rejected: it treats the symptom (duplication) without touching the
   cause (the switch-shaped design).
2. **A generic `Execute(name, ...)` on the service (adopted) vs. keeping
   the registry lookup in the HTTP layer and passing the resolved
   `Operation` straight to a leaner service method
   (`Execute(op Operation, values)`).** The chosen design resolves by
   *name* at the service boundary so the service's public contract stays
   string-in/float-out (matching "name" as it appears over HTTP) and the
   service, not the HTTP layer, owns "what does resolution failure mean"
   (`ErrUnknownOperation`). Passing a pre-resolved `Operation` in would
   leak the registry lookup into the HTTP layer's responsibilities.
   Rejected in favor of keeping resolution as part of "execute this
   named operation," which is what the service is for.
3. **A heavier structured-error system** (a `Kind` enum, an
   `error → HTTP status` registry keyed by type, per-error HTTP status
   methods). Would make "which errors map to which status code"
   fully declarative instead of the current flat "everything from
   decode/service is 400, unknown route/method are 404/405 via
   `http.ServeMux`, rate limit is 429" rule. Rejected for now: every
   error this API currently produces really is a 400 (client-caused) or
   comes from `http.ServeMux` itself, so a richer mapping mechanism has
   no error to exercise it — it would be speculative generality ahead of
   an actual need. Revisit if/when an error needs a different status
   than its neighbors.
4. **Code-generate the per-operation boilerplate** (a script or
   `go:generate` directive producing `add.go`, `subtract.go`, etc. from a
   declarative list). Each operation file is 8–20 lines; generation would
   add build tooling to save less code than the tooling itself costs.
   Rejected as premature for seven small, rarely-changing operations.

## Consequences

- Adding operation #8 (e.g. `modulo`) is: one new file implementing
  `Operation` in `internal/operations`, one line in `Default()`, one new
  test file in `tests/operations`. `internal/service` and `internal/api`
  — and their tests — need zero changes.
- `internal/api` now imports neither `internal/service` nor
  `internal/ratelimit` — only `internal/operations` (for the `Operation`
  type passed into `NewMux`) and its own two small interfaces. Dependency
  direction is enforced by the compiler, not just convention.
- The old `internal/calculator` package is gone — its math and errors
  now live with the operation that owns them (`divide.go` owns
  `ErrDivisionByZero`, `sqrt.go` owns `ErrNegativeSqrt`, etc.), and
  `tests/calculator` became `tests/operations` to match (see
  [0002](0002-separate-test-code-from-production-code.md)).
- New test surface: `tests/operations` (one file per operation plus
  registry tests — duplicate-name rejection, `Get`/`All`/`Default`
  coverage), `tests/ratelimit` (the token bucket in isolation), and
  `tests/api/middleware_test.go` (429 behavior, health-check exemption,
  IP-key extraction, and an integration test through the real
  `ratelimit.IPRateLimiter`) alongside the existing (updated)
  `tests/service` and `tests/api/handler_test.go`.
- One new runtime dependency: `golang.org/x/time/rate` (pinned to v0.11.0
  — v0.15.0 requires Go 1.25, newer than this project's `golang:1.23-alpine`
  Docker base image).
- Two new environment variables (`RATE_LIMIT_RPM`, `RATE_LIMIT_BURST`),
  documented in `docs/api.md` and the README, with safe defaults so
  deployment without them is unaffected.
