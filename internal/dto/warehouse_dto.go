package dto

import "github.com/gabrielleite03/kenjix_domain/model"

// WarehouseDTO representa a transferência de dados de Warehouse
type WarehouseDTO struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Address  string `json:"location"`
	Capacity *int64 `json:"capacity,omitempty"`
	Active   bool   `json:"active"`
}

// WarehousePlaceTypeDTO representa a transferência de dados de WarehousePlaceType
type WarehousePlaceTypeDTO struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Value  string `json:"value"`
	Active bool   `json:"active"`
}

// WarehousePlaceDTO representa a transferência de dados de WarehousePlace
type WarehousePlaceDTO struct {
	ID                   int64  `json:"id,string"`
	Capacity             int64  `json:"capacity"`
	Name                 string `json:"name"`
	Active               bool   `json:"active"`
	WarehousePlaceTypeID *int64 `json:"warehouse_place_type_id,omitempty"`
	WarehouseID          *int64 `json:"warehouseId,string"`
	Type                 *int64 `json:"type,omitempty"`

	WarehousePlaceType *WarehousePlaceTypeDTO `json:"warehouse_place_type,omitempty"`
	Warehouse          *WarehouseDTO          `json:"warehouse,omitempty"`
}

// -------------------- Funções de conversão --------------------

func (w *WarehouseDTO) ToWarehouseModel() *model.Warehouse {
	if w == nil {
		return nil
	}
	return &model.Warehouse{
		ID:       w.ID,
		Name:     w.Name,
		Address:  w.Address,
		Capacity: w.Capacity,
		Active:   w.Active,
	}
}

func (wpt *WarehousePlaceTypeDTO) ToWarehousePlaceTypeModel() *model.WarehousePlaceType {
	if wpt == nil {
		return nil
	}
	return &model.WarehousePlaceType{
		ID:     wpt.ID,
		Name:   wpt.Name,
		Value:  wpt.Value,
		Active: wpt.Active,
	}
}

func (wp *WarehousePlaceDTO) ToWarehousePlaceModel() *model.WarehousePlace {
	if wp == nil {
		return nil
	}
	return &model.WarehousePlace{
		ID:                   wp.ID,
		Name:                 wp.Name,
		Active:               wp.Active,
		Capacity:             &wp.Capacity,
		WarehousePlaceTypeID: coalesceInt64(wp.WarehousePlaceTypeID, wp.Type),
		WarehouseID:          wp.WarehouseID,
	}
}

func coalesceInt64(a, b *int64) *int64 {
	if a != nil {
		return a
	}
	return b
}

func (w *WarehouseDTO) ToWarehouseModelWithRelations() *model.Warehouse {
	if w == nil {
		return nil
	}
	return &model.Warehouse{
		ID:       w.ID,
		Name:     w.Name,
		Address:  w.Address,
		Capacity: w.Capacity,
		Active:   w.Active,
	}
}

func (wpt *WarehousePlaceTypeDTO) ToWarehousePlaceTypeModelWithRelations() *model.WarehousePlaceType {
	if wpt == nil {
		return nil
	}
	return &model.WarehousePlaceType{
		ID:     wpt.ID,
		Name:   wpt.Name,
		Value:  wpt.Value,
		Active: wpt.Active,
	}
}

// FromWarehouse converte model.Warehouse para WarehouseDTO
func FromWarehouse(w *model.Warehouse) *WarehouseDTO {
	if w == nil {
		return nil
	}
	return &WarehouseDTO{
		ID:       w.ID,
		Name:     w.Name,
		Address:  w.Address,
		Capacity: w.Capacity,
		Active:   w.Active,
	}
}

// FromWarehousePlaceType converte model.WarehousePlaceType para WarehousePlaceTypeDTO
func FromWarehousePlaceType(wpt *model.WarehousePlaceType) *WarehousePlaceTypeDTO {
	if wpt == nil {
		return nil
	}
	return &WarehousePlaceTypeDTO{
		ID:     wpt.ID,
		Name:   wpt.Name,
		Value:  wpt.Value,
		Active: wpt.Active,
	}
}

// FromWarehousePlace converte model.WarehousePlace para WarehousePlaceDTO
func FromWarehousePlace(wp *model.WarehousePlace) *WarehousePlaceDTO {
	if wp == nil {
		return nil
	}
	return &WarehousePlaceDTO{
		ID:                   wp.ID,
		Name:                 wp.Name,
		Active:               wp.Active,
		Capacity:             *wp.Capacity,
		WarehousePlaceTypeID: wp.WarehousePlaceTypeID,
		WarehouseID:          wp.WarehouseID,
	}
}
