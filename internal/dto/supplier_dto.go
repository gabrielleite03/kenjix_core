package dto

import model "github.com/gabrielleite03/kenjix_domain/model"

type SupplierDTO struct {
	ID         int64   `json:"id"`
	Name       string  `json:"name"`
	CNPJ       string  `json:"cnpj"`
	Country    *string `json:"country,omitempty"`
	Email      *string `json:"email,omitempty"`
	Phone      *string `json:"phone,omitempty"`
	Seller     *string `json:"seller,omitempty"`
	SellerFone *string `json:"seller_fone,omitempty"`
	Active     bool    `json:"active"`
	CategoryID *int64  `json:"category_id,omitempty"`

	Category *CategoryDTO `json:"category,omitempty"`
}

func FromSupplier(m *model.Supplier) *SupplierDTO {
	if m == nil {
		return nil
	}

	return &SupplierDTO{
		ID:         m.ID,
		Name:       m.Name,
		CNPJ:       m.CNPJ,
		Country:    m.Country,
		Email:      m.Email,
		Phone:      m.Phone,
		Seller:     m.Seller,
		SellerFone: m.SellerFone,
		Active:     m.Active,
		CategoryID: m.CategoryID,
	}
}
