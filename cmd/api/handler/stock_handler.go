package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gabrielleite03/kenjix_core/internal/dto"
	"github.com/gabrielleite03/kenjix_core/internal/service"
	"github.com/gabrielleite03/kenjix_persist/repository"
)

type StockHandler struct {
	service *service.StockService
}

// NewStockHandler cria uma instância do handler
func NewStockHandler() *StockHandler {
	repo := repository.NewStockDAO() // Ajuste para seu repositório real
	productDAO := repository.NewProductDAO()
	warehouseDAO := repository.NewWarehouseDAO()
	costCenterDAO := repository.NewCostCenterDAO()
	return &StockHandler{
		service: service.NewStockService(*repo, productDAO, warehouseDAO, *costCenterDAO),
	}
}

// Create cria um novo stock
func (h *StockHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input dto.StockDTO

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "JSON inválido: "+err.Error(), http.StatusBadRequest)
		return
	}

	result, err := h.service.Create(&input)
	if err != nil {
		http.Error(w, "Erro ao criar stock: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, result)
}

// FindAll retorna todos os stocks
func (h *StockHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.List()
	if err != nil {
		http.Error(w, "Erro ao buscar stocks: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// FindAll retorna todos os stocks
func (h *StockHandler) FindAllStockMovementsEager(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.FindAllStockMovementsEager()
	if err != nil {
		http.Error(w, "Erro ao buscar stocks: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, result)
}
