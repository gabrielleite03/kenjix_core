package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gabrielleite03/kenjix_core/internal/dto"
	"github.com/gabrielleite03/kenjix_core/internal/service"
)

type MarketplaceHandler struct {
	service service.MarketplaceService
}

func NewMarketplaceHandler(service service.MarketplaceService) *MarketplaceHandler {
	return &MarketplaceHandler{service: service}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{
		"error": message,
	})
}

func (h *MarketplaceHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input dto.MarketplaceDTO

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}

	result, err := h.service.Create(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, result)
}

func (h *MarketplaceHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r.URL.Path)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	var input dto.MarketplaceDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}

	result, err := h.service.Update(r.Context(), id, input)
	if err != nil {
		if err.Error() == "marketplace not found" {
			writeError(w, http.StatusNotFound, err.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *MarketplaceHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r.URL.Path)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	result, err := h.service.FindByID(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if result == nil {
		writeError(w, http.StatusNotFound, "marketplace not found")
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *MarketplaceHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.FindAll(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *MarketplaceHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r.URL.Path)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	err = h.service.Delete(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
