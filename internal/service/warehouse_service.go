package service

import (
	"errors"

	"github.com/gabrielleite03/kenjix_core/internal/dto"
	persist "github.com/gabrielleite03/kenjix_persist/repository"
)

// WarehouseRepo define os métodos que o repositório precisa implementar
type WarehouseService interface {
	Create(w *dto.WarehouseDTO) (*dto.WarehouseDTO, error)
	Update(w *dto.WarehouseDTO) (*dto.WarehouseDTO, error)
	FindByID(id int64) (*dto.WarehouseDTO, error)
	FindAll() ([]*dto.WarehouseDTO, error)
	Delete(id int64) error
	FindAllPlaceTypes() ([]*dto.WarehousePlaceTypeDTO, error)
	CreateWarehouseLocation(w *dto.WarehousePlaceDTO) (*dto.WarehousePlaceDTO, error)
	FindAllLocationsByWarehouseID(warehouseID int64) ([]*dto.WarehousePlaceDTO, error)
	DeleteWarehousePlace(id int64) error
	UpdateWarehousePlace(w *dto.WarehousePlaceDTO) (*dto.WarehousePlaceDTO, error)
}

// WarehouseService fornece a lógica de negócios para Warehouse
type warehouseService struct {
	repo persist.WarehouseDAO
}

// NewWarehouseService cria uma instância do service
func NewWarehouseService(repo persist.WarehouseDAO) WarehouseService {
	return &warehouseService{repo: repo}
}

// Create cria um novo warehouse
func (s *warehouseService) Create(w *dto.WarehouseDTO) (*dto.WarehouseDTO, error) {
	if w == nil {
		return nil, errors.New("warehouse DTO is nil")
	}

	modelWarehouse := w.ToWarehouseModel()

	created, err := s.repo.Create(modelWarehouse)
	if err != nil {
		return nil, err
	}

	return dto.FromWarehouse(created), nil
}

// Update atualiza um warehouse existente
func (s *warehouseService) Update(w *dto.WarehouseDTO) (*dto.WarehouseDTO, error) {
	if w == nil {
		return nil, errors.New("warehouse DTO is nil")
	}

	modelWarehouse := w.ToWarehouseModel()

	updated, err := s.repo.Update(modelWarehouse)
	if err != nil {
		return nil, err
	}

	return dto.FromWarehouse(updated), nil
}

// FindByID retorna um warehouse pelo ID
func (s *warehouseService) FindByID(id int64) (*dto.WarehouseDTO, error) {
	w, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return dto.FromWarehouse(w), nil
}

// FindAll retorna todos os warehouses
func (s *warehouseService) FindAll() ([]*dto.WarehouseDTO, error) {
	ws, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var dtos []*dto.WarehouseDTO
	for _, w := range ws {
		dtos = append(dtos, dto.FromWarehouse(w))
	}

	return dtos, nil
}

// FindAllPlaceTypes retorna todos os tipos de locais de armazenamento
func (s *warehouseService) FindAllPlaceTypes() ([]*dto.WarehousePlaceTypeDTO, error) {
	placeTypes, err := s.repo.FindAllWarehousePlaceType()
	if err != nil {
		return nil, err
	}

	var dtos []*dto.WarehousePlaceTypeDTO
	for _, pt := range placeTypes {
		dtos = append(dtos, dto.FromWarehousePlaceType(pt))
	}

	return dtos, nil
}

// Delete remove um warehouse pelo ID
func (s *warehouseService) Delete(id int64) error {
	return s.repo.Delete(id)
}

// CreateWarehouseLocation cria um novo local de armazenamento
func (s *warehouseService) CreateWarehouseLocation(w *dto.WarehousePlaceDTO) (*dto.WarehousePlaceDTO, error) {
	if w == nil {
		return nil, errors.New("warehouse place DTO is nil")
	}
	modelWarehousePlace := w.ToWarehousePlaceModel()

	created, err := s.repo.CreateWarehousePlace(modelWarehousePlace)
	if err != nil {
		return nil, err
	}
	return dto.FromWarehousePlace(created), nil
}

// FindAllLocationsByWarehouseID retorna todos os locais de armazenamento de um warehouse
func (s *warehouseService) FindAllLocationsByWarehouseID(warehouseID int64) ([]*dto.WarehousePlaceDTO, error) {
	places, err := s.repo.FindByWarehouseID(warehouseID)
	if err != nil {
		return nil, err
	}

	var result []*dto.WarehousePlaceDTO

	for _, p := range places {
		pDTO := dto.FromWarehousePlace(p)
		result = append(result, pDTO)
	}

	return result, nil
}

// DeleteWarehousePlace remove um local de armazenamento pelo ID
func (s *warehouseService) DeleteWarehousePlace(id int64) error {
	return s.repo.DeleteWarehousePlace(id)
}

// UpdateWarehousePlace atualiza um local de armazenamento existente
func (s *warehouseService) UpdateWarehousePlace(w *dto.WarehousePlaceDTO) (*dto.WarehousePlaceDTO, error) {
	if w == nil {
		return nil, errors.New("warehouse place DTO is nil")
	}
	modelWarehousePlace := w.ToWarehousePlaceModel()
	updated, err := s.repo.UpdateWarehousePlace(modelWarehousePlace)
	if err != nil {
		return nil, err
	}
	return dto.FromWarehousePlace(updated), nil
}
