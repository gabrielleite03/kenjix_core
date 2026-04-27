package dto

import (
	"time"

	model "github.com/gabrielleite03/kenjix_domain/model"
	"github.com/shopspring/decimal"
)

type StockDTO struct {
	ID               int64     `json:"id"`
	ProductID        int64     `json:"productId"`
	WarehousePlaceID int64     `json:"warehousePlaceId"`
	PurchaseItemID   int64     `json:"purchaseItemId"`
	Quantity         int       `json:"quantity"`
	Active           bool      `json:"active"`
	UpdatedAt        time.Time `json:"updatedAt"`
	// opcionais para tela
	Product            ProductHomeDTO  `json:"product"`
	ProductName        *string         `json:"productName,omitempty"`
	WarehouseName      *string         `json:"warehouseName,omitempty"`
	WarehousePlaceName *string         `json:"warehousePlaceName,omitempty"`
	MinPrice           decimal.Decimal `json:"minPrice,omitempty"`
	MaxPrice           decimal.Decimal `json:"maxPrice,omitempty"`
}

func (d *StockDTO) ToStockDTO(s *model.Stock) StockDTO {
	return StockDTO{
		ID:               s.ID,
		ProductID:        s.Product.ID,
		WarehousePlaceID: s.WarehousePlace.ID,
		PurchaseItemID:   s.PurchaseItem.ID,
		Quantity:         s.Quantity,
		Active:           s.Active,
		UpdatedAt:        s.UpdatedAt,
	}
}

func (d *StockDTO) ToStockModel() model.Stock {
	return model.Stock{
		ID: d.ID,
		Product: model.Product{
			ID: d.ProductID,
		},
		WarehousePlace: model.WarehousePlace{
			ID: d.WarehousePlaceID,
		},
		PurchaseItem: model.PurchaseItem{
			ID: d.PurchaseItemID,
		},
		Quantity:  d.Quantity,
		Active:    d.Active,
		UpdatedAt: d.UpdatedAt,
	}
}

type StockMovementDTO struct {
	ID               int64  `json:"id"`
	ProductID        int64  `json:"productId"`
	WarehousePlaceID int64  `json:"warehousePlaceId"`
	PurchaseItemID   int64  `json:"purchaseItemId"`
	Type             string `json:"type"`
	Quantity         int    `json:"quantity"`

	ReferenceID   *int64  `json:"referenceId,omitempty"`
	ReferenceType *string `json:"referenceType,omitempty"`

	Reason    string    `json:"reason,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}

func (s *StockMovementDTO) ToStockMovementDTO(m *model.StockMovement) StockMovementDTO {
	return StockMovementDTO{
		ID:               m.ID,
		ProductID:        m.ProductID,
		WarehousePlaceID: m.WarehousePlaceID,
		PurchaseItemID:   m.PurchseItemID,
		Type:             string(m.Type),
		Quantity:         m.Quantity,
		ReferenceID:      m.ReferenceID,
		ReferenceType:    m.ReferenceType,
		Reason:           m.Reason,
		CreatedAt:        m.CreatedAt,
	}
}

func (d StockMovementDTO) ToStockMovementModel() model.StockMovement {
	return model.StockMovement{
		ID:               d.ID,
		ProductID:        d.ProductID,
		WarehousePlaceID: d.WarehousePlaceID,
		PurchseItemID:    d.PurchaseItemID,
		Type:             model.StockMovementType(d.Type),
		Quantity:         d.Quantity,
		ReferenceID:      d.ReferenceID,
		ReferenceType:    d.ReferenceType,
		Reason:           d.Reason,
		CreatedAt:        d.CreatedAt,
	}
}

type StockMovementEagerDTO struct {
	ID int64 `json:"id"`

	Product        ProductDTO        `json:"product"`
	WarehousePlace WarehousePlaceDTO `json:"warehouse_place"`
	PurchaseItem   PurchaseItemDTO   `json:"purchase_item"`

	Type     string `json:"type"`
	Quantity int    `json:"quantity"`

	ReferenceID   *int64  `json:"reference_id,omitempty"`
	ReferenceType *string `json:"reference_type,omitempty"`

	Reason *string `json:"reason,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}

func (d *StockMovementEagerDTO) ToStockMovementDTO(m model.StockMovementEager) StockMovementEagerDTO {
	return StockMovementEagerDTO{
		ID: m.ID,

		Product: ProductDTO{
			ID:    m.Product.ID,
			Name:  m.Product.Name,
			SKU:   m.Product.SKU,
			Brand: m.Product.Marca,
		},

		WarehousePlace: WarehousePlaceDTO{
			ID:   m.WarehousePlace.ID,
			Name: m.WarehousePlace.Name,
		},

		PurchaseItem: PurchaseItemDTO{
			ProductID:    m.PurchaseItem.ProductID,
			ProductName:  m.Product.Name,
			Quantity:     m.PurchaseItem.Quantity,
			CostPrice:    m.PurchaseItem.CostPrice,
			CostCenterID: m.PurchaseItem.CostCenterID,
			Total:        m.PurchaseItem.Total,
		},

		Type:     string(m.Type),
		Quantity: m.Quantity,

		ReferenceID:   m.ReferenceID,
		ReferenceType: m.ReferenceType,
		Reason:        m.Reason,

		CreatedAt: m.CreatedAt,
	}
}
