# Agent Guide for `order-ecommerce-api`

## Project Type & Stack

This is a **Go 1.26.5 HTTP API** for an order/e-commerce backend.

- **Router:** [`go-chi/chi/v5`](https://github.com/go-chi/chi)
- **Database driver:** [`jackc/pgx/v5`](https://github.com/jackc/pgx) (direct `pgx.Conn`)
- **SQL/codegen:** [sqlc](https://docs.sqlc.dev/) v1.31.1, configured for `pgx/v5`
- **Database:** PostgreSQL 16 (started via Docker Compose)
- **Migrations:** Goose-style SQL files (`.env` references `GOOSE_*`, but no Goose CLI is configured yet)
- **Logging:** `log/slog` with JSON output

## Repository Layout

```
.
├── cmd/
│   ├── main.go          # Entry point: loads config, connects DB, runs server
│   └── api.go           # `application` struct, router setup, HTTP server
├── internal/
│   ├── products/
│   │   ├── handlers.go  # HTTP handlers for /products
│   │   └── service.go   # Business logic + Service interface
│   ├── json/
│   │   └── json.go      # Shared JSON response helper
│   ├── env/
│   │   └── env.go       # Environment variable helper
│   └── adapters/
│       └── postgres/
│           ├── migrations/
│           │   └── 00001_create_products.sql
│           └── sqlc/
│               ├── queries.sql       # Source SQL for sqlc
│               ├── queries.sql.go    # Generated query methods
│               ├── querier.go        # Generated Querier interface
│               ├── db.go             # Generated DBTX and Queries constructor
│               └── models.go         # Generated model structs
├── .env                 # Declares DB env vars for Goose/Postgres
├── compose.yaml         # Docker Compose for Postgres
├── sqlc.yaml           # sqlc configuration
├── requests.http       # VS Code REST Client sample requests
├── go.mod
└── go.sum
```

## Essential Commands

There is no `Makefile` or CI config; use standard Go tooling directly.

### Run

```bash
# 1. Start Postgres
sudo docker compose up -d

# 2. Run the API (defaults to :8080)
go run ./cmd
```

### Build

```bash
go build -o api ./cmd
```

### Test

```bash
# There are currently no *_test.go files in the repo.
go test ./...
```

### Database

```bash
# Start the database container
sudo docker compose up -d

# Run migrations manually with goose (install goose if needed)
# The .env file is set up for goose, but the binary is not a project dependency yet.
source .env
goose -dir "$GOOSE_MIGRATION_DIR" postgres "$GOOSE_DBSTRING" up
```

### Generate SQL code

```bash
sqlc generate
```

> Regenerate after any change to `internal/adapters/postgres/sqlc/queries.sql` or `internal/adapters/postgres/migrations/*.sql`.

## Application Architecture & Data Flow

### Layered design

1. **`cmd/`** — wiring and process start
   - `main.go`: reads env/config, opens a single `pgx.Conn`, and starts the HTTP server.
   - `api.go`: owns the `application` struct, mounts routes, applies middleware, and runs `http.Server`.

2. **`internal/products/`** — feature package
   - `service.go`: defines `Service` interface and `svc` implementation. The service depends only on the generated `repo.Querier` interface.
   - `handlers.go`: HTTP handlers depend on `Service`. Converts nil slices to empty arrays before serializing.

3. **`internal/adapters/postgres/sqlc/`** — generated data access layer
   - `repo.New(app.db)` produces a `*Queries` that satisfies `repo.Querier`.
   - All SQL lives in `queries.sql`; generated Go should not be edited by hand.

4. **`internal/json/`** — shared response helper
   - `json.Write(w, status, data)` sets `Content-Type: application/json`, writes the status, and encodes the payload.

5. **`internal/env/`** — tiny helper for env vars with fallbacks.

### Request lifecycle

```
HTTP Request
    ↓
chi router (middleware: RequestID, ClientIPFromRemoteAddr, Logger, Recoverer, Timeout)
    ↓
products.handler.ListProductsHandler
    ↓
products.Service.ListProducts(ctx)
    ↓
repo.Querier.ListProducts(ctx)
    ↓
PostgreSQL (pgx)
```

## Configuration & Environment

### Environment variables

| Variable         | Default                                                                      | Used for                     |
| ---------------- | ---------------------------------------------------------------------------- | ---------------------------- |
| `GOOSE_DBSTRING` | `host=localhost user=order password=order dbname=order port=5432 sslmode=disable` | Postgres DSN in `main.go`   |

### Hardcoded server address

- Server listens on `:8080` (`main.go` → `cmd/api.go`).

## Code Conventions & Patterns

### Naming

- Generated repository package is imported as `repo` (`import repo "github.com/thermcampos/ecom/internal/adapters/postgres/sqlc"`).
- Service type is an interface named `Service`; concrete implementation is `svc` (private).
- Handlers return a `*handler` via `NewHandler` constructor.
- Use `New...` constructors for wiring dependencies.

### Style

- No `Makefile` or CI; rely on `go` and `sqlc` CLI commands.
- Structured logging with `slog` is initialized in `main.go`.
- Handler errors are currently logged with `log.Println(err)` and returned as plain text `http.Error`.
- Money is stored as `INTEGER` cents (`price_in_cents`); no decimal type is used.
- `TIMESTAMPTZ` columns default to `CURRENT_TIMESTAMP`.

### Database

- Migrations are Goose-style (`-- +goose Up` / `-- +goose Down`) in `internal/adapters/postgres/migrations/`.
- sqlc reads the migration directory as the schema source and `queries.sql` as the query source.
- sqlc config (`sqlc.yaml`):
  - Engine: `postgresql`
  - Generated package: `repo`
  - Output: `internal/adapters/postgres/sqlc`
  - `sql_package: pgx/v5`
  - `emit_json_tags: true`
  - `emit_interface: true`

## API Endpoints

Currently implemented routes (from `cmd/api.go`):

- `GET /health` — returns plain text `hey`
- `GET /products` — lists all products as JSON

The `requests.http` file also references a `POST /products` endpoint, but it is **not implemented yet**.

## Testing Approach

- **There are currently no `*_test.go` files.**
- The service layer depends on `repo.Querier`, an interface generated by sqlc, so unit tests can inject a fake/mocked querier.
- Handler tests can use `httptest.NewRecorder` and a stubbed `Service`.

## Important Gotchas & Non-Obvious Patterns

1. **Do not edit generated files.** Anything in `internal/adapters/postgres/sqlc/` is generated by `sqlc`. Modify `queries.sql` and re-run `sqlc generate`.

2. **`GOOSE_DBSTRING` is the actual runtime DSN.** Although `.env` references Goose migration variables, the application itself reads `GOOSE_DBSTRING` via `env.GetString` in `main.go`. The fallback DSN matches the Docker Compose Postgres credentials.

3. **`/products` returns an empty array, never `null`.** The handler explicitly converts a nil slice to `[]repo.Product{}` before writing JSON.

4. **`ListProducts` can return `nil, err`.** If the database query fails, the error is logged and a 500 is returned.

5. **No request/response validation library is in use.** JSON decoding is not yet implemented in handlers; expect to add it manually if you add write endpoints.

6. **No authentication or authorization middleware is present.** Any new endpoints are currently public by default.

7. **No `Makefile` or task runner.** Every command is a direct `go`, `sqlc`, `docker compose`, or `goose` command.

8. **The `requests.http` file is stale.** It includes a `POST /products` example that does not yet have a corresponding handler or database query.

## When Adding a New Feature

1. Add/change SQL in `internal/adapters/postgres/sqlc/queries.sql`.
2. Add migrations in `internal/adapters/postgres/migrations/` if the schema changes.
3. Run `sqlc generate`.
4. Update or add a service in `internal/<feature>/service.go`.
5. Update or add handlers in `internal/<feature>/handlers.go`.
6. Wire the route in `cmd/api.go`.
7. Add tests.
