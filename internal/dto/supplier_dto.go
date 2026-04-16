package dto

import model "github.com/gabrielleite03/kenjix_domain/model"

type SupplierDTO struct {
	ID           int64   `json:"id"`
	RazaoSocial  string  `json:"razaoSocial"`
	NomeFantasia string  `json:"nomeFantasia"`
	CNPJ         string  `json:"cnpj"`
	IE           *string `json:"ie,omitempty"`
	Address      *string `json:"address,omitempty"`
	Salesperson  *string `json:"salesperson,omitempty"`
	Email        *string `json:"email,omitempty"`
	Phone        *string `json:"phone,omitempty"`
	Active       bool    `json:"active"`
	CategoryID   *int64  `json:"categoryId,omitempty"`

	Category *CategoryDTO `json:"category,omitempty"`
}

func FromSupplier(m *model.Supplier) *SupplierDTO {
	if m == nil {
		return nil
	}

	return &SupplierDTO{
		ID:           m.ID,
		RazaoSocial:  m.RazaoSocial,
		NomeFantasia: m.NomeFantasia,
		CNPJ:         m.CNPJ,
		IE:           m.IE,
		Address:      m.Address,
		Salesperson:  m.Salesperson,
		Email:        m.Email,
		Phone:        m.Phone,
		Active:       m.Active,
		CategoryID:   m.CategoryID,
		Category:     (*CategoryDTO)(m.Category),
	}
}
