package dto

import (
	"time"

	"github.com/shopspring/decimal"
)

type MarketplaceStatus string
type IntegrationType string

const (
	MarketplaceActive   MarketplaceStatus = "active"
	MarketplaceInactive MarketplaceStatus = "inactive"

	IntegrationAPI     IntegrationType = "api"
	IntegrationManual  IntegrationType = "manual"
	IntegrationWebhook IntegrationType = "webhook"
)

type MarketplaceDTO struct {
	ID              int64             `json:"id" db:"id"`
	Name            string            `json:"name"`
	Logo            *string           `json:"logo,omitempty"`
	Status          MarketplaceStatus `json:"status"`
	CommissionRate  decimal.Decimal   `json:"commissionRate"`
	IntegrationType IntegrationType   `json:"integrationType"`
	APIURL          *string           `json:"apiUrl,omitempty"`
	APIKey          *string           `json:"apiKey,omitempty"`
	APISecret       *string           `json:"apiSecret,omitempty"`
	APIEndpoint     *string           `json:"apiEndpoint,omitempty"`
	CreatedAt       time.Time         `json:"createdAt"`
}
