package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	repo "github.com/rmcampos/ecom/internal/adapters/postgres/sqlc"
	"github.com/rmcampos/ecom/internal/orders"
	"github.com/rmcampos/ecom/internal/products"
)

type application struct {
	config config
	db     *pgx.Conn
}

type dbConfig struct {
	dsn string // database connection string
}

type config struct {
	addr string // 12 factor document
	db   dbConfig
}

// mount
func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)              // important for rate-limiting
	r.Use(middleware.ClientIPFromRemoteAddr) // also important for rate-limiting, analytics and tracing
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer) // nice to recover from panics and crashes

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("OK")); err != nil {
			slog.Error("Failed to write health response", "error", err)
		}
	})

	productsService := products.NewService(repo.New(app.db))
	productsHandler := products.NewHandler(productsService)
	r.Get("/products", productsHandler.ListProductsHandler)
	r.Get("/products/{id}", productsHandler.FindProductByIDHandler)
	r.Post("/products", productsHandler.CreateProductHandler)

	ordersService := orders.NewService(repo.New(app.db), app.db)
	ordersHandler := orders.NewHandler(ordersService)
	r.Post("/orders", ordersHandler.CreateOrderHandler)

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
