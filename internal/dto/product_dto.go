package dto

import (
	model "github.com/gabrielleite03/kenjix_domain/model"
	"github.com/shopspring/decimal"
)

type ProductDTO struct {
	ID          int64           `json:"id"`
	Name        string          `json:"name"`
	SKU         string          `json:"sku"`
	Price       decimal.Decimal `json:"price"`
	Marca       string          `json:"marca"`
	Description string          `json:"description"`
	Active      bool            `json:"active"`
	CategoryID  *int64          `json:"category_id,omitempty"`

	Properties []ProductPropertyDTO `json:"properties,omitempty"`
	Images     []ProductImageDTO    `json:"images,omitempty"`
	Videos     []ProductVideoDTO    `json:"videos,omitempty"`
}

type ProductPropertyDTO struct {
	ID        int64  `json:"id"`
	ProductID int64  `json:"product_id"`
	Name      string `json:"name"`
	Value     string `json:"value"`
}

type ProductImageDTO struct {
	ID        int64  `json:"id"`
	ProductID int64  `json:"product_id"`
	URL       string `json:"url"`
	Position  int    `json:"position"`
	IsPrimary bool   `json:"is_primary"`
}

type ProductVideoDTO struct {
	ID        int64   `json:"id"`
	ProductID int64   `json:"product_id"`
	URL       string  `json:"url"`
	Provider  *string `json:"provider,omitempty"`
}

func FromProduct(m *model.Product) *ProductDTO {
	if m == nil {
		return nil
	}

	dto := &ProductDTO{
		ID:          m.ID,
		Name:        m.Name,
		SKU:         m.SKU,
		Price:       m.Price,
		Marca:       m.Marca,
		Description: m.Description,
		Active:      m.Active,
		CategoryID:  m.CategoryID,
	}

	for _, p := range m.Properties {
		dto.Properties = append(dto.Properties, ProductPropertyDTO{
			ID:        p.ID,
			ProductID: p.ProductID,
			Name:      p.Name,
			Value:     p.Value,
		})
	}

	for _, i := range m.Images {
		dto.Images = append(dto.Images, ProductImageDTO{
			ID:        i.ID,
			ProductID: i.ProductID,
			URL:       i.URL,
			Position:  i.Position,
			IsPrimary: i.IsPrimary,
		})
	}

	for _, v := range m.Videos {
		dto.Videos = append(dto.Videos, ProductVideoDTO{
			ID:        v.ID,
			ProductID: v.ProductID,
			URL:       v.URL,
			Provider:  v.Provider,
		})
	}

	return dto
}

func (d *ProductDTO) ToModel() *model.Product {
	if d == nil {
		return nil
	}

	m := &model.Product{
		ID:          d.ID,
		Name:        d.Name,
		SKU:         d.SKU,
		Price:       d.Price,
		Marca:       d.Marca,
		Description: d.Description,
		Active:      d.Active,
		CategoryID:  d.CategoryID,
	}

	for _, p := range d.Properties {
		m.Properties = append(m.Properties, model.ProductProperty{
			ID:        p.ID,
			ProductID: p.ProductID,
			Name:      p.Name,
			Value:     p.Value,
		})
	}

	for _, i := range d.Images {
		m.Images = append(m.Images, model.ProductImage{
			ID:        i.ID,
			ProductID: i.ProductID,
			URL:       i.URL,
			Position:  i.Position,
			IsPrimary: i.IsPrimary,
		})
	}

	for _, v := range d.Videos {
		m.Videos = append(m.Videos, model.ProductVideo{
			ID:        v.ID,
			ProductID: v.ProductID,
			URL:       v.URL,
			Provider:  v.Provider,
		})
	}

	return m
}
