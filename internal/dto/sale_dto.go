package dto

import (
	"time"

	model "github.com/gabrielleite03/kenjix_domain/model"
	"github.com/shopspring/decimal"
)

type PaymentMethodDTO struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

type SalesOrderDTO struct {
	ID              int64           `json:"id"`
	Price           decimal.Decimal `json:"price"`
	Discount        decimal.Decimal `json:"discount"`
	Status          string          `json:"status"`
	PaymentMethodID int64           `json:"payment_method_id"`
	Active          bool            `json:"active"`
	CreatedAt       time.Time       `json:"created_at"`

	PaymentMethod *PaymentMethodDTO `json:"payment_method,omitempty"`
}

type SalesOrderItemDTO struct {
	SalesOrderID int64           `json:"sales_order_id"`
	ProductID    int64           `json:"product_id"`
	Quantity     int             `json:"quantity"`
	UnitPrice    decimal.Decimal `json:"unit_price"`

	SalesOrder *SalesOrderDTO `json:"sales_order,omitempty"`
	Product    *ProductDTO    `json:"product,omitempty"`
}

func FromPaymentMethod(m *model.PaymentMethod) *PaymentMethodDTO {
	if m == nil {
		return nil
	}

	return &PaymentMethodDTO{
		ID:     m.ID,
		Name:   m.Name,
		Active: m.Active,
	}
}

func FromSalesOrder(m *model.SalesOrder) *SalesOrderDTO {
	if m == nil {
		return nil
	}

	return &SalesOrderDTO{
		ID:              m.ID,
		Price:           m.Price,
		Discount:        m.Discount,
		Status:          m.Status,
		PaymentMethodID: m.PaymentMethodID,
		Active:          m.Active,
		CreatedAt:       m.CreatedAt,
	}
}

func FromSalesOrderItem(m *model.SalesOrderItem) *SalesOrderItemDTO {
	if m == nil {
		return nil
	}

	return &SalesOrderItemDTO{
		SalesOrderID: m.SalesOrderID,
		ProductID:    m.ProductID,
		Quantity:     m.Quantity,
		UnitPrice:    m.UnitPrice,
	}
}
