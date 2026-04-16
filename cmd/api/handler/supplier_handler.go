package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gabrielleite03/kenjix_core/internal/dto"
	"github.com/gabrielleite03/kenjix_core/internal/service"
)

type SupplierHandler struct {
	SupplierService service.SupplierService
}

func NewSupplierHandler(supplierService service.SupplierService) *SupplierHandler {
	return &SupplierHandler{SupplierService: supplierService}
}

// GET /products
func (h *SupplierHandler) List(w http.ResponseWriter, r *http.Request) {

	suppliers, err := h.SupplierService.GetAll()

	if err != nil {
		println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, suppliers)
}

// GET /products/{id}
func (h *SupplierHandler) Get(w http.ResponseWriter, r *http.Request) {
	println("Handling get Supplier request")

	id, err := extractID(r.URL.Path)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	supplier, err := h.SupplierService.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, supplier)

}

func (h *SupplierHandler) Create(w http.ResponseWriter, r *http.Request) {

	println("Handling create category request")

	var s dto.SupplierDTO
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	created, err := h.SupplierService.Create(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

// PUT /products/{id}
func (h *SupplierHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r.URL.Path)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var s dto.SupplierDTO
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.ID = id

	err = h.SupplierService.Update(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, err)

}

// DELETE /products/{id}
func (h *SupplierHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// implementar
}
