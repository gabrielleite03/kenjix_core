package service

import (
	"errors"
	"fmt"

	"github.com/gabrielleite03/kenjix_core/internal/dto"
	"github.com/gabrielleite03/kenjix_domain/model"
	"github.com/gabrielleite03/kenjix_persist/repository"
	"github.com/shopspring/decimal"
)

type PurchaseService struct {
	dao              repository.PurchaseDAO
	supplierDAO      repository.SupplierDAO
	productDAO       repository.ProductDAO
	warehouseDAO     repository.WarehouseDAO
	stockDao         repository.StockDAO
	stockMovementDAO repository.StockMovementDAO
}

func NewPurchaseService() *PurchaseService {
	return &PurchaseService{
		dao:              *repository.NewPurchaseDAO(),
		supplierDAO:      repository.NewSupplierDAO(),
		productDAO:       repository.NewProductDAO(),
		warehouseDAO:     repository.NewWarehouseDAO(),
		stockDao:         *repository.NewStockDAO(),
		stockMovementDAO: *repository.NewStockMovementDAO(),
	}
}

func (s *PurchaseService) Create(input dto.PurchaseCreateDTO) (*dto.PurchaseResponseDTO, error) {

	total := calculateTotal(input.Items)

	purchase := &model.Purchase{
		InvoiceNumber: input.InvoiceNumber,
		InvoiceType:   input.InvoiceType,
		SupplierID:    input.SupplierID,
		Status:        input.Status,
		Total:         decimal.NewFromFloat(total),
	}

	items := make([]model.PurchaseItem, len(input.Items))

	for i, item := range input.Items {
		itemTotal := item.Quantity.Mul(item.CostPrice)

		items[i] = model.PurchaseItem{
			ProductID:    item.ProductID,
			Quantity:     item.Quantity,
			CostPrice:    item.CostPrice,
			Total:        itemTotal,
			CostCenterID: item.CostCenterID,
		}
	}

	purchase.Items = items

	err := s.dao.Create(purchase)
	if err != nil {
		return nil, errors.New("Falha ao criar purchase")
	}

	purchase, _ = s.dao.FindByID(purchase.ID)

	if input.Status == "Recebido" {
		err := s.atualizarStock(purchase.Items, &purchase.ID)
		if err != nil {
			purchase.Status = "Pendente"
			s.dao.Update(purchase)
			return nil, errors.Join(
				fmt.Errorf("falha ao atualizar stock"),
				err,
			)
		}
	}

	return s.toResponse(purchase)
}

func (s *PurchaseService) Update(input dto.PurchaseUpdateDTO) (*dto.PurchaseResponseDTO, error) {

	total := calculateTotal(input.Items)

	purchase := &model.Purchase{
		ID:            input.ID,
		InvoiceNumber: input.InvoiceNumber,
		InvoiceType:   input.InvoiceType,
		SupplierID:    input.SupplierID,
		Status:        input.Status,
		Total:         decimal.NewFromFloat(total),
	}

	items := make([]model.PurchaseItem, len(input.Items))

	for i, item := range input.Items {
		itemTotal := item.Quantity.Mul(item.CostPrice)

		items[i] = model.PurchaseItem{
			ProductID:    item.ProductID,
			Quantity:     item.Quantity,
			CostPrice:    item.CostPrice,
			Total:        itemTotal,
			CostCenterID: item.CostCenterID,
		}
	}

	purchase.Items = items

	// verificar se é para atualizar ou não o stock

	persistedPurchase, errP := s.dao.FindByID(input.ID)
	if errP != nil {
		return nil, errors.New("Falha ao recuperar purchase da base de dados")
	}
	if persistedPurchase.Status == "Recebido" {
		return nil, errors.New("Esta purchase já foi utilizada para atualizar o stock e não pode ser alterada")
	}
	err := s.dao.Update(purchase)
	if err != nil {
		return nil, errors.New("Falha ao atualizar purchase")
	}
	persistedPurchase, _ = s.dao.FindByID(input.ID)

	if input.Status == "Recebido" {
		err := s.atualizarStock(persistedPurchase.Items, &input.ID)
		if err != nil {
			return nil, err
		}
	}

	return s.toResponse(purchase)
}

func (s *PurchaseService) atualizarStock(items []model.PurchaseItem, purchaseID *int64) error {
	// carrega todos os produtos na memoria
	allProducts, _ := s.productDAO.List()

	// carrega o stock na memoria
	allStock, _ := s.stockDao.GetGroupedByProductAndWarehouse()

	// carregar todos warehouse_place na memoria
	allWarehousePlace, _ := s.warehouseDAO.FindAllWarehousePlace()

	// localizar todos warehouse_place que tem free space
	warehousePLaceWithFreeSpace := s.getWarehousePlaceWithFreeSpace(allWarehousePlace, allStock, allProducts)

	totalFreeSpace := decimal.Zero

	for _, freeSpace := range warehousePLaceWithFreeSpace {
		totalFreeSpace = totalFreeSpace.Add(freeSpace)
	}

	// validar se se tem espaço livre antes
	var estimatedSpace decimal.Decimal
	estimatedSpace = decimal.NewFromFloat(0.0)
	for _, item := range items {
		// buscar o produto
		produto := s.findProductByID(item.ProductID, allProducts)
		if produto == nil {
			continue
		}
		volumeTotalDoItem := produto.Volume.Mul(item.Quantity)
		estimatedSpace = estimatedSpace.Add(volumeTotalDoItem)
	}

	if estimatedSpace.GreaterThan(totalFreeSpace) {
		return errors.New("Não Há espaço suficiente para armazenar os itens")
	}

	// Verificar se já existe registro em Stock para o produto.
	for _, item := range items {
		// buscar o produto
		produto := s.findProductByID(item.ProductID, allProducts)
		volume := produto.Volume
		volumeTotalDoItem := volume.Mul(item.Quantity)

		//localizar place com freeSapce que ja tenha o mesmo produto
		placeID := s.getWarehousePlaceIdWithFreeSpaceByProductId(
			item.ProductID, allStock, warehousePLaceWithFreeSpace, volumeTotalDoItem)
		if placeID == nil {
			placeID = s.getWarehousePlaceIdWithFreeSpace(warehousePLaceWithFreeSpace, volumeTotalDoItem)
		}
		if placeID == nil {
			return errors.New("Não foi encontrado place com capacidade suficiente disponível")
		}

		// se não, inserir um novo stok
		stock := &model.Stock{
			Product:        *s.findProductByID(produto.ID, allProducts),
			WarehousePlace: *s.findWarehousePLaceByID(*placeID, allWarehousePlace),
			PurchaseItem:   item,
			Quantity:       int(item.Quantity.IntPart()),
			Active:         true,
		}
		err := s.stockDao.Create(stock)
		if err != nil {
			return errors.New("Falha ao criar Stock")
		}

		referenceType := "PURCHASE"
		stockMovement := &model.StockMovement{
			ProductID:        produto.ID,
			WarehousePlaceID: *placeID,
			PurchseItemID:    item.ID,
			Quantity:         int(item.Quantity.IntPart()),
			ReferenceID:      purchaseID,
			Type:             model.StockMovementIn,
			ReferenceType:    &referenceType,
			Reason:           "Stock Adjustment",
		}

		err = s.stockMovementDAO.Create(stockMovement)
		if err != nil {
			return errors.New("Falha ao criar Movement")
		}

	}
	return nil
}

func (s *PurchaseService) FindByID(id int64) (*dto.PurchaseResponseDTO, error) {
	purchase, err := s.dao.FindByID(id)
	if err != nil {
		return nil, err
	}

	return s.toResponse(purchase)
}

func (s *PurchaseService) FindAll() ([]dto.PurchaseResponseDTO, error) {
	list, err := s.dao.FindAll()
	if err != nil {
		return nil, err
	}

	response := make([]dto.PurchaseResponseDTO, len(list))

	for i, p := range list {
		resp, _ := s.toResponse(&p)
		response[i] = *resp
	}

	return response, nil
}

func calculateTotal(items []dto.PurchaseItemCreateDTO) float64 {
	total := decimal.NewFromFloat(0.0)

	for _, item := range items {
		total = total.Add(item.Quantity.Mul(item.CostPrice))
	}

	return total.InexactFloat64()
}

func (s *PurchaseService) toResponse(m *model.Purchase) (*dto.PurchaseResponseDTO, error) {

	supplier, _ := s.supplierDAO.FindByID(m.SupplierID)

	items := make([]dto.PurchaseItemDTO, len(m.Items))

	for i, item := range m.Items {

		product, _ := s.productDAO.GetByID(item.ProductID)

		items[i] = dto.PurchaseItemDTO{
			ProductID:    item.ProductID,
			ProductName:  product.Name,
			Quantity:     item.Quantity,
			CostPrice:    item.CostPrice,
			Total:        item.Total,
			CostCenterID: item.CostCenterID,
		}
	}

	return &dto.PurchaseResponseDTO{
		ID:            m.ID,
		InvoiceNumber: m.InvoiceNumber,
		InvoiceType:   m.InvoiceType,
		SupplierID:    m.SupplierID,
		SupplierName:  supplier.NomeFantasia,
		Items:         items,
		Total:         m.Total,
		Status:        m.Status,
		CreatedAt:     m.CreatedAt,
	}, nil
}

func (s *PurchaseService) findProductByID(productId int64, allProducts []model.Product) *model.Product {

	for _, product := range allProducts {
		if product.ID == productId {
			return &product
		}
	}
	return nil
}

func (s *PurchaseService) findWarehousePLaceByID(placeID int64, allWarehousePlace []*model.WarehousePlace) *model.WarehousePlace {

	for _, place := range allWarehousePlace {
		if place.ID == placeID {
			return place
		}
	}
	return nil
}

func (s *PurchaseService) getWarehousePlaceWithFreeSpace(
	allWarehousePlace []*model.WarehousePlace,
	allStock []model.Stock,
	allProducts []model.Product) map[int64]decimal.Decimal {
	placesWithFreeSpace := make(map[int64]decimal.Decimal)
	for _, place := range allWarehousePlace {
		var placeUsedSpace decimal.Decimal
		for _, stock := range allStock {
			if stock.WarehousePlace.ID == place.ID && stock.Active {
				// buscar o produto
				product := s.findProductByID(stock.Product.ID, allProducts)
				volume := product.Volume.Mul(decimal.NewFromInt(int64(stock.Quantity)))
				placeUsedSpace = placeUsedSpace.Add(volume)
			}
		}
		freeSpace := decimal.NewFromInt(*place.Capacity).Sub(placeUsedSpace)
		if freeSpace.GreaterThan(decimal.Zero) {
			placesWithFreeSpace[place.ID] = freeSpace
		}
	}
	return placesWithFreeSpace
}

func (s *PurchaseService) getWarehousePlaceIdWithFreeSpaceByProductId(
	productId int64,
	allStock []model.Stock,
	allWarehousePlaceFreeSpace map[int64]decimal.Decimal,
	estimatedSpace decimal.Decimal) *int64 {
	for placeID, freeSpace := range allWarehousePlaceFreeSpace {
		if len(allStock) == 0 && freeSpace.GreaterThan(estimatedSpace) {
			return &placeID
		}
		for _, stock := range allStock {
			if stock.Product.ID == productId && stock.WarehousePlace.ID == placeID && freeSpace.GreaterThanOrEqual(estimatedSpace) {
				return &placeID
			}
		}
	}
	return nil
}

func (s *PurchaseService) getWarehousePlaceIdWithFreeSpace(
	allWarehousePlaceFreeSpace map[int64]decimal.Decimal,
	estimatedSpace decimal.Decimal) *int64 {
	for placeID, freeSpace := range allWarehousePlaceFreeSpace {
		if freeSpace.GreaterThanOrEqual(estimatedSpace) {
			return &placeID
		}
	}
	return nil
}

type Volume struct {
	Item        model.PurchaseItem `json:"items,omitempty"`
	VolumeTotal decimal.Decimal    `json:"volume_total" db:"volume_total"`
}
