package service

import (
	"context"

	"github.com/gabrielleite03/kenjix_core/internal/dto"
	persist "github.com/gabrielleite03/kenjix_persist/repository"
)

type CategoryService interface {
	List(ctx context.Context) ([]*dto.CategoryDTO, error)
	GetByID(ctx context.Context, id int64) (*dto.CategoryDTO, error)
	Create(ctx context.Context, d *dto.CategoryDTO) (*dto.CategoryDTO, error)
	Update(ctx context.Context, d *dto.CategoryDTO) (*dto.CategoryDTO, error)
	Delete(ctx context.Context, id int64) error
}

type categoryService struct {
	repo persist.CategoryDAO
}

func NewCategoryService(repo persist.CategoryDAO) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) List(ctx context.Context) ([]*dto.CategoryDTO, error) {
	categories, err := s.repo.List()
	if err != nil {
		return nil, err
	}

	result := make([]*dto.CategoryDTO, 0, len(categories))
	for _, c := range categories {
		result = append(result, dto.ToCategoryDTO(&c))
	}

	return result, nil
}

func (s *categoryService) GetByID(ctx context.Context, id int64) (*dto.CategoryDTO, error) {
	category, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return dto.ToCategoryDTO(category), nil
}

func (s *categoryService) Create(ctx context.Context, d *dto.CategoryDTO) (*dto.CategoryDTO, error) {
	model := d.ToModel()

	if err := s.repo.Create(model); err != nil {
		return nil, err
	}

	return dto.ToCategoryDTO(model), nil
}

func (s *categoryService) Update(ctx context.Context, d *dto.CategoryDTO) (*dto.CategoryDTO, error) {
	model := d.ToModel()

	if err := s.repo.Update(model); err != nil {
		return nil, err
	}

	return dto.ToCategoryDTO(model), nil
}

func (s *categoryService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(id)
}
