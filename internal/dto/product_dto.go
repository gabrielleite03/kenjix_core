package dto

import (
	"mime/multipart"

	model "github.com/gabrielleite03/kenjix_domain/model"
	"github.com/shopspring/decimal"
)

type ProductDTO struct {
	ID          int64           `json:"id"`
	Name        string          `json:"name"`
	SKU         string          `json:"sku"`
	Price       decimal.Decimal `json:"price"`
	Brand       string          `json:"brand"`
	Description string          `json:"description"`
	Active      bool            `json:"active"`
	Volume      decimal.Decimal `json:"volume"`
	CategoryID  *int64          `json:"category_id,omitempty"`
	Category    *CategoryDTO    `json:"category,omitempty"`

	Prices []ProductPriceDTO `json:"prices,omitempty"` // preços por marketplace
	Stocks []StockDTO        `json:"stocks,omitempty"` // estoque por warehouse

	Properties     []ProductPropertyDTO    `json:"properties,omitempty"`
	Images         []ProductImageDTO       `json:"images,omitempty"`
	ExistingImages []string                `json:"existingImages,omitempty"`
	DeletedImages  []string                `json:"deletedImages,omitempty"`
	ImageFiles     []*multipart.FileHeader `json:"image_files,omitempty"`
	Videos         []ProductVideoDTO       `json:"videos,omitempty"`
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

func FromProduct(m *model.Product) *ProductDTO {
	if m == nil {
		return nil
	}

	dto := &ProductDTO{
		ID:          m.ID,
		Name:        m.Name,
		SKU:         m.SKU,
		Price:       m.Price,
		Brand:       m.Marca,
		Description: m.Description,
		Active:      m.Active,
		Volume:      m.Volume,
		CategoryID:  m.CategoryID,
		Category:    (*CategoryDTO)(m.Category),
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

	return dto
}

func (d *ProductDTO) ToModel() *model.Product {
	if d == nil {
		return nil
	}

	m := &model.Product{
		ID:          d.ID,
		Name:        d.Name,
		SKU:         d.SKU,
		Price:       d.Price,
		Marca:       d.Brand,
		Description: d.Description,
		Active:      d.Active,
		Volume:      d.Volume,
		CategoryID:  d.CategoryID,
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

	return m
}

type ProductHomeDTO struct {
	ID           int64                `json:"id"`
	SKU          string               `json:"sku"`
	Name         string               `json:"name"`
	Brand        string               `json:"brand"`
	Description  string               `json:"description"`
	Category     string               `json:"category"`
	Price        decimal.Decimal      `json:"price"`
	Rating       int                  `json:"rating"`
	Reviews      int                  `json:"reviews"`
	Images       []string             `json:"images"`
	CurrentIndex int                  `json:"currentIndex"`
	Available    bool                 `json:"available"`
	Properties   []ProductPropertyDTO `json:"properties,omitempty"`
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
		ID:           p.ID,
		SKU:          p.SKU,
		Name:         p.Name,
		Brand:        p.Brand,
		Description:  p.Description,
		Category:     categoryName,
		Price:        p.Price,
		Rating:       0, // ajustar conforme necessário
		Reviews:      0, // ajustar conforme necessário
		Images:       imageURLs,
		CurrentIndex: 0,
		Properties:   p.Properties,
	}
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
