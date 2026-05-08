package dto

import (
	"mime/multipart"
	"time"

	model "github.com/gabrielleite03/kenjix_domain/model"
	"github.com/shopspring/decimal"
)

type ProductDTO struct {
	ID          int64            `json:"id"`
	Name        string           `json:"name"`
	SKU         string           `json:"sku"`
	Price       decimal.Decimal  `json:"price"`
	Brand       string           `json:"brand"`
	Weight      *decimal.Decimal `json:"weight,omitempty"`
	EAN         *string          `json:"ean"`
	NCM         *string          `json:"ncm"`
	Description string           `json:"description"`
	Active      bool             `json:"active"`
	Volume      decimal.Decimal  `json:"volume"`
	CategoryID  *int64           `json:"category_id,omitempty"`
	Category    *CategoryDTO     `json:"category,omitempty"`

	Prices []ProductPriceDTO `json:"prices,omitempty"` // preços por marketplace
	Stocks []StockDTO        `json:"stocks,omitempty"` // estoque por warehouse

	Properties             []ProductPropertyDTO    `json:"properties,omitempty"`
	Images                 []ProductImageDTO       `json:"images,omitempty"`
	ExistingImages         []string                `json:"existingImages,omitempty"`
	DeletedImages          []string                `json:"deletedImages,omitempty"`
	ImageFiles             []*multipart.FileHeader `json:"image_files,omitempty"`
	Videos                 []ProductVideoDTO       `json:"videos,omitempty"`
	ProductMarketplaceDTOs []ProductMarketplaceDTO `json:"product_marketplaces,omitempty"`
	IsKit                  bool                    `json:"is_kit"` // indica se é um kit (produto pai)
	KitItemsDTOs           []ProductKitDTO         `json:"kit_items,omitempty"`
}

type ProductMarketplaceDTO struct {
	ID            int64            `json:"id" db:"id"`
	ProductID     int64            `json:"productId"`     // ✅
	MarketplaceID int64            `json:"marketplaceId"` // ✅
	ExternalID    *string          `json:"externalId"`
	ProductURL    string           `json:"productUrl"`
	Price         *decimal.Decimal `json:"price"`
	ListingType   *string          `json:"listingType"`
	Status        *string          `json:"status"`
	Active        bool             `json:"active"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type ProductPropertyDTO struct {
	ID        int64  `json:"id"`
	ProductID int64  `json:"product_id"`
	Name      string `json:"name"`
	Value     string `json:"value"`
}

type ProductImageDTO struct {
	ID        int64  `json:"id"`
	ProductID int64  `json:"product_id"`
	URL       string `json:"url"`
	Position  int    `json:"position"`
	IsPrimary bool   `json:"is_primary"`
}

type ProductVideoDTO struct {
	ID        int64   `json:"id"`
	ProductID int64   `json:"product_id"`
	URL       string  `json:"url"`
	Provider  *string `json:"provider,omitempty"`
}

type ProductKitDTO struct {
	ID                 int64 `json:"id"`
	ProductID          int64 `json:"product_id"`
	ComponentProductID int64 `json:"productId"`
	Quantity           int   `json:"quantity"`
}

func ConvertToDTO(items []model.ProductKit) []ProductKitDTO {
	result := make([]ProductKitDTO, 0, len(items))

	for _, i := range items {
		result = append(result, ProductKitDTO{
			ID:                 i.ID,
			ProductID:          i.ProductID,
			ComponentProductID: i.ComponentProductID,
			Quantity:           i.Quantity,
		})
	}

	return result
}

func FromProduct(m *model.Product) *ProductDTO {
	if m == nil {
		return nil
	}

	dto := &ProductDTO{
		ID:           m.ID,
		Name:         m.Name,
		SKU:          m.SKU,
		Price:        m.Price,
		Brand:        m.Marca,
		Weight:       m.Weight,
		Description:  m.Description,
		Active:       m.Active,
		Volume:       m.Volume,
		CategoryID:   m.CategoryID,
		Category:     (*CategoryDTO)(m.Category),
		EAN:          m.EAN,
		NCM:          m.NCM,
		KitItemsDTOs: ConvertToDTO(m.KitComponents),
		IsKit:        m.IsKit,
	}

	for _, p := range m.Properties {
		dto.Properties = append(dto.Properties, ProductPropertyDTO{
			ID:        p.ID,
			ProductID: p.ProductID,
			Name:      p.Name,
			Value:     p.Value,
		})
	}

	for _, i := range m.Images {
		dto.Images = append(dto.Images, ProductImageDTO{
			ID:        i.ID,
			ProductID: i.ProductID,
			URL:       i.URL,
			Position:  i.Position,
			IsPrimary: i.IsPrimary,
		})
	}

	for _, v := range m.Videos {
		dto.Videos = append(dto.Videos, ProductVideoDTO{
			ID:        v.ID,
			ProductID: v.ProductID,
			URL:       v.URL,
			Provider:  v.Provider,
		})
	}

	for _, pm := range m.ProductMarketplaces {
		dto.ProductMarketplaceDTOs = append(dto.ProductMarketplaceDTOs, ProductMarketplaceDTO{
			ID:            pm.ID,
			ProductID:     pm.ProductID,
			MarketplaceID: pm.MarketplaceID,
			ExternalID:    pm.ExternalID,
			ProductURL:    pm.ProductURL,
			Price:         pm.Price,
			ListingType:   pm.ListingType,
			Status:        pm.Status,
			Active:        pm.Active,
			CreatedAt:     pm.CreatedAt,
			UpdatedAt:     pm.UpdatedAt,
		})
	}

	return dto
}

func FromProducts(products []model.Product) []ProductDTO {
	result := make([]ProductDTO, 0, len(products))

	for i := range products {
		dto := FromProduct(&products[i])
		if dto != nil {
			result = append(result, *dto)
		}
	}

	return result
}

func (d *ProductDTO) ToModel() *model.Product {
	if d == nil {
		return nil
	}

	m := &model.Product{
		ID:            d.ID,
		Name:          d.Name,
		SKU:           d.SKU,
		Price:         d.Price,
		Marca:         d.Brand,
		Description:   d.Description,
		Active:        d.Active,
		Volume:        d.Volume,
		CategoryID:    d.CategoryID,
		NCM:           d.NCM,
		EAN:           d.EAN,
		Weight:        d.Weight,
		IsKit:         d.IsKit,
		KitComponents: ProductKitDTOsToModel(d.KitItemsDTOs),
	}

	for _, p := range d.Properties {
		m.Properties = append(m.Properties, model.ProductProperty{
			ID:        p.ID,
			ProductID: p.ProductID,
			Name:      p.Name,
			Value:     p.Value,
		})
	}

	for _, i := range d.Images {
		m.Images = append(m.Images, model.ProductImage{
			ID:        i.ID,
			ProductID: i.ProductID,
			URL:       i.URL,
			Position:  i.Position,
			IsPrimary: i.IsPrimary,
		})
	}

	for _, v := range d.Videos {
		m.Videos = append(m.Videos, model.ProductVideo{
			ID:        v.ID,
			ProductID: v.ProductID,
			URL:       v.URL,
			Provider:  v.Provider,
		})
	}

	for _, pm := range d.ProductMarketplaceDTOs {
		m.ProductMarketplaces = append(m.ProductMarketplaces, model.ProductMarketplace{
			ID:            pm.ID,
			ProductID:     pm.ProductID,
			MarketplaceID: pm.MarketplaceID,
			ExternalID:    pm.ExternalID,
			ProductURL:    pm.ProductURL,
			Price:         pm.Price,
			ListingType:   pm.ListingType,
			Status:        pm.Status,
			Active:        pm.Active,
			CreatedAt:     pm.CreatedAt,
			UpdatedAt:     pm.UpdatedAt,
		})
	}

	return m
}

type ProductHomeDTO struct {
	ID                     int64                   `json:"id"`
	SKU                    string                  `json:"sku"`
	Name                   string                  `json:"name"`
	Brand                  string                  `json:"brand"`
	Weight                 *decimal.Decimal        `json:"weight,omitempty"`
	Description            string                  `json:"description"`
	Category               string                  `json:"category"`
	Price                  decimal.Decimal         `json:"price"`
	Rating                 int                     `json:"rating"`
	Reviews                int                     `json:"reviews"`
	Images                 []string                `json:"images"`
	CurrentIndex           int                     `json:"currentIndex"`
	Available              bool                    `json:"available"`
	Properties             []ProductPropertyDTO    `json:"properties,omitempty"`
	ProductMarketplaceDTOs []ProductMarketplaceDTO `json:"product_marketplaces,omitempty"`
	IsKit                  bool                    `json:"is_kit"` // indica se é um kit (produto pai)
	KitItemsDTOs           []ProductKitDTO         `json:"kit_items,omitempty"`
}

type ProductPriceDTO struct {
	ID            int64           `json:"id"`
	ProductID     int64           `json:"product_id"`
	MarketplaceID int64           `json:"marketplace_id"`
	Price         decimal.Decimal `json:"price"`
	Active        bool            `json:"active"`

	Product     *ProductDTO     `json:"product,omitempty"`
	Marketplace *MarketplaceDTO `json:"marketplace,omitempty"`
}

func ToProductHomeDTO(p *ProductDTO) ProductHomeDTO {
	// Converte imagens para slice de URLs
	imageURLs := make([]string, 0, len(p.Images))
	for _, img := range p.Images {
		if img.URL != "" {
			imageURLs = append(imageURLs, img.URL)
		}
	}

	for _, v := range p.Videos {
		if v.URL != "" {
			imageURLs = append(imageURLs, v.URL)
		}
	}

	categoryName := ""
	if p.Category != nil {
		categoryName = p.Category.Name
	}

	for _, pp := range p.Properties {
		p.Properties = append(p.Properties, ProductPropertyDTO{
			ID:        pp.ID,
			ProductID: pp.ProductID,
			Name:      pp.Name,
			Value:     pp.Value,
		})
	}

	return ProductHomeDTO{
		ID:                     p.ID,
		SKU:                    p.SKU,
		Name:                   p.Name,
		Brand:                  p.Brand,
		Weight:                 p.Weight,
		Description:            p.Description,
		Category:               categoryName,
		Price:                  p.Price,
		Rating:                 0, // ajustar conforme necessário
		Reviews:                0, // ajustar conforme necessário
		Images:                 imageURLs,
		CurrentIndex:           0,
		Properties:             p.Properties,
		ProductMarketplaceDTOs: p.ProductMarketplaceDTOs,
		KitItemsDTOs:           p.KitItemsDTOs,
		IsKit:                  p.IsKit,
		Available:              p.Active,
	}
}

func ToProductHomeDTOList(products []ProductDTO) []ProductHomeDTO {
	result := make([]ProductHomeDTO, 0, len(products))

	for i := range products {
		dto := ToProductHomeDTO(&products[i])
		result = append(result, dto)
	}

	return result
}

func GetPrice(product ProductDTO, marketplaceID int64) decimal.Decimal {
	price := product.Price
	for _, p := range product.Prices {
		if p.MarketplaceID == marketplaceID && p.Active {
			price = p.Price
			break
		}
	}
	return price
}

func ProductKitToDTO(m model.ProductKit) ProductKitDTO {
	return ProductKitDTO{
		ID:                 m.ID,
		ProductID:          m.ProductID,
		ComponentProductID: m.ComponentProductID,
		Quantity:           m.Quantity,
	}
}

func ProductKitsToDTO(items []model.ProductKit) []ProductKitDTO {
	if len(items) == 0 {
		return nil
	}

	dtos := make([]ProductKitDTO, 0, len(items))
	for _, item := range items {
		dtos = append(dtos, ProductKitToDTO(item))
	}

	return dtos
}

func ProductKitDTOToModel(d ProductKitDTO) model.ProductKit {
	return model.ProductKit{
		ID:                 d.ID,
		ProductID:          d.ProductID,
		ComponentProductID: d.ComponentProductID,
		Quantity:           d.Quantity,
	}
}

func ProductKitDTOsToModel(items []ProductKitDTO) []model.ProductKit {
	if len(items) == 0 {
		return nil
	}

	models := make([]model.ProductKit, 0, len(items))
	for _, item := range items {
		models = append(models, ProductKitDTOToModel(item))
	}

	return models
}
