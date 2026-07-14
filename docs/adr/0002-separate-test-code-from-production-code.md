# 0002. Separate test code from production code

## Status

Accepted

## Context

Both apps previously colocated tests with the code they exercised:

- Backend: `internal/calculator/calculator_test.go` sat next to
  `calculator.go`, same package, same directory — same for
  `internal/service` and `internal/api`.
- Frontend: `src/features/calculator/**/*.test.ts(x)` sat next to the
  module under test, and `src/setupTests.ts` (test environment setup)
  lived directly in `src/`.

This is Go's and Vitest's default/idiomatic layout, and it worked, but it
mixes two different audiences in the same directory listing: someone
reading `internal/api/` to understand the HTTP layer also has to step
over `handler_test.go`; someone auditing what ships in the Docker image
has to mentally filter out every `*.test.tsx` from `src/`. As the
calculator UI grew a real state machine and a multi-file keypad
implementation (`engine/`, `hooks/`, `components/`), the test-to-source
file ratio in some directories exceeded 1:1, and it became the deciding
factor to separate the two before it got worse.

## Decision

Move all test code into a dedicated `tests/` directory per app, mirroring
the production package/feature structure, and keep production source free
of test files.

### Backend: `apps/backend/tests/`, external test packages only

```
apps/backend/
├── internal/
│   ├── calculator/calculator.go
│   ├── service/service.go
│   └── api/handler.go, types.go
└── tests/
    ├── calculator/calculator_test.go   (package calculator_test)
    ├── service/service_test.go        (package service_test)
    └── api/handler_test.go            (package api_test)
```

Go ties a test file's access level to its package declaration, not its
directory: a whitebox test (`package calculator`, in the same directory
as the code, seeing unexported identifiers) **cannot** be relocated to a
different directory — Go's package system is directory-based. The only
way to physically separate test files from source is external
("blackbox") testing: `package calculator_test`, living wherever you like,
importing the production package and using only its exported API. Since
`tests/` sits alongside (not inside) `internal/`, and Go's `internal/`
visibility rule allows import by anything rooted at `internal`'s parent
(the module root), `calculator-backend/tests/calculator` can import
`calculator-backend/internal/calculator` without issue.

One concrete consequence: `apps/backend/internal/api/handler_test.go`
used to construct requests directly from the package's own unexported DTOs
(`binaryRequest{A: 2, B: 3}`). The relocated `tests/api/handler_test.go`
instead builds `map[string]float64{"a": 2, "b": 3}` and defines its own
local `resultBody`/`errorBody` structs matching the documented JSON shape
(`docs/api.md`). That's not a workaround — it's the actual improvement:
the test now asserts against the JSON contract the frontend and any other
client depend on, not against a Go type that's free to change internally
as long as the contract holds.

Coverage requires one flag change: with tests relocated, `go test
./tests/... -coverprofile=...` alone would only measure coverage of the
`tests/` packages themselves (near-zero, uninteresting). The Makefile's
`coverage-backend` target now passes `-coverpkg=./internal/...` explicitly,
telling Go which packages to instrument even though the tests exercising
them live elsewhere — this is the standard Go pattern for external test
suites and is exactly what "coverage of production code, driven by tests
elsewhere" requires.

### Frontend: `apps/frontend/tests/`, importing via the `@/` alias

```
apps/frontend/
├── src/features/calculator/{api,engine,hooks,validation,components}/*.ts(x)
└── tests/
    ├── setup.ts
    └── features/calculator/{api,engine,hooks,validation,components}/*.test.ts(x)
```

`tests/` mirrors `src/features/calculator/` exactly, so "where's the test
for X" is always "the same path under `tests/`." Every test imports
production code via the existing `@/` alias (e.g.
`@/features/calculator/engine/format`) instead of relative paths like
`./format`. That's what makes moving a test file a non-event for its
imports (the alias doesn't care where the test physically lives) and
what guarantees the dependency only runs one way: `tests/` imports `src/`,
never the reverse, so there's no circular-dependency risk between the two
trees. `vite.config.ts`'s `test.include` now points at `tests/**/*.{test,spec}.{ts,tsx}`
and `test.setupFiles` at `tests/setup.ts`; `test.coverage.include` is
pinned to `src/**/*.{ts,tsx}` so coverage reports what ships, not the
tests themselves.

### Tooling touched

- `apps/backend/tests/**` — new test files (blackbox packages).
- `apps/frontend/tests/**` — new test files + `setup.ts`.
- `apps/frontend/vite.config.ts` — `test.include`, `test.setupFiles`,
  `test.coverage.include`/`exclude`.
- `apps/frontend/tsconfig.json` — `include` now covers `["src", "tests"]`
  so `tsc --noEmit` (part of `npm run build`) still typechecks tests.
- Root `Makefile` — `coverage-backend` gained `-coverpkg=./internal/...`;
  every other target (`test`, `test-backend`, `test-frontend`,
  `coverage-frontend`, `build`, `up`, `down`, `clean`) is unchanged text,
  because `go test ./...` and `npm test`/`npm run coverage` already
  discover the relocated tests without configuration.
- `.dockerignore` — excludes `apps/backend/tests` and
  `apps/frontend/tests` (plus `coverage.out`/`coverage/`) from the
  production image's build context. Neither `go build` nor `vite build`
  touches test code, so this only shrinks/stabilizes the build context
  (test-only changes no longer bust that Docker layer's cache); the dev
  Dockerfiles (`Dockerfile.dev`) still `COPY . .` and keep `tests/`
  available for running tests inside a dev container.

## Alternatives considered

1. **Keep tests colocated (status quo).**
   Zero migration cost, and it's the ecosystem default for a reason (test
   right next to what it tests is easy to find by proximity). Rejected
   because it was explicitly the thing being asked to change, and because
   as this project grows past a single feature, colocated test files are
   the main source of directory clutter when scanning "what does this
   package actually do" versus "how is it verified."

2. **A single repo-root `tests/` covering both frontend and backend.**
   Would centralize "where are all the tests" into one place, but Go
   tests and Vitest tests have nothing in common (different toolchains,
   different module roots, different CI steps) — merging them under one
   directory buys a shorter README sentence at the cost of making each
   app's own test suite harder to run in isolation (`go test ./...` from
   a repo-root `tests/` would need to reach back across the whole
   monorepo). Rejected in favor of one `tests/` per app, consistent with
   `apps/backend` and `apps/frontend` already being independent build
   units.

3. **Backend: keep whitebox tests, just move the file and fix the
   package name to match the new directory via a build tag or symlink.**
   Not viable: Go has no mechanism to make a `_test.go` file whitebox
   (access unexported identifiers) for a package whose source lives in a
   different directory. Any real separation requires external
   (`package foo_test`) testing, which was chosen deliberately — see
   "Decision" above — not as a fallback.

4. **Frontend: keep relative imports from `tests/` back into `src/`
   (e.g. `../../../src/features/calculator/engine/format`).**
   Works, but produces exactly the kind of fragile, rename-sensitive path
   that colocated tests never had to deal with, and the codebase already
   has a clean alias (`@/`) built for this. Rejected in favor of `@/`
   imports, which also make it visually obvious in a diff that a test
   file is importing FROM production code, never the other way around.

5. **Split coverage into two profiles (one for `tests/`'s own coverage,
   one `-coverpkg` profile for production code) and report both.**
   The `tests/` packages have no interesting logic of their own (just
   assertions and fixtures), so a coverage number for them is noise.
   Rejected — only the `-coverpkg=./internal/...` profile is generated
   and reported, matching what the previous colocated setup effectively
   measured.

## Consequences

- **Discoverability**: `internal/` and `src/` now contain only the code
  that ships; `tests/` is the one place to look for "how is this
  verified," mirroring the production layout 1:1 in both apps.
- **Public-interface discipline (backend)**: external test packages
  physically cannot reach into unexported internals, so a test that
  starts depending on an implementation detail is a compile error, not a
  code-review nit. This is a real constraint gained, not just a cosmetic
  move — see the `handler_test.go` JSON-body rewrite above.
- **No new indirection cost (frontend)**: because the `@/` alias already
  existed, tests reading `@/features/calculator/...` are exactly as clear
  as they were reading `./...` — arguably clearer, since the import now
  states the full feature path instead of relying on directory context.
- **CI/tooling surface grew by one flag**: `-coverpkg=./internal/...` on
  the backend coverage command. Everything else a CI pipeline would call
  (`go test ./...`, `npm test`, `npm run coverage`, `make test`, `make
  coverage`) is textually unchanged.
- **Scaling**: adding a second feature (beyond `calculator`) means adding
  a matching subtree under both `tests/` directories — the mirrored
  structure means there's no ambiguity about where its tests belong.
