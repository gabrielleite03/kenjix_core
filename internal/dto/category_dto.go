package dto

import model "github.com/gabrielleite03/kenjix_domain/model"

type CategoryDTO struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (c *CategoryDTO) ToCategoryModel() *model.Category {
	if c == nil {
		return nil
	}
	return &model.Category{
		ID:   c.ID,
		Name: c.Name,
	}
}
