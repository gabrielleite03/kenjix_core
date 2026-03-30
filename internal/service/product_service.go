package service

import (
	"context"

	"github.com/gabrielleite03/kenjix_core/internal/dto"
	model "github.com/gabrielleite03/kenjix_domain/model"
	persist "github.com/gabrielleite03/kenjix_persist/repository"
)

type ProductService interface {
	CreateProduct(ctx context.Context, prod *dto.ProductDTO) (int64, error)
	GetProduct(ctx context.Context, id int64) (*dto.ProductDTO, error)
	UpdateProduct(ctx context.Context, prod *model.Product) error
	DeleteProduct(ctx context.Context, id int64) error
	ListProducts(ctx context.Context) ([]model.Product, error)
}

// ProductService provides product-related operations
type productServiceImpl struct {
	repo persist.ProductDAO
}

// NewProductService creates a new ProductService
func NewProductService(repo persist.ProductDAO) ProductService {
	return &productServiceImpl{
		repo: repo,
	}
}

// CreateProduct creates a new product
func (s *productServiceImpl) CreateProduct(ctx context.Context, prod *dto.ProductDTO) (int64, error) {
	productModel := prod.ToModel()
	err := s.repo.Create(productModel)
	if err != nil {
		return 0, err
	}
	return productModel.ID, err
}

// GetProduct retrieves a product by ID
func (s *productServiceImpl) GetProduct(ctx context.Context, id int64) (*dto.ProductDTO, error) {
	productModel, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return dto.FromProduct(productModel), nil

}

// UpdateProduct updates an existing product
func (s *productServiceImpl) UpdateProduct(ctx context.Context, prod *model.Product) error {
	return s.repo.Update(prod)
}

// DeleteProduct deletes a product by ID
func (s *productServiceImpl) DeleteProduct(ctx context.Context, id int64) error {
	return s.repo.Delete(id)
}

// ListProducts lists all products
func (s *productServiceImpl) ListProducts(ctx context.Context) ([]model.Product, error) {
	return s.repo.List()
}
