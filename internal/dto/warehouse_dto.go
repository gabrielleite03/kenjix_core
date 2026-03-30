package dto

import model "github.com/gabrielleite03/kenjix_domain/model"

type WarehouseDTO struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"address"`
	Capacity *int64 `json:"capacity,omitempty"`
	Active   bool   `json:"active"`
}

type WarehousePlaceTypeDTO struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Value  string `json:"value"`
	Active bool   `json:"active"`
}

type WarehousePlaceDTO struct {
	ID                   int64  `json:"id"`
	Name                 string `json:"name"`
	Active               bool   `json:"active"`
	WarehousePlaceTypeID *int64 `json:"warehouse_place_type_id,omitempty"`
	WarehouseID          *int64 `json:"warehouse_id,omitempty"`

	WarehousePlaceType *WarehousePlaceTypeDTO `json:"warehouse_place_type,omitempty"`
	Warehouse          *WarehouseDTO          `json:"warehouse,omitempty"`
}

func FromWarehouse(m *model.Warehouse) *WarehouseDTO {
	if m == nil {
		return nil
	}

	return &WarehouseDTO{
		ID:       m.ID,
		Name:     m.Name,
		Address:  m.Address,
		Capacity: m.Capacity,
		Active:   m.Active,
	}
}

func FromWarehousePlaceType(m *model.WarehousePlaceType) *WarehousePlaceTypeDTO {
	if m == nil {
		return nil
	}

	return &WarehousePlaceTypeDTO{
		ID:     m.ID,
		Name:   m.Name,
		Value:  m.Value,
		Active: m.Active,
	}
}

func FromWarehousePlace(m *model.WarehousePlace) *WarehousePlaceDTO {
	if m == nil {
		return nil
	}

	return &WarehousePlaceDTO{
		ID:                   m.ID,
		Name:                 m.Name,
		Active:               m.Active,
		WarehousePlaceTypeID: m.WarehousePlaceTypeID,
		WarehouseID:          m.WarehouseID,
	}
}
