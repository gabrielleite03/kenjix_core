package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

type PurchaseCreateDTO struct {
	InvoiceNumber *string                 `json:"invoiceNumber,omitempty"`
	InvoiceType   string                  `json:"invoiceType"`
	SupplierID    int64                   `json:"supplierId,string"`
	Status        string                  `json:"status"`
	Items         []PurchaseItemCreateDTO `json:"items"`
}

type PurchaseUpdateDTO struct {
	ID            int64                   `json:"id"`
	InvoiceNumber *string                 `json:"invoiceNumber,omitempty"`
	InvoiceType   string                  `json:"invoiceType"`
	SupplierID    int64                   `json:"supplierId"`
	Status        string                  `json:"status"`
	Items         []PurchaseItemCreateDTO `json:"items"`
}

type PurchaseResponseDTO struct {
	ID            int64             `json:"id"`
	InvoiceNumber *string           `json:"invoiceNumber,omitempty"`
	InvoiceType   string            `json:"invoiceType"`
	SupplierID    int64             `json:"supplierId"`
	SupplierName  string            `json:"supplierName"`
	Items         []PurchaseItemDTO `json:"items"`
	Total         decimal.Decimal   `json:"total"`
	Status        string            `json:"status"`
	CreatedAt     time.Time         `json:"createdAt"`
}

type PurchaseItemCreateDTO struct {
	ProductID    int64           `json:"productId"`
	Quantity     decimal.Decimal `json:"quantity"`
	CostPrice    decimal.Decimal `json:"costPrice"`
	CostCenterID *int64          `json:"costCenterId,omitempty"`
}

type PurchaseItemDTO struct {
	ProductID    int64           `json:"productId"`
	ProductName  string          `json:"productName"`
	Quantity     decimal.Decimal `json:"quantity"`
	CostPrice    decimal.Decimal `json:"costPrice"`
	Total        decimal.Decimal `json:"total"`
	CostCenterID *int64          `json:"costCenterId,omitempty"`
}
