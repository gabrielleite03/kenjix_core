package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gabrielleite03/kenjix_core/internal/service"
)

type ProductHandler struct {
	ProductService service.ProductService
}

type ProductService interface {
	ListProducts() error
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{ProductService: productService}
}

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	products, err := h.ProductService.ListProducts(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(products); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
