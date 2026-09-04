package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/thermcampos/ecom/internal/domain/orders"
	"github.com/thermcampos/ecom/internal/domain/products"
	"github.com/thermcampos/ecom/internal/platform/health"
	"github.com/thermcampos/ecom/internal/platform/logging"
)

// constructor arguments type
type application struct {
	config config
	db     *pgx.Conn
}

// nested type for db config, in the config argument of application
type dbConfig struct {
	dsn string
}

// config type, argument for the application constructor
type config struct {
	addr string
	db   dbConfig
}

// mount
func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)              // important for rate-limiting
	r.Use(middleware.ClientIPFromRemoteAddr) // also important for rate-limiting, analytics and tracing
	r.Use(middleware.RequestLogger(logging.NewFormatter()))
	r.Use(middleware.Recoverer) // nice to recover from panics and crashes

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	health.Mount(r)

	// Domains registers
	products.Mount(r, app.db)
	orders.Mount(r, app.db)

	return r
}

// run
func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	slog.Info("Starting server", "port", app.config.addr)

	return srv.ListenAndServe()
}
