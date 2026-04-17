package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gabrielleite03/kenjix_core/internal/dto"
	"github.com/gabrielleite03/kenjix_core/internal/service"
)

type PurchaseHandler struct {
	service *service.PurchaseService
}

func NewPurchaseHandler() *PurchaseHandler {
	return &PurchaseHandler{
		service: service.NewPurchaseService(),
	}
}

func (h *PurchaseHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input dto.PurchaseCreateDTO

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := h.service.Create(input)
	if err != nil {
		http.Error(w, errors.Join(err).Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, result)
}

func (h *PurchaseHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r.URL.Path)

	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var input dto.PurchaseUpdateDTO

	err = json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	input.ID = id

	result, err := h.service.Update(input)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *PurchaseHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r.URL.Path)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	result, err := h.service.FindByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *PurchaseHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.FindAll()
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, result)
}
