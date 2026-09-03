package products

import (
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	repo "github.com/rmcampos/ecom/internal/adapters/postgres/sqlc"
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

func (h *handler) ListProductsHandler(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.ListProducts(r.Context())
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If products is nil, return an empty array instead of null
	if products == nil {
		products = []repo.Product{}
	}

	json.Write(w, http.StatusOK, products)
}

func (h *handler) FindProductByIDHandler(w http.ResponseWriter, r *http.Request) {
	// Get from URL path parameter
	idParam := chi.URLParam(r, "id")
	if idParam == "" {
		http.Error(w, "missing product ID", http.StatusBadRequest)
		return
	}

	// parse to int64
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		http.Error(w, "invalid product ID", http.StatusBadRequest)
		return
	}

	product, err := h.service.FindProductByID(r.Context(), id)
	if err != nil {
		// check if err is "no rows in result set" error
		if err.Error() == "no rows in result set" {
			http.Error(w, "product not found", http.StatusNotFound)
			return
		}

		// otherwise, return internal server error
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, product)
}

func (h *handler) CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	var productPayload createProductDto
	if err := json.Read(r, &productPayload); err != nil {
		log.Println(err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	createdProduct, err := h.service.CreateProduct(r.Context(), productPayload)
	if err != nil {
		log.Println(err)
		// if error contains "is required" or "invalid" returns 400
		if err.Error() == "name is required" || err.Error() == "price is required" || err.Error() == "quantity is required" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, "failed to create product", http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusCreated, createdProduct)
}