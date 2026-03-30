package dto

import model "github.com/gabrielleite03/kenjix_domain/model"

type CategoryDTO struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
}

func ToCategoryDTO(c *model.Category) *CategoryDTO {
	if c == nil {
		return nil
	}

	return &CategoryDTO{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		Active:      c.Active,
	}
}

func (d *CategoryDTO) ToModel() *model.Category {
	if d == nil {
		return nil
	}

	return &model.Category{
		ID:          d.ID,
		Name:        d.Name,
		Description: d.Description,
		Active:      d.Active,
	}
}
