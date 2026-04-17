// internal/router/router.go
package router

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gabrielleite03/kenjix_core/cmd/api/handler"
	"github.com/gabrielleite03/kenjix_core/cmd/api/middleware"
	"github.com/gabrielleite03/kenjix_core/internal/service"
	persist "github.com/gabrielleite03/kenjix_persist/repository"
)

const (
	AuthServiceURL = "http://localhost:7020"
	//AuthServiceURL = "http://koto-server01:81"
)

type Router struct {
	authMiddleware            *middleware.AuthMiddleware
	authHandler               *handler.AuthHandler
	productHandler            *handler.ProductHandler
	categoryHandler           *handler.CategoryHandler
	warehouseHandler          *handler.WarehouseHandler
	warehousePlaceHandler     *handler.WarehousePlaceHandler
	warehouseLocationsHandler *handler.WarehouseLocationsHandler
	expenseHandler            *handler.ExpenseHandler
	supplierHandler           *handler.SupplierHandler
	costCenterHandler         *handler.CostCenterHandler
	purchaseHandler           *handler.PurchaseHandler
	stockHandler              *handler.StockHandler
	markeplaceHandler         *handler.MarketplaceHandler
}

func NewRouter() *Router {
	repo := persist.NewProductDAO()
	svc := service.NewProductService(repo)
	h := handler.NewProductHandler(svc)
	return &Router{
		authHandler:               handler.NewAuthHandler(service.NewAuthService(AuthServiceURL)),
		authMiddleware:            middleware.NewAuthMiddleware(service.NewAuthService(AuthServiceURL)),
		productHandler:            h,
		categoryHandler:           handler.NewCategoryHandler(service.NewCategoryService(persist.NewCategoryDAO())),
		warehouseHandler:          handler.NewWarehouseHandler(service.NewWarehouseService(persist.NewWarehouseDAO())),
		warehousePlaceHandler:     handler.NewWarehousePlaceHandler(service.NewWarehouseService(persist.NewWarehouseDAO())),
		warehouseLocationsHandler: handler.NewWarehouseLocationsHandler(service.NewWarehouseService(persist.NewWarehouseDAO())),
		expenseHandler:            handler.NewExpenseHandler(service.NewExpenseService(persist.NewExpenseDAO())),
		supplierHandler:           handler.NewSupplierHandler(service.NewSupplierService(persist.NewSupplierDAO())),
		costCenterHandler:         handler.NewCostCenterHandler(*service.NewCostCenterService(persist.NewCostCenterDAO())),
		purchaseHandler:           handler.NewPurchaseHandler(),
		stockHandler:              handler.NewStockHandler(),
		markeplaceHandler:         handler.NewMarketplaceHandler(service.NewMarketplaceService(persist.NewMarketplaceDAO())),
	}
}

func (r *Router) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/products", r.productHandler.List)

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Kenjix Persist API")
	})
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "ok")
	})

	return mux
}

func (r *Router) Register() {

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, `{"status":"ok"}`)
	})

	// Product
	http.HandleFunc("/products", middleware.CorsMiddleware(r.handleProducts))
	http.HandleFunc("/products/", middleware.CorsMiddleware(r.handleProduct))

	// Category
	http.HandleFunc("/categories", middleware.CorsMiddleware(r.handleCategories))
	http.HandleFunc("/categories/", middleware.CorsMiddleware(r.handleCategory))
	http.HandleFunc("/category", middleware.CorsMiddleware(r.handleCategory))

	//Warehouse
	http.HandleFunc("/warehouses", middleware.CorsMiddleware(r.handleWarehouses))
	http.HandleFunc("/warehouses/", middleware.CorsMiddleware(r.handleWarehousesByID))
	http.HandleFunc("/storage-location-types", middleware.CorsMiddleware(r.handleWarehousePlaceTypes))
	http.HandleFunc("/warehouses-storage-locations/", middleware.CorsMiddleware(r.handleWarehouseLocations))

	// Login
	http.HandleFunc("/login", middleware.CorsMiddleware(r.authHandler.Login))

	// Expenses
	http.HandleFunc("/expenses", middleware.CorsMiddleware(r.handleExpense))
	http.HandleFunc("/expenses/", middleware.CorsMiddleware(r.handleExpenseByID))
	http.HandleFunc("/expenses/categories", middleware.CorsMiddleware(r.handleExpenseByCategory))

	// Supplier
	http.HandleFunc("/suppliers", middleware.CorsMiddleware(r.handleSuppliers))
	http.HandleFunc("/suppliers/", middleware.CorsMiddleware(r.handleSuppliersById))

	// Supplier
	http.HandleFunc("/cost-centers", middleware.CorsMiddleware(r.handleCostCenters))
	http.HandleFunc("/cost-centers/", middleware.CorsMiddleware(r.handleCostCentersById))

	// Purchase
	http.HandleFunc("/purchases", middleware.CorsMiddleware(r.handlePurchases))
	http.HandleFunc("/purchases/", middleware.CorsMiddleware(r.handlePurchasesById))

	// Stock
	http.HandleFunc("/stocks", middleware.CorsMiddleware(r.handleStocks))

	// MarketPlaces
	http.HandleFunc("/marketplaces", middleware.CorsMiddleware(r.handleMarketplace))
	http.HandleFunc("/marketplaces/", middleware.CorsMiddleware(r.handleMarketplaceByID))
}

func (r *Router) handleProducts(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		r.productHandler.List(w, req)
	case http.MethodPost:
		r.authMiddleware.Middleware(r.productHandler.Create)(w, req)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (r *Router) handleProduct(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		r.productHandler.Get(w, req)
	case http.MethodPut:
		r.authMiddleware.Middleware(r.productHandler.Update)(w, req)
	case http.MethodDelete:
		r.authMiddleware.Middleware(r.productHandler.Delete)(w, req)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (r *Router) handleCategories(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		r.categoryHandler.List(w, req)
	case http.MethodOptions:
		r.categoryHandler.List(w, req)
	case http.MethodPost:
		r.authMiddleware.Middleware(r.categoryHandler.Create)(w, req)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (r *Router) handleWarehouses(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		r.warehouseHandler.FindAll(w, req)
	case http.MethodOptions:
		r.warehouseHandler.FindAll(w, req)
	case http.MethodPost:
		r.authMiddleware.Middleware(r.warehouseHandler.Create)(w, req)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (r *Router) handleWarehousesByID(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPut:
		r.authMiddleware.Middleware(r.warehouseHandler.Update)(w, req)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (r *Router) handleWarehousePlaceTypes(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		r.warehousePlaceHandler.FindAllPlaceType(w, req)
	case http.MethodOptions:
		r.warehousePlaceHandler.FindAllPlaceType(w, req)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (r *Router) handleWarehouseLocations(w http.ResponseWriter, req *http.Request) {
	// /warehouses-storage-locations/{id}
	path := strings.TrimPrefix(req.URL.Path, "/warehouses-storage-locations/")
	parts := strings.Split(path, "/")

	if len(parts) == 0 || parts[0] == "" {
		http.NotFound(w, req)
		return
	}

	warehouseID, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		http.Error(w, "invalid warehouse id", http.StatusBadRequest)
		return
	}

	switch req.Method {
	case http.MethodGet:
		r.warehouseLocationsHandler.FindAllLocationByWarehouseID(w, req, warehouseID)
	case http.MethodOptions:
		r.warehousePlaceHandler.FindAllPlaceType(w, req)
	case http.MethodPost:
		r.authMiddleware.Middleware(func(w http.ResponseWriter, req *http.Request) {
			r.warehouseLocationsHandler.CreateWarehouseLocation(w, req, warehouseID)
		})(w, req)
	case http.MethodPut:
		r.authMiddleware.Middleware(r.warehouseLocationsHandler.UpdateWarehousePlace)(w, req)
	case http.MethodDelete:
		warehousePlaceID, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			http.Error(w, "invalid warehouse place id", http.StatusBadRequest)
			return
		}
		r.authMiddleware.Middleware(func(w http.ResponseWriter, req *http.Request) {
			r.warehouseLocationsHandler.DeleteWarehouseLocation(w, req, warehousePlaceID)
		})(w, req)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (r *Router) handleCategory(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		r.categoryHandler.Get(w, req)
	case http.MethodPost:
		r.authMiddleware.Middleware(r.categoryHandler.Create)(w, req)
	case http.MethodPut:
		r.authMiddleware.Middleware(r.categoryHandler.Update)(w, req)
	case http.MethodDelete:
		r.authMiddleware.Middleware(r.categoryHandler.Delete)(w, req)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (r *Router) handleExpense(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		r.authMiddleware.Middleware(r.expenseHandler.FindAll)(w, req)
	case http.MethodPost:
		r.authMiddleware.Middleware(r.expenseHandler.Create)(w, req)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (r *Router) handleExpenseByID(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		r.authMiddleware.Middleware(r.expenseHandler.FindByID)(w, req)
	case http.MethodPut:
		r.authMiddleware.Middleware(r.expenseHandler.Update)(w, req)
	case http.MethodDelete:
		r.authMiddleware.Middleware(r.expenseHandler.Delete)(w, req)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (r *Router) handleExpenseByCategory(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		r.authMiddleware.Middleware(r.expenseHandler.FindAllExpenseCategories)(w, req)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (r *Router) handleSuppliers(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		r.authMiddleware.Middleware(r.supplierHandler.List)(w, req)
	case http.MethodPost:
		r.authMiddleware.Middleware(r.supplierHandler.Create)(w, req)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (r *Router) handleSuppliersById(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		r.authMiddleware.Middleware(r.supplierHandler.Get)(w, req)
	case http.MethodPut:
		r.authMiddleware.Middleware(r.supplierHandler.Update)(w, req)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (r *Router) handleCostCenters(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		r.authMiddleware.Middleware(r.costCenterHandler.FindAll)(w, req)
	case http.MethodPost:
		r.authMiddleware.Middleware(r.costCenterHandler.Create)(w, req)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (r *Router) handleCostCentersById(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		r.authMiddleware.Middleware(r.costCenterHandler.FindByID)(w, req)
	case http.MethodPut:
		r.authMiddleware.Middleware(r.costCenterHandler.Update)(w, req)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (r *Router) handlePurchases(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		r.authMiddleware.Middleware(r.purchaseHandler.FindAll)(w, req)
	case http.MethodPost:
		r.authMiddleware.Middleware(r.purchaseHandler.Create)(w, req)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (r *Router) handlePurchasesById(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		r.authMiddleware.Middleware(r.purchaseHandler.FindByID)(w, req)
	case http.MethodPut:
		r.authMiddleware.Middleware(r.purchaseHandler.Update)(w, req)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (r *Router) handleStocks(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		r.authMiddleware.Middleware(r.stockHandler.FindAll)(w, req)
	case http.MethodPost:
		r.authMiddleware.Middleware(r.stockHandler.Create)(w, req)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (r *Router) handleMarketplace(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		r.markeplaceHandler.FindAll(w, req)
	case http.MethodPost:
		r.authMiddleware.Middleware(r.markeplaceHandler.Create)(w, req)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (r *Router) handleMarketplaceByID(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		r.authMiddleware.Middleware(r.markeplaceHandler.FindByID)(w, req)
	case http.MethodPut:
		r.authMiddleware.Middleware(r.markeplaceHandler.Update)(w, req)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
