package dto

import (
	"time"

	model "github.com/gabrielleite03/kenjix_domain/model"
	"github.com/shopspring/decimal"
)

type ImportProcessDTO struct {
	ID              int64           `json:"id"`
	PurchaseOrderID int64           `json:"purchase_order_id"`
	Incoterm        string          `json:"incoterm"`
	ExchangeRate    decimal.Decimal `json:"exchange_rate"`
	Status          string          `json:"status"`
	ArrivalDate     *time.Time      `json:"arrival_date,omitempty"`

	PurchaseOrder *PurchaseOrderDTO `json:"purchase_order,omitempty"`
}

type ImportCostDTO struct {
	ID              int64           `json:"id"`
	ImportProcessID int64           `json:"import_process_id"`
	Type            string          `json:"type"`
	Description     *string         `json:"description,omitempty"`
	Amount          decimal.Decimal `json:"amount"`
	Currency        string          `json:"currency"`

	ImportProcess *ImportProcessDTO `json:"import_process,omitempty"`
}

type ImportCostAllocationDTO struct {
	ImportCostID    int64           `json:"import_cost_id"`
	ProductID       int64           `json:"product_id"`
	AllocatedAmount decimal.Decimal `json:"allocated_amount"`

	ImportCost *ImportCostDTO `json:"import_cost,omitempty"`
	Product    *ProductDTO    `json:"product,omitempty"`
}

func FromImportProcess(m *model.ImportProcess) *ImportProcessDTO {
	if m == nil {
		return nil
	}

	return &ImportProcessDTO{
		ID:              m.ID,
		PurchaseOrderID: m.PurchaseOrderID,
		Incoterm:        m.Incoterm,
		ExchangeRate:    m.ExchangeRate,
		Status:          m.Status,
		ArrivalDate:     m.ArrivalDate,
	}
}

func FromImportCost(m *model.ImportCost) *ImportCostDTO {
	if m == nil {
		return nil
	}

	return &ImportCostDTO{
		ID:              m.ID,
		ImportProcessID: m.ImportProcessID,
		Type:            m.Type,
		Description:     m.Description,
		Amount:          m.Amount,
		Currency:        m.Currency,
	}
}

func FromImportCostAllocation(m *model.ImportCostAllocation) *ImportCostAllocationDTO {
	if m == nil {
		return nil
	}

	return &ImportCostAllocationDTO{
		ImportCostID:    m.ImportCostID,
		ProductID:       m.ProductID,
		AllocatedAmount: m.AllocatedAmount,
	}
}
