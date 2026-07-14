# API Reference

Base path: `/api`

All responses are JSON. Every calculation endpoint accepts `POST` with a
JSON body and returns `{ "result": <number> }` on success, or
`{ "error": "<message>" }` on failure.

## Architecture

Requests flow through the following layers (see
[architecture.md](architecture.md)):

1. **`internal/api`** — HTTP handlers and rate-limit middleware.
   Decode/encode JSON, route to the service, map errors to status codes.
2. **`internal/service`** — validates that operands are finite numbers
   (rejects `NaN`/`Infinity`), resolves the named operation, and
   delegates to it.
3. **`internal/operations`** — one independent implementation per
   operation (add, subtract, multiply, divide, power, sqrt, percentage)
   plus the domain errors each can produce (division by zero, negative
   square root, undefined results).

## Error responses

| Status | When                                                                 |
| ------ | ---------------------------------------------------------------------- |
| `400`  | Malformed/unknown JSON fields, non-finite operand, or a domain error (division by zero, negative square root, undefined exponentiation result) |
| `404`  | Unknown route                                                           |
| `405`  | Wrong HTTP method for a known route                                    |
| `429`  | Rate limit exceeded for the client's IP (see [Rate limiting](#rate-limiting)) |

## `GET /api/health`

Liveness check.

**Response `200`**

```json
{ "status": "ok" }
```

## `POST /api/add`

Returns `a + b`.

**Request**

```json
{ "a": 2, "b": 3 }
```

**Response `200`**

```json
{ "result": 5 }
```

## `POST /api/subtract`

Returns `a - b`.

**Request**

```json
{ "a": 5, "b": 3 }
```

**Response `200`**

```json
{ "result": 2 }
```

## `POST /api/multiply`

Returns `a * b`.

**Request**

```json
{ "a": 4, "b": 3 }
```

**Response `200`**

```json
{ "result": 12 }
```

## `POST /api/divide`

Returns `a / b`.

**Request**

```json
{ "a": 9, "b": 3 }
```

**Response `200`**

```json
{ "result": 3 }
```

**Response `400` — division by zero**

```json
{ "error": "division by zero" }
```

## `POST /api/power`

Returns `base` raised to the power of `exponent`.

**Request**

```json
{ "base": 2, "exponent": 10 }
```

**Response `200`**

```json
{ "result": 1024 }
```

**Response `400` — undefined result**

Zero raised to a negative exponent:

```json
{ "error": "division by zero" }
```

A negative base with a fractional exponent (no real result):

```json
{ "error": "operation has no defined real result" }
```

## `POST /api/sqrt`

Returns the square root of `value`.

**Request**

```json
{ "value": 9 }
```

**Response `200`**

```json
{ "result": 3 }
```

**Response `400` — negative operand**

```json
{ "error": "cannot take the square root of a negative number" }
```

## `POST /api/percentage`

Returns `percent`% of `value`, i.e. `(value * percent) / 100`.

**Request**

```json
{ "value": 200, "percent": 10 }
```

**Response `200`**

```json
{ "result": 20 }
```

## Common error cases

**Malformed JSON body / unknown fields — `400`**

```json
{ "error": "invalid request body" }
```

**Non-finite operand (`NaN` / `Infinity`) — `400`**

```json
{ "error": "input must be a finite number" }
```

## Rate limiting

Every operation endpoint above (everything except `GET /api/health`) is
rate-limited **per client IP**, using a token-bucket limiter
([`golang.org/x/time/rate`](https://pkg.go.dev/golang.org/x/time/rate)).
There's no authentication in this app, so IP is the only client identity
available — see
[docs/adr/0003-solid-refactor-and-rate-limiting.md](adr/0003-solid-refactor-and-rate-limiting.md)
for the reasoning, limitations, and how this evolves if auth is added.

Configured via environment variables on the server (defaults shown):

| Variable            | Default | Meaning                                      |
| -------------------- | ------- | ----------------------------------------------- |
| `RATE_LIMIT_RPM`      | `60`    | Requests per minute allowed per IP, steady-state |
| `RATE_LIMIT_BURST`    | `10`    | Requests a single IP may fire instantaneously before the per-minute rate applies |

**Response `429` — rate limit exceeded**

Includes a `Retry-After` header (seconds); body uses the same shape as
every other error:

```json
{ "error": "rate limit exceeded, try again later" }
```
