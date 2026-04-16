package service

import (
	"github.com/gabrielleite03/kenjix_core/internal/dto"
	"github.com/gabrielleite03/kenjix_domain/model"
	persist "github.com/gabrielleite03/kenjix_persist/repository"
)

type SupplierService interface {
	GetAll() ([]dto.SupplierDTO, error)
	GetByID(id int64) (*dto.SupplierDTO, error)
	Create(supplier dto.SupplierDTO) (*dto.SupplierDTO, error)
	Update(supplier dto.SupplierDTO) error
	Delete(id int64) error
}

type supplierService struct {
	dao persist.SupplierDAO
}

func NewSupplierService(dao persist.SupplierDAO) SupplierService {
	return &supplierService{dao: dao}
}

func (s *supplierService) GetAll() ([]dto.SupplierDTO, error) {
	suppliers, err := s.dao.FindAll()
	if err != nil {
		return nil, err
	}

	var result []dto.SupplierDTO
	for _, sup := range suppliers {
		result = append(result, toSupplierDTO(sup))
	}

	return result, nil
}

func (s *supplierService) GetByID(id int64) (*dto.SupplierDTO, error) {
	supplier, err := s.dao.FindByID(id)
	if err != nil {
		return nil, err
	}

	dto := toSupplierDTO(*supplier)
	return &dto, nil
}

func (s *supplierService) Create(supplier dto.SupplierDTO) (*dto.SupplierDTO, error) {
	model := toSupplierModel(supplier)

	err := s.dao.Create(&model)
	if err != nil {
		return nil, err
	}

	result := toSupplierDTO(model)
	return &result, nil
}

func (s *supplierService) Update(supplier dto.SupplierDTO) error {
	model := toSupplierModel(supplier)
	return s.dao.Update(&model)
}

func (s *supplierService) Delete(id int64) error {
	return s.dao.Delete(id)
}

func toSupplierModel(d dto.SupplierDTO) model.Supplier {
	return model.Supplier{
		ID:           d.ID,
		RazaoSocial:  d.RazaoSocial,
		NomeFantasia: d.NomeFantasia,
		CNPJ:         d.CNPJ,
		IE:           d.IE,
		Address:      d.Address,
		Salesperson:  d.Salesperson,
		Email:        d.Email,
		Phone:        d.Phone,
		Active:       d.Active,
		CategoryID:   d.CategoryID,
	}
}

func toSupplierDTO(m model.Supplier) dto.SupplierDTO {
	return dto.SupplierDTO{
		ID:           m.ID,
		RazaoSocial:  m.RazaoSocial,
		NomeFantasia: m.NomeFantasia,
		CNPJ:         m.CNPJ,
		IE:           m.IE,
		Address:      m.Address,
		Salesperson:  m.Salesperson,
		Email:        m.Email,
		Phone:        m.Phone,
		Active:       m.Active,
		CategoryID:   m.CategoryID,
		Category:     (*dto.CategoryDTO)(m.Category),
	}
}
