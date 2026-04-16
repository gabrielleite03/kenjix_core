package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gabrielleite03/kenjix_core/internal/dto"
	"github.com/gabrielleite03/kenjix_core/internal/service"
)

type CategoryHandler struct {
	CategoryService service.CategoryService
}

func NewCategoryHandler(categoryService service.CategoryService) *CategoryHandler {
	return &CategoryHandler{CategoryService: categoryService}
}

// GET /categories
func (h *CategoryHandler) List(w http.ResponseWriter, r *http.Request) {
	categories, err := h.CategoryService.List(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, categories)
}

// GET /categories/{id}
func (h *CategoryHandler) Get(w http.ResponseWriter, r *http.Request) {

	id, err := extractID(r.URL.Path)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	category, err := h.CategoryService.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, category)
}

// POST /categories
func (h *CategoryHandler) Create(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	var d dto.CategoryDTO

	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	created, err := h.CategoryService.Create(r.Context(), &d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

// PUT /categories/{id}
func (h *CategoryHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r.URL.Path)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var d dto.CategoryDTO
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	d.ID = id

	updated, err := h.CategoryService.Update(r.Context(), &d)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

// DELETE /categories/{id}
func (h *CategoryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r.URL.Path)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.CategoryService.Delete(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func extractID(path string) (int64, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	return strconv.ParseInt(parts[len(parts)-1], 10, 64)
}
