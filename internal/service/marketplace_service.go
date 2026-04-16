package service

import (
	"context"
	"errors"
	"time"

	"github.com/gabrielleite03/kenjix_core/internal/dto"
	"github.com/gabrielleite03/kenjix_domain/model"
	"github.com/gabrielleite03/kenjix_persist/repository"
)

type MarketplaceService interface {
	Create(ctx context.Context, input dto.MarketplaceDTO) (*dto.MarketplaceDTO, error)
	Update(ctx context.Context, id int64, input dto.MarketplaceDTO) (*dto.MarketplaceDTO, error)
	FindByID(ctx context.Context, id int64) (*dto.MarketplaceDTO, error)
	FindAll(ctx context.Context) ([]dto.MarketplaceDTO, error)
	Delete(ctx context.Context, id int64) error
}

type marketplaceService struct {
	repo repository.MarketplaceDAO
}

func NewMarketplaceService(repo repository.MarketplaceDAO) MarketplaceService {
	return &marketplaceService{repo: repo}
}

func (s *marketplaceService) Create(ctx context.Context, input dto.MarketplaceDTO) (*dto.MarketplaceDTO, error) {
	entity := model.Marketplace{
		Name:            input.Name,
		Logo:            input.Logo,
		Status:          string(input.Status),
		CommissionRate:  input.CommissionRate,
		IntegrationType: string(input.IntegrationType),
		APIURL:          input.APIURL,
		APIKey:          input.APIKey,
		APISecret:       input.APISecret,
		APIEndpoint:     input.APIEndpoint,
		CreatedAt:       time.Now(),
	}

	if err := s.repo.Create(ctx, &entity); err != nil {
		return nil, err
	}

	result := toDTO(entity)
	return &result, nil
}

func (s *marketplaceService) Update(ctx context.Context, id int64, input dto.MarketplaceDTO) (*dto.MarketplaceDTO, error) {
	entity, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, errors.New("marketplace not found")
	}

	entity.Name = input.Name
	entity.Logo = input.Logo
	entity.Status = string(input.Status)
	entity.CommissionRate = input.CommissionRate
	entity.IntegrationType = string(input.IntegrationType)
	entity.APIURL = input.APIURL
	entity.APIKey = input.APIKey
	entity.APISecret = input.APISecret
	entity.APIEndpoint = input.APIEndpoint

	if err := s.repo.Update(ctx, entity); err != nil {
		return nil, err
	}

	result := toDTO(*entity)
	return &result, nil
}

func (s *marketplaceService) FindByID(ctx context.Context, id int64) (*dto.MarketplaceDTO, error) {
	entity, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, nil
	}

	result := toDTO(*entity)
	return &result, nil
}

func (s *marketplaceService) FindAll(ctx context.Context) ([]dto.MarketplaceDTO, error) {
	entities, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	var result []dto.MarketplaceDTO
	for _, e := range entities {
		result = append(result, toDTO(e))
	}

	return result, nil
}

func (s *marketplaceService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func toDTO(m model.Marketplace) dto.MarketplaceDTO {
	return dto.MarketplaceDTO{
		ID:              m.ID,
		Name:            m.Name,
		Logo:            m.Logo,
		Status:          dto.MarketplaceStatus(m.Status),
		CommissionRate:  m.CommissionRate,
		IntegrationType: dto.IntegrationType(m.IntegrationType),
		APIURL:          m.APIURL,
		APIKey:          m.APIKey,
		APISecret:       m.APISecret,
		APIEndpoint:     m.APIEndpoint,
		CreatedAt:       m.CreatedAt,
	}
}
