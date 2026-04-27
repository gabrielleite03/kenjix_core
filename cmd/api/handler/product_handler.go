package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gabrielleite03/kenjix_core/internal/dto"
	"github.com/gabrielleite03/kenjix_core/internal/service"
	"github.com/shopspring/decimal"
)

type ProductHandler struct {
	ProductService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{ProductService: productService}
}

// GET /products
func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	marketplaceParam := r.URL.Query().Get("marketplace")

	if marketplaceParam != "" {
		products, err := h.ProductService.ListProductsByMarketplace(context.Background(), marketplaceParam)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Erro ao buscar produto por marketplace")
			return
		}
		writeJSON(w, http.StatusOK, products)
		return
	}
	products, err := h.ProductService.ListProducts(r.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	productsDTO := make([]*dto.ProductDTO, len(products))
	for i, p := range products {
		productsDTO[i] = dto.FromProduct(&p)
	}

	writeJSON(w, http.StatusOK, productsDTO)
}

// GET /products/{id}
func (h *ProductHandler) Get(w http.ResponseWriter, r *http.Request) {
	marketplaceParam := r.URL.Query().Get("marketplace")
	id := getID(r)

	product, err := h.ProductService.GetProductByMarketplace(context.Background(), id, marketplaceParam)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusOK, product)

}

// POST /products
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(50 << 50) // 50MB
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	categoryID, _ := strconv.ParseInt(r.FormValue("category"), 10, 64)

	volume, err := decimal.NewFromString(r.FormValue("volume"))
	if err != nil {
		http.Error(w, "Invalid volume", http.StatusBadRequest)
		return
	}

	ean := r.FormValue("ean")
	ncm := r.FormValue("ncm")

	productDTO := dto.ProductDTO{
		Name:        r.FormValue("name"),
		SKU:         r.FormValue("sku"),
		CategoryID:  &categoryID,
		Brand:       r.FormValue("brand"),
		Description: r.FormValue("description"),
		Volume:      volume,
		Active:      r.FormValue("active") == "true",
		ImageFiles:  r.MultipartForm.File["images"],
		Videos: []dto.ProductVideoDTO{
			{URL: r.FormValue("videoUrl")},
		},
		EAN: &ean,
		NCM: &ncm,
	}

	// propriedades
	properties := r.FormValue("properties")
	if properties != "" {
		json.Unmarshal([]byte(properties), &productDTO.Properties)
	}

	_, kk := h.ProductService.CreateProduct(r.Context(), &productDTO)
	if kk != nil {
		http.Error(w, kk.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	// implementar
}

// PUT /products/{id}
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r.URL.Path)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}
	err = r.ParseMultipartForm(50 << 50) // 50MB
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	categoryID, _ := strconv.ParseInt(r.FormValue("category"), 10, 64)

	volume, err := decimal.NewFromString(r.FormValue("volume"))
	if err != nil {
		http.Error(w, "Invalid volume", http.StatusBadRequest)
		return
	}

	ean := r.FormValue("ean")
	ncm := r.FormValue("ncm")

	productDTO := dto.ProductDTO{
		Name:        r.FormValue("name"),
		SKU:         r.FormValue("sku"),
		CategoryID:  &categoryID,
		Brand:       r.FormValue("brand"),
		Description: r.FormValue("description"),
		Volume:      volume,
		Active:      r.FormValue("active") == "true",
		ImageFiles:  r.MultipartForm.File["images"],
		Videos: []dto.ProductVideoDTO{
			{URL: r.FormValue("videoUrl")},
		},
		EAN: &ean,
		NCM: &ncm,
	}
	productDTO.ID = id

	existing := r.FormValue("existingImages")
	var existingImages []string

	if existing != "" {
		err := json.Unmarshal([]byte(existing), &existingImages)
		if err != nil {
			http.Error(w, "invalid existingImages", http.StatusBadRequest)
			return
		}

		productDTO.ExistingImages = existingImages
	}

	deleting := r.FormValue("deletedImages")
	var deletedImages []string

	if deleting != "" {
		err := json.Unmarshal([]byte(deleting), &deletedImages)
		if err != nil {
			http.Error(w, "invalid deletedImages", http.StatusBadRequest)
			return
		}

		productDTO.DeletedImages = deletedImages
	}

	// propriedades
	properties := r.FormValue("properties")
	if properties != "" {
		json.Unmarshal([]byte(properties), &productDTO.Properties)
	}

	err = h.ProductService.UpdateProduct(r.Context(), &productDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

}

// DELETE /products/{id}
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// implementar
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
