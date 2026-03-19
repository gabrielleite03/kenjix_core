package dto

import model "github.com/gabrielleite03/kenjix_domain/model"

// ProductDTO represents the data transfer object for a product
type ProductDTO struct {
	ID         int64   `json:"id"`
	Name       string  `json:"name"`
	SKU        string  `json:"sku"`
	Price      float64 `json:"price"`
	Active     bool    `json:"active"`
	CategoryID *int64  `json:"category,omitempty"`

	Properties map[string]string `json:"properties,omitempty"`
	Images     []string          `json:"images,omitempty"`
	Videos     []string          `json:"videos,omitempty"`
}

func (s *ProductDTO) ToProductModel() *model.Product {
	if s == nil {
		return nil
	}
	return &model.Product{
		ID:         s.ID,
		Name:       s.Name,
		SKU:        s.SKU,
		Price:      s.Price,
		Active:     s.Active,
		CategoryID: s.CategoryID,
	}
}
