package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gabrielleite03/kenjix_core/internal/service"
)

type ProductHandler struct {
	ProductService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{ProductService: productService}
}

// GET /products
func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	products, err := h.ProductService.ListProducts(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, products)
}

// GET /products/{id}
func (h *ProductHandler) Get(w http.ResponseWriter, r *http.Request) {
	// implementar
}

// POST /products
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	// implementar
}

// PUT /products/{id}
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	// implementar
}

// DELETE /products/{id}
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// implementar
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
