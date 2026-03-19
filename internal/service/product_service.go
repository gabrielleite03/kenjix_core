package service

import (
	"context"

	"github.com/gabrielleite03/kenjix_core/internal/dto"
	"github.com/gabrielleite03/kenjix_core/internal/wrapper"
	model "github.com/gabrielleite03/kenjix_domain/model"
	persist "github.com/gabrielleite03/kenjix_persist/repository"
)

type ProductService interface {
	CreateProduct(ctx context.Context, prod *dto.ProductDTO) (int64, error)
	GetProduct(ctx context.Context, id int64) (*dto.ProductDTO, error)
	UpdateProduct(ctx context.Context, prod *model.Product) (int64, error)
	DeleteProduct(ctx context.Context, id int64) (int64, error)
	ListProducts(ctx context.Context) ([]*model.Product, error)
}

// ProductService provides product-related operations
type productServiceImpl struct {
	repo persist.ProductRepository
}

// NewProductService creates a new ProductService
func NewProductService(repo persist.ProductRepository) ProductService {
	return &productServiceImpl{
		repo: repo,
	}
}

// CreateProduct creates a new product
func (s *productServiceImpl) CreateProduct(ctx context.Context, prod *dto.ProductDTO) (int64, error) {
	productModel := prod.ToProductModel()
	prodId, err := s.repo.Create(ctx, productModel)
	s.repo.CreateProductProperties(ctx, prodId, prod.Properties)
	return prodId, err
}

// GetProduct retrieves a product by ID
func (s *productServiceImpl) GetProduct(ctx context.Context, id int64) (*dto.ProductDTO, error) {
	productModel, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return wrapper.NewProductDTOFromModel(productModel), nil

}

// UpdateProduct updates an existing product
func (s *productServiceImpl) UpdateProduct(ctx context.Context, prod *model.Product) (int64, error) {
	return s.repo.Update(ctx, prod)
}

// DeleteProduct deletes a product by ID
func (s *productServiceImpl) DeleteProduct(ctx context.Context, id int64) (int64, error) {
	return s.repo.Delete(ctx, id)
}

// ListProducts lists all products
func (s *productServiceImpl) ListProducts(ctx context.Context) ([]*model.Product, error) {
	return s.repo.List(ctx)
}
