package service

import (
	"time"

	"github.com/gabrielleite03/kenjix_core/internal/dto"
	"github.com/gabrielleite03/kenjix_domain/model"
	persist "github.com/gabrielleite03/kenjix_persist/repository"
	"github.com/shopspring/decimal"
)

type StockService struct {
	repo             persist.StockDAO
	productDAO       persist.ProductDAO
	warehouseDAO     persist.WarehouseDAO
	costCenterDAO    persist.CostCenterDAO
	purchaseDAO      persist.PurchaseDAO
	stockMovementDAO persist.StockMovementDAO
}

// NewStockService cria uma instância do service
func NewStockService(repo persist.StockDAO, productDAO persist.ProductDAO, warehouseDAO persist.WarehouseDAO, costCenterDAO persist.CostCenterDAO) *StockService {
	return &StockService{
		repo:             repo,
		productDAO:       productDAO,
		warehouseDAO:     warehouseDAO,
		costCenterDAO:    costCenterDAO,
		purchaseDAO:      *persist.NewPurchaseDAO(),
		stockMovementDAO: *persist.NewStockMovementDAO(),
	}
}

// Create cria um novo stock
func (s *StockService) Create(w *dto.StockDTO) (*dto.StockDTO, error) {
	model := w.ToStockModel()
	err := s.repo.Create(&model)

	if err != nil {
		return nil, err
	}
	return w, nil
}

// FindAll retorna todos os stocks
func (s *StockService) List() ([]*dto.StockDTO, error) {
	stocks, err := s.repo.GetGroupedByProductAndWarehouse()
	if err != nil {
		return nil, err
	}

	allCostCenters, _ := s.costCenterDAO.FindAll()

	allStocks, err := s.repo.GetAllActive()
	s.fillWithPurchaseItem(allStocks)

	allProducts, _ := s.productDAO.List()
	allWarehousePLace, _ := s.warehouseDAO.FindAllWarehousePlace()
	allWarehouse, _ := s.warehouseDAO.FindAll()

	var stockDTOs []*dto.StockDTO
	for _, stock := range stocks {
		stockdto := dto.ToStockDTO(&stock)
		product := s.findProductByID(stockdto.ProductID, allProducts)

		stockdto.ProductName = &product.Name
		productDTO := dto.FromProduct(product)
		stockdto.Product = dto.ToProductHomeDTO(productDTO)
		whp := s.findWarehousePLaceByID(stockdto.WarehousePlaceID, allWarehousePLace)
		wh := s.findWarehouseByID(*whp.WarehouseID, allWarehouse)
		stockdto.WarehousePlaceName = &whp.Name
		stockdto.WarehouseName = &wh.Name

		s.fillMinPriceAndMaxPrice(&stockdto, stock, allCostCenters, allStocks)

		stockDTOs = append(stockDTOs, &stockdto)
	}

	stockKits, _ := s.findAllWithKits(allStocks, allProducts, allCostCenters)
	for i := range stockKits {
		stockDTOs = append(stockDTOs, &stockKits[i])
	}

	return stockDTOs, nil
}

func (s *StockService) findAllWithKits(normalStocks []model.Stock, products []model.Product, allCostCenters []model.CostCenter) ([]dto.StockDTO, error) {

	result := BuildStockWithKits(dto.ToStockDTOList(normalStocks), dto.ToProductHomeDTOList(dto.FromProducts(products)), normalStocks, allCostCenters)
	return result, nil
}

func CalculateKitQuantity(componentStocks map[int64]int, kitItems []dto.ProductKitDTO) int {
	if len(kitItems) == 0 {
		return 0
	}

	minQty := -1

	for _, item := range kitItems {
		componentQty := componentStocks[item.ComponentProductID]
		possibleKits := componentQty / item.Quantity

		if minQty == -1 || possibleKits < minQty {
			minQty = possibleKits
		}
	}

	if minQty < 0 {
		return 0
	}

	return minQty
}

func BuildStockWithKits(
	normalStocks []dto.StockDTO,
	products []dto.ProductHomeDTO,
	allStocks []model.Stock,
	allCostCenters []model.CostCenter,
) []dto.StockDTO {

	result := make([]dto.StockDTO, 0, len(normalStocks))

	componentStocks := map[int64]int{}
	componentPriceMap := map[int64]dto.StockDTO{}
	existing := map[int64]bool{}

	for _, stock := range normalStocks {
		componentStocks[stock.ProductID] += stock.Quantity
		componentPriceMap[stock.ProductID] = stock
		existing[stock.ProductID] = true
	}

	for _, product := range products {
		if !product.IsKit {
			continue
		}

		if existing[product.ID] {
			continue
		}

		kitQty := CalculateKitQuantity(componentStocks, product.KitItemsDTOs)
		if kitQty <= 0 {
			continue
		}

		minPrice, maxPrice := CalculateKitPrices(product.KitItemsDTOs, allCostCenters, allStocks)

		name := product.Name

		result = append(result, dto.StockDTO{
			ID:        0,
			ProductID: product.ID,
			Quantity:  kitQty,
			Active:    product.Available && kitQty > 0,
			UpdatedAt: time.Time{},
			Product:   product,

			ProductName: &name,
			MinPrice:    minPrice,
			MaxPrice:    maxPrice,
		})
	}

	return result
}

func CalculateKitPrices(
	kitItems []dto.ProductKitDTO,
	allCostCenters []model.CostCenter,
	allStocks []model.Stock,
) (decimal.Decimal, decimal.Decimal) {

	minTotal := decimal.Zero
	maxTotal := decimal.Zero

	for _, item := range kitItems {
		componentMin, componentMax := calculateComponentPrice(
			item.ComponentProductID,
			allCostCenters,
			allStocks,
		)

		qty := decimal.NewFromInt(int64(item.Quantity))

		minTotal = minTotal.Add(componentMin.Mul(qty))

		// se max vier zerado, assume o min
		if componentMax.Equal(decimal.Zero) {
			componentMax = componentMin
		}

		maxTotal = maxTotal.Add(componentMax.Mul(qty))
	}

	if minTotal.Equal(maxTotal) {
		return minTotal, decimal.Zero
	}

	return minTotal, maxTotal
}

func calculateComponentPrice(
	productID int64,
	allCostCenters []model.CostCenter,
	allStocks []model.Stock,
) (decimal.Decimal, decimal.Decimal) {

	var prices []decimal.Decimal

	for _, stock := range allStocks {
		if stock.Product.ID != productID {
			continue
		}

		if stock.PurchaseItem.CostCenterID == nil {
			continue
		}
		cc := findCostCenter(*stock.PurchaseItem.CostCenterID, allCostCenters)
		if cc == nil {
			continue
		}

		basePrice := stock.PurchaseItem.CostPrice
		price := basePrice

		for _, p := range cc.Properties {
			switch p.Type {
			case "index":
				index := p.Value.Div(decimal.NewFromInt(100))
				price = price.Add(basePrice.Mul(index))

			case "value":
				price = price.Add(p.Value)
			}
		}

		prices = append(prices, price)
	}

	if len(prices) == 0 {
		return decimal.Zero, decimal.Zero
	}

	minPrice := prices[0]
	maxPrice := prices[0]

	for _, p := range prices {
		if p.LessThan(minPrice) {
			minPrice = p
		}
		if p.GreaterThan(maxPrice) {
			maxPrice = p
		}
	}

	if minPrice.Equal(maxPrice) {
		return minPrice, decimal.Zero
	}

	return minPrice, maxPrice
}

func findCostCenter(
	costCenterID int64,
	allCostCenter []model.CostCenter,
) *model.CostCenter {
	for _, c := range allCostCenter {
		if c.ID == costCenterID {
			return &c
		}
	}
	return nil
}

func (s *StockService) fillMinPriceAndMaxPrice(
	stockdto *dto.StockDTO,
	stock model.Stock,
	allCostCenters []model.CostCenter,
	allStocks []model.Stock,
) {

	stockWithPurchaseItem := s.findStockPurchaseItem(stock, allStocks)

	var prices []decimal.Decimal

	for _, v := range stockWithPurchaseItem {

		if v.PurchaseItem.CostCenterID == nil {
			continue
		}

		cc := s.findCostCenter(*v.PurchaseItem.CostCenterID, allCostCenters)
		if cc == nil {
			continue
		}

		basePrice := v.PurchaseItem.CostPrice
		price := basePrice
		for _, p := range cc.Properties {

			switch p.Type {

			case "index":
				index := p.Value.Div(decimal.NewFromInt(100))
				price = price.Add(basePrice.Mul(index))

			case "value":
				price = price.Add(p.Value)
			}

		}

		prices = append(prices, price)
	}

	if len(prices) == 0 {
		stockdto.MinPrice = decimal.Zero
		stockdto.MaxPrice = decimal.Zero
		return
	}

	minPrice := prices[0]
	maxPrice := prices[0]

	for _, p := range prices {
		if p.LessThan(minPrice) {
			minPrice = p
		}
		if p.GreaterThan(maxPrice) {
			maxPrice = p
		}
	}

	stockdto.MinPrice = minPrice

	if minPrice.Equal(maxPrice) {
		stockdto.MaxPrice = decimal.Zero
		return
	}

	stockdto.MaxPrice = maxPrice
}

func (s *StockService) findProductByID(productId int64, allProducts []model.Product) *model.Product {

	for _, product := range allProducts {
		if product.ID == productId {
			return &product
		}
	}
	return nil
}

func (s *StockService) findWarehousePLaceByID(placeID int64, allWarehousePlace []*model.WarehousePlace) *model.WarehousePlace {

	for _, place := range allWarehousePlace {
		if place.ID == placeID {
			return place
		}
	}
	return nil
}

func (s *StockService) findWarehouseByID(warehouseID int64, allWarehouse []*model.Warehouse) *model.Warehouse {

	for _, warehouse := range allWarehouse {
		if warehouse.ID == warehouseID {
			return warehouse
		}
	}
	return nil
}

func (s *StockService) fillWithPurchaseItem(stocks []model.Stock) {
	for i := range stocks {
		pi, _ := s.purchaseDAO.GetPurchaseItemByID(stocks[i].PurchaseItem.ID)
		if pi != nil {
			stocks[i].PurchaseItem = *pi
		}
	}
}

func (s *StockService) findStockPurchaseItem(
	stock model.Stock,
	stocks []model.Stock,
) []model.Stock {
	var result []model.Stock
	for i := range stocks {
		st := stocks[i]
		if stock.Product.ID == st.Product.ID &&
			stock.WarehousePlace.ID == st.WarehousePlace.ID {
			result = append(result, st)
		}
	}

	return result
}

func (s *StockService) findCostCenter(
	costCenterID int64,
	allCostCenter []model.CostCenter,
) *model.CostCenter {
	for _, c := range allCostCenter {
		if c.ID == costCenterID {
			return &c
		}
	}
	return nil
}

func (d *StockService) FindAllStockMovementsEager() ([]dto.StockMovementEagerDTO, error) {
	allStocks, err := d.stockMovementDAO.FindAllEager()
	if err != nil {
		return nil, err
	}

	stocks := make([]dto.StockMovementEagerDTO, 0, len(allStocks))

	for _, s := range allStocks {
		st := &dto.StockMovementEagerDTO{}
		stocks = append(stocks, st.ToStockMovementDTO(s))
	}

	return stocks, nil
}
