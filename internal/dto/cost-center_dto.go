package dto

import (
	"strings"

	model "github.com/gabrielleite03/kenjix_domain/model"
	"github.com/shopspring/decimal"
)

type CostCenterCreateDTO struct {
	Name        string                        `json:"name"`
	Code        string                        `json:"code"`
	Description string                        `json:"description"`
	Properties  []CostCenterPropertyCreateDTO `json:"properties,omitempty"`
}

type CostCenterUpdateDTO struct {
	ID          int64                         `json:"id"`
	Name        string                        `json:"name"`
	Code        string                        `json:"code"`
	Description string                        `json:"description"`
	Active      bool                          `json:"active"`
	Properties  []CostCenterPropertyCreateDTO `json:"properties,omitempty"`
}

type CostCenterResponseDTO struct {
	ID          int64                   `json:"id"`
	Name        string                  `json:"name"`
	Code        string                  `json:"code"`
	Description string                  `json:"description"`
	Active      bool                    `json:"active"`
	Properties  []CostCenterPropertyDTO `json:"properties,omitempty"`
}

type CostCenterPropertyCreateDTO struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

type CostCenterPropertyDTO struct {
	ID           int64           `json:"id"`
	CostCenterID int64           `json:"costCenterId"`
	Name         string          `json:"name"`
	Value        decimal.Decimal `json:"value"`
	Type         string          `json:"type"`
}

func ToCostCenterModelCreate(d CostCenterCreateDTO) *model.CostCenter {
	properties := make([]model.CostCenterProperty, len(d.Properties))

	for i, p := range d.Properties {
		dec, _ := decimal.NewFromString(strings.ReplaceAll(p.Value, ",", "."))
		properties[i] = model.CostCenterProperty{
			Name:  p.Name,
			Value: dec,
			Type:  model.CostCenterPropertyType(p.Type),
		}
	}

	return &model.CostCenter{
		Name:        d.Name,
		Code:        d.Code,
		Description: d.Description,
		Properties:  properties,
	}
}

func ToCostCenterResponse(m *model.CostCenter) CostCenterResponseDTO {
	properties := make([]CostCenterPropertyDTO, len(m.Properties))

	for i, p := range m.Properties {
		properties[i] = CostCenterPropertyDTO{
			ID:           p.ID,
			CostCenterID: m.ID,
			Name:         p.Name,
			Value:        p.Value,
			Type:         string(p.Type),
		}
	}

	return CostCenterResponseDTO{
		ID:          m.ID,
		Name:        m.Name,
		Code:        m.Code,
		Description: m.Description,
		Active:      m.Active,
		Properties:  properties,
	}
}
