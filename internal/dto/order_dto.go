package dto

import (
	"time"

	model "github.com/gabrielleite03/kenjix_domain/model"
	"github.com/shopspring/decimal"
)

type PurchaseStatusDTO struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Active      bool    `json:"active"`
}

type FiscalNumberTypeDTO struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Active      bool    `json:"active"`
}

type PurchaseOrderDTO struct {
	ID                 int64     `json:"id"`
	PurchaseStatusID   int64     `json:"purchase_status_id"`
	FiscalNumberTypeID int64     `json:"fiscal_number_type_id"`
	SupplierID         int64     `json:"supplier_id"`
	FiscalNumber       *string   `json:"fiscal_number,omitempty"`
	Status             string    `json:"status"`
	CreatedAt          time.Time `json:"created_at"`
	Active             bool      `json:"active"`

	PurchaseStatus   *PurchaseStatusDTO   `json:"purchase_status,omitempty"`
	FiscalNumberType *FiscalNumberTypeDTO `json:"fiscal_number_type,omitempty"`
	Supplier         *SupplierDTO         `json:"supplier,omitempty"`
}

type PurchaseOrderItemDTO struct {
	PurchaseOrderID int64           `json:"purchase_order_id"`
	ProductID       int64           `json:"product_id"`
	Quantity        int             `json:"quantity"`
	UnitPrice       decimal.Decimal `json:"unit_price"`
	Active          bool            `json:"active"`

	PurchaseOrder *PurchaseOrderDTO `json:"purchase_order,omitempty"`
	Product       *ProductDTO       `json:"product,omitempty"`
}

func FromPurchaseStatus(m *model.PurchaseStatus) *PurchaseStatusDTO {
	if m == nil {
		return nil
	}

	return &PurchaseStatusDTO{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		Active:      m.Active,
	}
}

func FromFiscalNumberType(m *model.FiscalNumberType) *FiscalNumberTypeDTO {
	if m == nil {
		return nil
	}

	return &FiscalNumberTypeDTO{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		Active:      m.Active,
	}
}

func FromPurchaseOrder(m *model.PurchaseOrder) *PurchaseOrderDTO {
	if m == nil {
		return nil
	}

	return &PurchaseOrderDTO{
		ID:                 m.ID,
		PurchaseStatusID:   m.PurchaseStatusID,
		FiscalNumberTypeID: m.FiscalNumberTypeID,
		SupplierID:         m.SupplierID,
		FiscalNumber:       m.FiscalNumber,
		Status:             m.Status,
		CreatedAt:          m.CreatedAt,
		Active:             m.Active,
	}
}

func FromPurchaseOrderItem(m *model.PurchaseOrderItem) *PurchaseOrderItemDTO {
	if m == nil {
		return nil
	}

	return &PurchaseOrderItemDTO{
		PurchaseOrderID: m.PurchaseOrderID,
		ProductID:       m.ProductID,
		Quantity:        m.Quantity,
		UnitPrice:       m.UnitPrice,
		Active:          m.Active,
	}
}
