package health

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/thermcampos/ecom/internal/platform/json"
)

// mount
func Mount(r chi.Router) {
	healthHandler := NewHandler()
	healthHandler.RegisterRoutes(r)
}

// routes declaration
func (h *handler) RegisterRoutes(r chi.Router) {
	r.Get("/health", h.GetHealth)
	r.Get("/ready", h.GetLive)
	r.Get("/live", h.GetLive)
}

// constructor arguments type
type handler struct{}

// constructor logic
func NewHandler() *handler {
	return &handler{}
}

// rountes implementation
func (h *handler) GetHealth(w http.ResponseWriter, r *http.Request) {
	msg := healthCheck{
		Message: "OK",
	}
	json.Write(w, http.StatusOK, msg)
}

func (h *handler) GetLive(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte(nil)); err != nil {
		slog.Error("Failed to write ready/live response", "error", err)
	}
}
