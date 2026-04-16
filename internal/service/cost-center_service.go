package service

import (
	"strings"

	"github.com/gabrielleite03/kenjix_core/internal/dto"
	"github.com/gabrielleite03/kenjix_domain/model"
	"github.com/gabrielleite03/kenjix_persist/repository"
	persist "github.com/gabrielleite03/kenjix_persist/repository"
	"github.com/shopspring/decimal"
)

type CostCenterService struct {
	dao *repository.CostCenterDAO
}

func NewCostCenterService(repo *persist.CostCenterDAO) *CostCenterService {
	return &CostCenterService{
		dao: repo,
	}
}

func (s *CostCenterService) Create(input dto.CostCenterCreateDTO) (*dto.CostCenterResponseDTO, error) {
	model := dto.ToCostCenterModelCreate(input)

	err := s.dao.Create(model)
	if err != nil {
		return nil, err
	}

	response := dto.ToCostCenterResponse(model)
	return &response, nil
}

func (s *CostCenterService) Update(input dto.CostCenterUpdateDTO) (*dto.CostCenterResponseDTO, error) {
	properties := make([]model.CostCenterProperty, len(input.Properties))

	for i, p := range input.Properties {
		dec, _ := decimal.NewFromString(strings.ReplaceAll(p.Value, ",", "."))
		properties[i] = model.CostCenterProperty{
			Name:  p.Name,
			Value: dec,
			Type:  model.CostCenterPropertyType(p.Type),
		}
	}

	cc := &model.CostCenter{
		ID:          input.ID,
		Name:        input.Name,
		Code:        input.Code,
		Description: input.Description,
		Active:      input.Active,
		Properties:  properties,
	}

	err := s.dao.Update(cc)
	if err != nil {
		return nil, err
	}

	response := dto.ToCostCenterResponse(cc)
	return &response, nil
}

func (s *CostCenterService) FindByID(id int64) (*dto.CostCenterResponseDTO, error) {
	cc, err := s.dao.FindByID(id)
	if err != nil {
		return nil, err
	}

	response := dto.ToCostCenterResponse(cc)
	return &response, nil
}

func (s *CostCenterService) FindAll() ([]dto.CostCenterResponseDTO, error) {
	list, err := s.dao.FindAll()
	if err != nil {
		return nil, err
	}

	response := make([]dto.CostCenterResponseDTO, len(list))

	for i, cc := range list {
		response[i] = dto.ToCostCenterResponse(&cc)
	}

	return response, nil
}

func (s *CostCenterService) Delete(id int64) error {
	return s.dao.Delete(id)
}
