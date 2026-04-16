package service

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gabrielleite03/kenjix_core/internal/dto"
	model "github.com/gabrielleite03/kenjix_domain/model"
	"github.com/gabrielleite03/kenjix_persist/repository"
	persist "github.com/gabrielleite03/kenjix_persist/repository"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ProductService interface {
	CreateProduct(ctx context.Context, prod *dto.ProductDTO) (int64, error)
	GetProduct(ctx context.Context, id int64) (*dto.ProductDTO, error)
	UpdateProduct(ctx context.Context, prod *dto.ProductDTO) error
	DeleteProduct(ctx context.Context, id int64) error
	ListProducts(ctx context.Context) ([]model.Product, error)
	ListProductsByMarketplace(ctx context.Context, marketplaceParam string) ([]dto.ProductHomeDTO, error)
	GetProductByMarketplace(ctx context.Context, id int64, marketplaceParam string) (*dto.ProductHomeDTO, error)
}

// ProductService provides product-related operations
type productServiceImpl struct {
	repo           persist.ProductDAO
	marketplaceDao persist.MarketplaceDAO
	stockDAO       persist.StockDAO
	categoryDAO    persist.CategoryDAO
	costCenterDAO  *repository.CostCenterDAO
	purchaseDAO    persist.PurchaseDAO
}

// NewProductService creates a new ProductService
func NewProductService(repo persist.ProductDAO) ProductService {
	return &productServiceImpl{
		repo:           repo,
		marketplaceDao: persist.NewMarketplaceDAO(),
		stockDAO:       *persist.NewStockDAO(),
		categoryDAO:    persist.NewCategoryDAO(),
		costCenterDAO:  persist.NewCostCenterDAO(),
		purchaseDAO:    *persist.NewPurchaseDAO(),
	}
}

// CreateProduct creates a new product
func (s *productServiceImpl) CreateProduct(ctx context.Context, prod *dto.ProductDTO) (int64, error) {

	// imagens
	e := saveImageToS3(prod)
	if e != nil {
		return 0, e
	}

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

func (s *productServiceImpl) GetProductByMarketplace(ctx context.Context, id int64, marketplaceParam string) (*dto.ProductHomeDTO, error) {
	marketplaces, _ := s.marketplaceDao.FindAll(ctx)

	var marketplace *model.Marketplace
	for i := range marketplaces {
		if strings.EqualFold(marketplaces[i].Name, marketplaceParam) {
			marketplace = &marketplaces[i]
			break
		}
	}
	if marketplace == nil {
		return nil, errors.New("Marketplace não localizado")
	}

	productModel, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	cat, err := s.categoryDAO.GetByID(*productModel.CategoryID)
	productModel.Category = cat

	// calcular price
	allStocks, _ := s.stockDAO.GetByProduct(productModel.ID)
	s.fillWithPurchaseItem(allStocks)

	// buscar os centros de custo
	allCostCenters, _ := s.costCenterDAO.FindAll()
	var prices []decimal.Decimal
	for _, stock := range allStocks {
		if stock.IsActive() {
			stockWithPurchaseItem := s.findStockPurchaseItem(stock, allStocks)

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
		}
	}

	p := dto.FromProduct(productModel)

	productHomeDTO := dto.ToProductHomeDTO(p)

	stocksWithQuantity, _ := s.stockDAO.GetGroupedByProduct()
	s.putAvailble(&productHomeDTO, stocksWithQuantity)

	productHomeDTO.Price = *s.getMaxPrice(prices)

	maxPrice := *s.getMaxPrice(prices)

	if marketplace != nil {

		rate := marketplace.CommissionRate.Div(decimal.NewFromInt(100))
		multiplier := decimal.NewFromInt(1).Add(rate)

		productHomeDTO.Price = maxPrice.Mul(multiplier)

	} else {
		productHomeDTO.Price = maxPrice
	}

	return &productHomeDTO, nil

}

func (s *productServiceImpl) getMaxPrice(prices []decimal.Decimal) *decimal.Decimal {
	if len(prices) == 0 {
		zero := decimal.Zero
		return &zero
	}
	max := prices[0]

	for _, p := range prices {
		if p.GreaterThan(max) {
			max = p
		}
	}

	return &max
}

func (s *productServiceImpl) fillWithPurchaseItem(stocks []model.Stock) {
	for i := range stocks {
		pi, _ := s.purchaseDAO.GetPurchaseItemByID(stocks[i].PurchaseItem.ID)
		if pi != nil {
			stocks[i].PurchaseItem = *pi
		}
	}
}
func (s *productServiceImpl) findStockPurchaseItem(
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

func (s *productServiceImpl) findCostCenter(
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

// UpdateProduct updates an existing product
func (s *productServiceImpl) UpdateProduct(ctx context.Context, prod *dto.ProductDTO) error {
	prodDb, err := s.repo.GetByID(prod.ID)
	if err != nil {
		return err
	}

	// deletar imagens antigas do s3
	for _, img := range prod.DeletedImages {
		err = deleteFromS3(img)
		if err != nil {
			return err
		}
	}

	// salvar novas imagens no s3
	err = saveImageToS3(prod)
	if err != nil {
		return err
	}

	// remover da base de dados as imagens que foram deletadas

	prodDb.Name = prod.Name
	prodDb.SKU = prod.SKU
	prodDb.Price = prod.Price
	prodDb.Marca = prod.Brand
	prodDb.Description = prod.Description
	prodDb.Active = prod.Active
	prodDb.Volume = prod.Volume
	prodDb.CategoryID = prod.CategoryID
	prodDb.Properties = prod.ToModel().Properties

	prodDb.Images = filterDeletedImages(prod.Images, prodDb.Images, prod.DeletedImages)

	return s.repo.Update(prodDb)
}

func filterDeletedImages(
	savedImages []dto.ProductImageDTO,
	images []model.ProductImage,
	deleted []string,
) []model.ProductImage {

	deletedMap := make(map[string]struct{})

	for _, url := range deleted {
		deletedMap[url] = struct{}{}
	}

	var result []model.ProductImage

	// 🔥 adiciona savedImages ao result
	for _, img := range savedImages {
		if _, found := deletedMap[img.URL]; !found {
			result = append(result, model.ProductImage{
				URL: img.URL,
			})
		}
	}

	// mantém images do banco
	for _, img := range images {
		if _, found := deletedMap[img.URL]; !found {
			result = append(result, img)
		}
	}

	return result
}

// DeleteProduct deletes a product by ID
func (s *productServiceImpl) DeleteProduct(ctx context.Context, id int64) error {
	return s.repo.Delete(id)
}

// ListProducts lists all products
func (s *productServiceImpl) ListProducts(ctx context.Context) ([]model.Product, error) {
	return s.repo.List()
}

func (s *productServiceImpl) ListProductsByMarketplace(ctx context.Context, marketplaceParam string) ([]dto.ProductHomeDTO, error) {
	marketplaces, _ := s.marketplaceDao.FindAll(ctx)

	var marketplace *model.Marketplace
	for i := range marketplaces {
		if strings.EqualFold(marketplaces[i].Name, marketplaceParam) {
			marketplace = &marketplaces[i]
			break
		}
	}
	if marketplace == nil {
		return nil, errors.New("Marketplace não localizado")
	}

	// recuperar stocks

	stocksWithQuantity, _ := s.stockDAO.GetGroupedByProduct()

	products, _ := s.repo.List()

	// buscar os centros de custo
	allCostCenters, _ := s.costCenterDAO.FindAll()

	productHomeDTO := make([]dto.ProductHomeDTO, len(products))
	for i, p := range products {
		productHomeDTO[i] = dto.ToProductHomeDTO(dto.FromProduct(&p))
		s.putAvailble(&productHomeDTO[i], stocksWithQuantity)

		allStocks, _ := s.stockDAO.GetByProduct(p.ID)
		s.fillWithPurchaseItem(allStocks)

		var prices []decimal.Decimal
		for _, stock := range allStocks {
			if stock.IsActive() {
				stockWithPurchaseItem := s.findStockPurchaseItem(stock, allStocks)

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
			}
		}

		productHomeDTO[i].Price = *s.getMaxPrice(prices)
		maxPrice := *s.getMaxPrice(prices)
		if marketplace != nil {

			rate := marketplace.CommissionRate.Div(decimal.NewFromInt(100))
			multiplier := decimal.NewFromInt(1).Add(rate)

			productHomeDTO[i].Price = maxPrice.Mul(multiplier)

		} else {
			productHomeDTO[i].Price = maxPrice
		}

	}
	return productHomeDTO, nil
}

func (s *productServiceImpl) putAvailble(
	productHomeDTO *dto.ProductHomeDTO,
	stocksWithQuantity []model.Stock,
) {

	for _, stock := range stocksWithQuantity {

		if stock.Product.ID == productHomeDTO.ID {
			productHomeDTO.Available = stock.Quantity > 0
			return
		}
	}

	// se não encontrou, garante false
	productHomeDTO.Available = false
}

func saveImageToS3(prod *dto.ProductDTO) error {
	const s3Path = "uploads/product/images/"

	imageList := []dto.ProductImageDTO{}
	files := prod.ImageFiles
	for pos, fileHeader := range files {

		file, err := fileHeader.Open()
		if err != nil {
			return err
		}

		ext := filepath.Ext(fileHeader.Filename)
		fileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)

		err = uploadToS3(file, fileName, s3Path)
		file.Close()

		if err != nil {
			return err
		}

		// adicionar na lista
		imageList = append(imageList, dto.ProductImageDTO{
			URL:      "https://aws-s3-site-kejipet.s3.us-east-1.amazonaws.com/" + s3Path + fileName,
			Position: pos,
		})
	}
	prod.Images = imageList
	return nil
}

func uploadToS3(file multipart.File, fileName string, urlS3 string) error {
	ctx := context.TODO()

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		return err
	}

	client := s3.NewFromConfig(cfg)

	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String("aws-s3-site-kejipet"),
		Key:    aws.String(urlS3 + fileName),
		Body:   file,
	})

	return err
}

func deleteFromS3(fileName string) error {
	ctx := context.TODO()

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		return err
	}

	client := s3.NewFromConfig(cfg)

	_, err = client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String("aws-s3-site-kejipet"),
		Key:    aws.String(fileName),
	})

	if err != nil {
		return err
	}

	return err
}
