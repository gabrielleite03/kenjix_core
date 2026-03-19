package wrapper

import (
	"github.com/gabrielleite03/kenjix_core/internal/dto"
	model "github.com/gabrielleite03/kenjix_domain/model"
)

func NewProductDTOFromModel(m *model.Product) *dto.ProductDTO {
	if m == nil {
		return nil
	}

	return &dto.ProductDTO{
		ID:         m.ID,
		Name:       m.Name,
		SKU:        m.SKU,
		Price:      m.Price,
		Active:     m.Active,
		CategoryID: m.CategoryID,

		Properties: map[string]string{},
		Images:     []string{},
		Videos:     []string{},
	}
}

func NewProductDTOListFromModel(models []*model.Product) []*dto.ProductDTO {
	if models == nil {
		return nil
	}

	dtos := make([]*dto.ProductDTO, 0, len(models))
	for _, m := range models {
		dtos = append(dtos, NewProductDTOFromModel(m))
	}
	return dtos
}

func NewProductModelFromDTO(p *dto.ProductDTO) *model.Product {
	if p == nil {
		return nil
	}

	return &model.Product{
		ID:         p.ID,
		Name:       p.Name,
		SKU:        p.SKU,
		Price:      p.Price,
		Active:     p.Active,
		CategoryID: p.CategoryID,
	}
}
