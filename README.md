# order-ecommerce-api

A learning project that builds a small e-commerce order API in Go. It currently exposes a single `GET /products` endpoint backed by PostgreSQL, with code-generated queries via sqlc and a layered service/handler architecture.

## What is this?

This repository is a minimal Go HTTP API for managing products and (eventually) orders. It is intended as a hands-on exercise in building a production-ish Go service: structured logging, dependency injection, sqlc-generated repositories, migrations, and Docker-based local development.

## Features

- HTTP server using [`go-chi/chi/v5`](https://github.com/go-chi/chi)
- PostgreSQL 16 for data persistence
- SQL queries compiled to type-safe Go code with [sqlc](https://docs.sqlc.dev/)
- Goose-style SQL migrations
- Docker Compose for local database
- Sample HTTP requests in `requests.http` for the REST Client VS Code extension

## Dependencies

- [Go](https://go.dev/) 1.26.5 or later
- [PostgreSQL](https://www.postgresql.org/) 16 (via Docker)
- [sqlc](https://docs.sqlc.dev/) — generates Go from SQL
- [goose](https://github.com/pressly/goose) — runs SQL migrations
- (Optional) VS Code with `REST Client` and `SQLTools` extensions

## Tooling Setup

### 1. Install Go

Download and install the latest Go release for your platform:

```bash
# macOS / Linux (using the official installer)
# Visit https://go.dev/dl/ and download the installer for Go 1.26.5 or later.

# Debian / Ubuntu example (adjust version as needed)
wget https://go.dev/dl/go1.26.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.26.5.linux-amd64.tar.gz
```

Add Go to your `PATH` (add this to `~/.bashrc`, `~/.zshrc`, or equivalent):

```bash
export PATH=$PATH:/usr/local/go/bin
export PATH=$PATH:$(go env GOPATH)/bin
```

Verify the installation:

```bash
go version
```

### 2. Install sqlc

Install `sqlc` to generate type-safe Go from SQL:

```bash
# macOS / Linux
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Verify installation
sqlc version
```

> Make sure `$(go env GOPATH)/bin` is on your `PATH`.

### 3. Install goose

Install `goose` to manage database migrations:

```bash
# macOS / Linux
go install github.com/pressly/goose/v3/cmd/goose@latest

# Verify installation
goose -version
```

### 4. VS Code Extensions (Optional)

These extensions make it easier to interact with the project:

#### Go

- Open VS Code.
- Go to **Extensions** (`Ctrl+Shift+X` / `Cmd+Shift+X`).
- Search for `Go lang` by the Golang team.
- Click **Install**

#### REST Client

- Open VS Code.
- Go to **Extensions** (`Ctrl+Shift+X` / `Cmd+Shift+X`).
- Search for `REST Client` by Huachao Mao.
- Click **Install**.
- Open `requests.http` at the root of the project and click **Send Request** above any request.

#### SQLTools

- Open VS Code.
- Go to **Extensions** (`Ctrl+Shift+X` / `Cmd+Shift+X`).
- Search for `SQLTools` by Matheus Teixeira.
- Click **Install**.
- Add a PostgreSQL connection with these default values:
  - **Server/Host:** `localhost`
  - **Port:** `5432`
  - **Database:** `order`
  - **Username:** `order`
  - **Password:** `order`
- You can then inspect tables and run queries directly against the local database.

In case you're on **Open VSX**, try installing with:

```
# Go
vscodium --install-extension golang.go

# REST Client
vscodium --install-extension humao.rest-client

# SQLTools UI
vscodium --install-extension mtxr.sqltools

# Choose your driver, Postgres or MySQL
vscodium --install-extension mtxr.sqltools-driver-pg
vscodium --install-extension mtxr.sqltools-driver-mysql
```

## Project Structure

```
.
├── cmd/
│   ├── main.go          # Application entry point
│   └── api.go           # Router, middleware, and server setup
├── internal/
│   ├── products/
│   │   ├── handlers.go  # HTTP handlers for /products
│   │   └── service.go   # Business logic and Service interface
│   ├── json/
│   │   └── json.go      # Shared JSON response helper
│   ├── env/
│   │   └── env.go       # Environment variable helper
│   └── adapters/
│       └── postgres/
│           ├── migrations/        # Goose migration files
│           └── sqlc/              # sqlc-generated repository code
├── .env                 # Environment variables for Goose / DB
├── compose.yaml         # Docker Compose for PostgreSQL
├── sqlc.yaml            # sqlc configuration
├── requests.http        # Sample HTTP requests for REST Client
├── go.mod
└── go.sum
```

## Getting Started

### 1. Clone the repository

```bash
git clone https://lightroasted.vps-kinghost.net/rmcampos/go-cli-course.git
cd order-ecommerce-api
```

### 2. Start the database

```bash
docker compose up -d
```

This starts a PostgreSQL 16 container named `order-db` on port `5432`.

### 3. Run the migrations

```bash
set -a; source .env; set +a
goose up
```

### 4. Run the API

```bash
go run cmd/*.go
```

The server starts on `http://localhost:8080`.

### 5. Test the endpoint

Using `curl`:

```bash
curl http://localhost:8080/health
curl http://localhost:8080/products
```

Using the `requests.http` file:

- Open `requests.http` in VS Code.
- Click **Send Request** above `Get products`.

## Build

```bash
go build -o api ./cmd
```

## Regenerate SQL code

After changing `internal/adapters/postgres/sqlc/queries.sql` or adding migrations:

```bash
sqlc generate
```

## Useful Commands

| Command | Description |
| ------- | ----------- |
| `go run ./cmd` | Run the API in development |
| `go build -o api ./cmd` | Build a binary named `api` |
| `go test ./...` | Run all tests |
| `docker compose up -d` | Start the PostgreSQL container |
| `docker compose down` | Stop the PostgreSQL container |
| `sqlc generate` | Regenerate Go code from SQL |
| `goose -dir "$GOOSE_MIGRATION_DIR" postgres "$GOOSE_DBSTRING" up` | Apply migrations |
| `goose -dir "$GOOSE_MIGRATION_DIR" postgres "$GOOSE_DBSTRING" down` | Roll back migrations |

## Environment Variables

| Variable | Default | Purpose |
| -------- | ------- | ------- |
| `GOOSE_DBSTRING` | `host=localhost user=order password=order dbname=order port=5432 sslmode=disable` | PostgreSQL connection string used by the app and Goose |
| `GOOSE_DRIVER` | `postgres` | Goose driver |
| `GOOSE_MIGRATION_DIR` | `./internal/adapters/postgres/migrations` | Location of Goose migration files |

## Current Endpoints

| Method | Path | Description |
| ------ | ---- | ----------- |
| GET | `/health` | Health check, returns `hey` |
| GET | `/products` | List all products |

## Roadmap / Next Steps

- Implement `POST /products` to create products (currently only documented in `requests.http`)
- Add full CRUD for products
- Add orders and order items
- Add request validation
- Add unit and integration tests
- Add CI/CD pipeline

## License

This project is for educational purposes.
