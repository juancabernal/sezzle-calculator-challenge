# Calculator — Full-Stack Technical Challenge

[![CI](https://github.com/juancabernal/sezzle-calculator-challenge/actions/workflows/ci.yml/badge.svg)](https://github.com/juancabernal/sezzle-calculator-challenge/actions/workflows/ci.yml)

A full-stack calculator built with a **Go backend** (hexagonal architecture,
REST API, PostgreSQL persistence) and a **React + TypeScript frontend**,
fully containerized with Docker Compose.

Built as a technical exercise for the Software Engineer Intern role at
Sezzle.

## Table of contents

- [Architecture](#architecture)
  - [Diagrams](#diagrams)
- [Tech stack](#tech-stack)
- [Project structure](#project-structure)
- [Getting started](#getting-started)
  - [Option A: Docker (recommended)](#option-a-docker-recommended)
  - [Option B: Running locally without Docker](#option-b-running-locally-without-docker)
- [API reference](#api-reference)
- [Testing](#testing)
- [Design decisions](#design-decisions)
- [Possible future improvements](#possible-future-improvements)

## Architecture

The backend follows **hexagonal architecture** (ports & adapters). The
core business logic never depends on HTTP, a database driver, or any
external framework — those are all "adapters" plugged in at the edges.
See the [Diagrams](#diagrams) section below for the full visual
breakdown (deployment, package structure, and the conceptual view).

The frontend mirrors the same principle: components never call `fetch`
directly. A `CalculatorApi` interface sits between them, so tests can
substitute a fake implementation without mocking libraries — the same
dependency-inversion idea applied on both ends of the stack.

### Diagrams

**Deployment diagram — simplified:**

![Deployment diagram — simplified](docs/CalculatorArchitecture.png)

A high-level, C4-style container view: the user's browser loads the
frontend once, then talks to the backend directly for every
calculation, and the backend talks to PostgreSQL.

**Deployment diagram — detailed:**

![Deployment diagram — detailed](docs/ReferenceArchetype.png)

Same system, one level more detailed: exact hostnames and ports for
each container, plus an explicit note on why the browser calls the
backend directly instead of routing through the frontend container
(the frontend is a Single-Page Application — after the initial load,
all interaction is handled by JavaScript already running in the
user's browser, which is also why CORS is enabled on the backend).

**Package diagram:**

![Package diagram](docs/Package.png)

Traces the real `import` relationships between backend packages,
confirming the dependency rule holds in practice: arrows only ever
point inward, toward `domain`, never outward from it.

**Hexagonal architecture — conceptual view:**

![Hexagonal architecture — conceptual view](docs/Hexagonal.png)

A simplified version of the same idea, independent of the specific
package names: adapters (inbound and outbound) depend on the
application layer, which depends on the domain — and the domain
depends on nothing at all. That single rule (every arrow points
inward) is what makes the domain trivially testable in isolation and
safe to change without breaking anything on the outside.

## Tech stack

| Layer | Choice | Notes |
|---|---|---|
| Backend language | Go 1.27 | Standard library `net/http` only, no framework |
| Backend testing | `testing`, `net/http/httptest` | Table-driven tests |
| Database | PostgreSQL 16 | Plain SQL via `database/sql`, no ORM |
| Frontend | React 19 + TypeScript | Vite build tooling |
| Frontend testing | Vitest + React Testing Library | |
| Styling | CSS Modules (no UI library) | |
| Containerization | Docker, Docker Compose | Multi-stage builds, non-root containers |
| CI | GitHub Actions | Backend + frontend tests, lint, build, and a Docker build check on every push/PR |

## Project structure

```
sezzle/
├── backend/
│   ├── cmd/server/            # composition root (main.go)
│   ├── internal/
│   │   ├── domain/            # entities, operations, Strategy pattern
│   │   ├── application/       # use cases (CalculatorService), ports
│   │   ├── adapters/
│   │   │   ├── http/          # inbound: REST handlers, DTOs, router
│   │   │   └── persistence/
│   │   │       ├── memory/    # outbound: in-memory repository (tests)
│   │   │       └── postgres/  # outbound: production repository
│   │   └── config/            # environment-based configuration
│   ├── migrations/            # SQL schema
│   ├── Dockerfile
│   └── Makefile
├── frontend/
│   ├── src/
│   │   ├── api/                # CalculatorApi interface + HTTP client
│   │   ├── hooks/               # useCalculator (state orchestration)
│   │   ├── components/          # CalculatorForm, ResultDisplay, HistoryList
│   │   ├── utils/                # pure input validation
│   │   └── types/                 # types mirroring backend DTOs
│   └── Dockerfile
├── docker-compose.yml
└── README.md
```

## Getting started

### Option A: Docker (recommended)

This runs PostgreSQL, the Go backend, and the React frontend together,
with the database schema applied automatically. This is the fastest
way to evaluate the whole project.

**Prerequisites:** Docker Desktop.

```bash
git clone https://github.com/juancabernal/sezzle-calculator-challenge.git
cd sezzle-calculator-challenge
docker compose up --build
```

Once running:
- Frontend: http://localhost:5173
- Backend API: http://localhost:8080
- Interactive API docs (Swagger UI): http://localhost:8080/swagger/index.html
- PostgreSQL: `localhost:5432` (user/password/db: `calculator`)

Stop everything with `Ctrl+C`, or `docker compose down` from another
terminal (add `-v` to also remove the database volume and start fresh).

### Option B: Running locally without Docker

**Prerequisites:** Go 1.23+, Node.js 18+, npm.

**Backend** (uses an in-memory repository automatically when no
`DATABASE_URL` is set — no database installation required):

```bash
cd backend
go run ./cmd/server
# listening on :8080
```

**Frontend**, in a separate terminal:

```bash
cd frontend
cp .env.example .env
npm install
npm run dev
# open http://localhost:5173
```

## API reference

Interactive documentation (Swagger UI, generated from code comments
via [`swaggo/swag`](https://github.com/swaggo/swag)) is available at
`http://localhost:8080/swagger/index.html` once the backend is
running — you can try every endpoint directly from the browser. The
raw OpenAPI spec is at `backend/docs/swagger.json`.

All responses are JSON. Successful calculations return `200 OK`;
invalid input or business-rule violations (e.g. division by zero)
return `400 Bad Request` with an `error` message.

### `POST /calculate`

**Request body:**

```json
{
  "operation": "add",
  "operand_a": 10,
  "operand_b": 4
}
```

`operation` is one of: `add`, `subtract`, `multiply`, `divide`,
`power`, `sqrt`, `percentage`. `operand_b` is omitted for the unary
operation `sqrt`.

**Example — addition:**

```bash
curl -X POST http://localhost:8080/calculate \
  -H "Content-Type: application/json" \
  -d '{"operation":"add","operand_a":10,"operand_b":4}'
```

```json
{
  "id": "8f137272-c37d-45ab-85d1-7036ea3405a8",
  "operation": "add",
  "operand_a": 10,
  "operand_b": 4,
  "result": 14,
  "created_at": "2026-09-09T21:16:15Z"
}
```

**Example — square root (unary):**

```bash
curl -X POST http://localhost:8080/calculate \
  -H "Content-Type: application/json" \
  -d '{"operation":"sqrt","operand_a":9}'
```

```json
{
  "id": "22daae9e-7887-4c07-8025-6af15e829f0c",
  "operation": "sqrt",
  "operand_a": 9,
  "result": 3,
  "created_at": "2026-09-09T21:16:31Z"
}
```

**Example — error case (division by zero):**

```bash
curl -X POST http://localhost:8080/calculate \
  -H "Content-Type: application/json" \
  -d '{"operation":"divide","operand_a":10,"operand_b":0}'
```

```json
{ "error": "division by zero" }
```

Returns `400 Bad Request`.

**Example — error case (non-finite result):**

Some otherwise well-formed inputs don't have a real-number answer —
e.g. a negative base raised to a fractional power, or an exponent
large enough to overflow a `float64`. JSON has no way to represent
`NaN` or `Infinity`, so these are rejected the same way division by
zero is, instead of the API returning `200 OK` with a broken body:

```bash
curl -X POST http://localhost:8080/calculate \
  -H "Content-Type: application/json" \
  -d '{"operation":"power","operand_a":-8,"operand_b":0.5}'
```

```json
{ "error": "this operation does not produce a finite real number for the given operands" }
```

Returns `400 Bad Request`.

### `GET /history`

Returns the most recent successful calculations (up to 20), newest
first. Failed calculations are never persisted.

```bash
curl http://localhost:8080/history
```

```json
{
  "calculations": [
    {
      "id": "8f137272-c37d-45ab-85d1-7036ea3405a8",
      "operation": "add",
      "operand_a": 10,
      "operand_b": 4,
      "result": 14,
      "created_at": "2026-09-09T21:16:15Z"
    }
  ]
}
```

## Testing

**Backend** (from `backend/`):

```bash
go test ./... -v
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

Current coverage: `domain` 95.7%, `application` 94.7%, `adapters/http`
93.4%, `adapters/persistence/memory` 100%.

`adapters/persistence/postgres` has its own integration test
(`repository_integration_test.go`) that runs against a **real**
Postgres instance instead of a mock — it verifies the actual SQL
round-trips correctly (including the NULL-vs-zero distinction for
`operand_b` on unary operations, and that history comes back in the
right order). It's gated behind an environment variable so `go test
./...` never fails — or even slows down — on a machine with no
database running:

```bash
# Skipped by default:
go test ./...

# Runs against a real database (any reachable Postgres works, e.g.
# the one docker-compose already starts, or a throwaway container):
docker run -d --rm -p 5432:5432 \
  -e POSTGRES_USER=calculator -e POSTGRES_PASSWORD=calculator -e POSTGRES_DB=calculator \
  postgres:16-alpine

TEST_DATABASE_URL="postgres://calculator:calculator@localhost:5432/calculator?sslmode=disable" \
  go test ./... -v
```

CI (see [`.github/workflows/ci.yml`](.github/workflows/ci.yml)) always
runs with a Postgres service container available, so this suite runs
on every push and pull request — it isn't just documented, it's
actually exercised.

**Frontend** (from `frontend/`):

```bash
npm test              # equivalent to: vitest run
npm run test:coverage # equivalent to: vitest run --coverage
```

Current coverage: 98.7% overall (statements), 25 tests.

## Design decisions

**Hexagonal architecture, applied but not over-applied.** The backend
separates domain, application, and adapters strictly — but there's no
CQRS, no event bus, no extra layers beyond what this problem actually
needs. The goal was the right amount of structure for a small,
well-tested service, not architecture for its own sake.

**Strategy pattern for operations.** Each arithmetic operation
(`AddOperation`, `DivideOperation`, ...) implements a single
`Operation` interface. Adding an 8th operation later means adding one
new file — no existing code changes (Open/Closed Principle).

**Two repository adapters, one port.** `CalculationRepository` is
implemented both by an in-memory store (used exclusively in tests —
no mocking library needed) and by a real PostgreSQL adapter (used in
production/Docker). `main.go` is the only file that decides which one
to use, based on whether `DATABASE_URL` is set.

**No web framework on the backend.** Go 1.22+'s `net/http.ServeMux`
supports method-based routing natively, which is enough for 2 routes.
Avoiding Gin/Echo here is a deliberate choice to demonstrate standard
library fluency, not an oversight — this also aligns with Sezzle's
stated preference to "build what we can before buying."

**No global state management on the frontend.** All state lives in a
single `useCalculator` hook. Redux/Zustand would add complexity this
app's size doesn't justify — the same restraint principle applied on
both ends of the stack.

**CSS Modules, no UI component library.** Keeps the bundle small and
avoids a dependency not mentioned in the target stack, at the cost of
writing a bit more CSS by hand.

**Trunk-based Git workflow (GitHub Flow), not Git Flow.** Every
feature lived in a short-lived `feature/*` branch, merged to `main`
via a Pull Request with a written description. No long-lived
`develop` branch — chosen because Sezzle's stack (Kubernetes + GitLab
CI/CD) implies continuous deployment, where `develop` as an
intermediate branch adds ceremony without solving a real problem for
a single-contributor project.

**Validation exists on both frontend and backend.** The frontend
rejects obviously invalid input (empty fields, division by zero)
before making a network call, purely for responsiveness. The backend
validates independently and is the actual source of truth — the
frontend is never trusted as the only line of defense.

**CORS is permissive (`*`) in this project.** Acceptable for local
development and this evaluation; a real production deployment would
restrict `Access-Control-Allow-Origin` to the frontend's exact domain.

**Non-finite results are rejected centrally, not per-operation.**
`math.Pow` can return `NaN` (a negative base with a fractional
exponent — there's no real-number answer) or `+Inf`/`-Inf` (the
exponent is large enough to overflow `float64`), and `encoding/json`
cannot marshal either value. Rather than special-casing this inside
`PowerOperation` — the only strategy where it's likely to come up in
practice, but not the only one where it's theoretically possible —
`CalculatorService.Calculate` checks the result of *every* operation
once, in one place, and maps it to a normal domain error. This is the
same reasoning as the existing division-by-zero handling, applied
consistently: fail with a clear `400` and a message, never with a
`200 OK` and a broken or empty body.

**Explicit `http.Server` timeouts, and a body-size cap on
`/calculate`.** A bare `http.ListenAndServe` has no read/write/idle
timeouts at all, which leaves the server exposed to a client that
opens a connection and simply never finishes sending its request
(Slowloris-style resource exhaustion). `main.go` builds an explicit
`http.Server` with sane timeouts instead. Similarly, `HandleCalculate`
wraps the request body in `http.MaxBytesReader` — a calculation
request is at most a few dozen bytes, so capping it at 4 KiB costs
nothing legitimate and closes off a trivial way to waste server memory
and CPU decoding an oversized payload.

**A `seq BIGSERIAL` column drives history ordering, not `created_at`
alone.** This was caught by the Postgres integration test, not
guessed upfront: two calculations saved by concurrent goroutines can
land in the very same microsecond (`TIMESTAMPTZ`'s own resolution),
which makes `ORDER BY created_at DESC` alone non-deterministic on
ties. `seq` is assigned atomically and monotonically by Postgres on
insert, so it can never tie. The in-memory repository has the same
class of fix for a different reason: it now returns history by
reversing insertion order (`Save` only ever appends) instead of
sorting by `CreatedAt`, since `sort.Slice` isn't even guaranteed
stable on ties in the first place.

**Docker containers run as an unprivileged user.** Both the backend
binary and (via the official `nginx:alpine` image's own default) the
frontend's web server drop root before serving traffic — neither
process has any legitimate reason to run as root, so there's no
reason to grant it.

**Line endings are normalized via `.gitattributes`.** Go tooling
(`gofmt`, and by extension any editor or CI step that runs it)
normalizes to LF; a Windows checkout without this file can end up
with CRLF baked into the committed blobs themselves, which makes
`gofmt -l` flag every file as "unformatted" for a contributor on
Linux or macOS, and turns unrelated diffs into noisy, all-lines-changed
ones. `* text=auto eol=lf` makes every checkout — regardless of
platform — converge on the same bytes.

## Possible future improvements

Deliberately left out — not overlooked, but genuinely more than this
problem needs right now, and adding them speculatively would be
exactly the kind of unnecessary complexity the rest of this project
tries to avoid:

- Authentication / per-user calculation history
- Rate limiting on the API
- A proper migration tool (e.g. `golang-migrate`) if the schema were
  expected to evolve frequently and across multiple environments
- End-to-end tests (e.g. Playwright) covering the full Docker Compose
  stack, on top of the unit and integration tests that already exist
