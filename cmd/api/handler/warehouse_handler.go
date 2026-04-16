package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/gabrielleite03/kenjix_core/internal/dto"
	"github.com/gabrielleite03/kenjix_core/internal/service"
)

type WarehouseHandler struct {
	service service.WarehouseService
}

func NewWarehouseHandler(s service.WarehouseService) *WarehouseHandler {
	return &WarehouseHandler{service: s}
}

type WarehousePlaceHandler struct {
	service service.WarehouseService
}

func NewWarehousePlaceHandler(s service.WarehouseService) *WarehousePlaceHandler {
	return &WarehousePlaceHandler{service: s}
}

type WarehouseLocationsHandler struct {
	service service.WarehouseService
}

func NewWarehouseLocationsHandler(s service.WarehouseService) *WarehouseLocationsHandler {
	return &WarehouseLocationsHandler{service: s}
}

// POST /warehouses
func (h *WarehouseHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dto dto.WarehouseDTO

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	created, err := h.service.Create(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

// PUT /warehouses
func (h *WarehouseHandler) Update(w http.ResponseWriter, r *http.Request) {
	var dto dto.WarehouseDTO

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updated, err := h.service.Update(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, updated)
}

// GET /warehouses
func (h *WarehouseHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.FindAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, list)
}

// GET /warehouses
func (h *WarehousePlaceHandler) FindAllPlaceType(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.FindAllPlaceTypes()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, list)
}

// GET /warehouses/{id}
func (h *WarehouseHandler) FindByID(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	result, err := h.service.FindByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if result == nil {
		http.NotFound(w, r)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// DELETE /warehouses/{id}
func (h *WarehouseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// POST /warehouses
func (h *WarehouseLocationsHandler) CreateWarehouseLocation(w http.ResponseWriter, r *http.Request, warehouseID int64) {

	var dto dto.WarehousePlaceDTO

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dto.WarehouseID = &warehouseID

	created, err := h.service.CreateWarehouseLocation(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

func (h *WarehouseLocationsHandler) UpdateWarehousePlace(w http.ResponseWriter, r *http.Request) {

	var dto dto.WarehousePlaceDTO

	body, _ := io.ReadAll(r.Body)

	if err := json.Unmarshal(body, &dto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updated, err := h.service.UpdateWarehousePlace(&dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, updated)
}

func (h *WarehouseLocationsHandler) DeleteWarehouseLocation(w http.ResponseWriter, r *http.Request, warehousePlaceID int64) {

	err := h.service.DeleteWarehousePlace(warehousePlaceID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, nil)
}

// POST /warehouses
func (h *WarehouseLocationsHandler) FindAllLocationByWarehouseID(w http.ResponseWriter, r *http.Request, warehouseID int64) {

	places, err := h.service.FindAllLocationsByWarehouseID(warehouseID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, places)
}

// ---------------- helpers ----------------

func getIDFromPath(path string) (int64, error) {
	parts := strings.Split(path, "/")
	idStr := parts[len(parts)-1]
	return strconv.ParseInt(idStr, 10, 64)
}
