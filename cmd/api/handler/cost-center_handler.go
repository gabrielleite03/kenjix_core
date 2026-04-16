package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gabrielleite03/kenjix_core/internal/dto"
	"github.com/gabrielleite03/kenjix_core/internal/service"
)

type CostCenterHandler struct {
	service service.CostCenterService
}

func NewCostCenterHandler(costCenterService service.CostCenterService) *CostCenterHandler {
	return &CostCenterHandler{
		service: costCenterService,
	}
}

func (h *CostCenterHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input dto.CostCenterCreateDTO

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := h.service.Create(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, result)
}

func (h *CostCenterHandler) Update(w http.ResponseWriter, r *http.Request) {
	var input dto.CostCenterUpdateDTO

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := h.service.Update(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *CostCenterHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r.URL.Path)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	result, err := h.service.FindByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *CostCenterHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.FindAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *CostCenterHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r.URL.Path)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
