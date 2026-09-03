package orders

import (
	"log"
	"net/http"

	"github.com/rmcampos/ecom/internal/json"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

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
