package dto

import (
	"time"

	model "github.com/gabrielleite03/kenjix_domain/model"
)

type StockMovementType string

const (
	StockMovementIn         StockMovementType = "IN"
	StockMovementOut        StockMovementType = "OUT"
	StockMovementAdjustment StockMovementType = "ADJUSTMENT"
)

type StockDTO struct {
	ProductID   int64     `json:"product_id"`
	WarehouseID int64     `json:"warehouse_id"`
	Quantity    int       `json:"quantity"`
	Active      bool      `json:"active"`
	UpdatedAt   time.Time `json:"updated_at,omitempty"`
}

type StockMovementDTO struct {
	ID          int64             `json:"id"`
	ProductID   int64             `json:"product_id"`
	WarehouseID int64             `json:"warehouse_id"`
	Type        StockMovementType `json:"type"`
	Quantity    int               `json:"quantity"`
	CreatedAt   time.Time         `json:"created_at"`
	Reason      string            `json:"reason"`
}

func FromStock(m *model.Stock) *StockDTO {
	if m == nil {
		return nil
	}

	return &StockDTO{
		ProductID:   m.ProductID,
		WarehouseID: m.WarehouseID,
		Quantity:    m.Quantity,
		Active:      m.Active,
		UpdatedAt:   m.UpdatedAt,
	}
}

func FromStockMovement(m *model.StockMovement) *StockMovementDTO {
	if m == nil {
		return nil
	}

	return &StockMovementDTO{
		ID:          m.ID,
		ProductID:   m.ProductID,
		WarehouseID: m.WarehouseID,
		Type:        StockMovementType(m.Type),
		Quantity:    m.Quantity,
		CreatedAt:   m.CreatedAt,
		Reason:      m.Reason,
	}
}
