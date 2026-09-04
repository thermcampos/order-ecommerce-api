package orders

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	repo "github.com/thermcampos/ecom/internal/adapters/postgres/sqlc"
	"github.com/thermcampos/ecom/internal/platform/json"
)

// mount
func Mount(r chi.Router, db *pgx.Conn) {
	ordersService := NewService(repo.New(db), db)
	ordersHandler := NewHandler(ordersService)
	ordersHandler.RegisterRoutes(r)
}

// routes declaration
func (h *handler) RegisterRoutes(r chi.Router) {
	r.Post("/orders", h.CreateOrderHandler)
}

// constructor arguments type
type handler struct {
	service Service
}

// constructor logic
func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

// rountes implementation
func (h *handler) CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	var orderRequest createOrderParam
	if err := json.Read(r, &orderRequest); err != nil {
		log.Println(err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	createdOrder, err := h.service.CreateOrder(r.Context(), orderRequest)
	if err != nil {
		log.Println(err)
		http.Error(w, "failed to create order", http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusCreated, createdOrder)
}
