// internal/router/router.go
package router

import (
	"fmt"
	"net/http"

	"github.com/gabrielleite03/kenjix_core/cmd/api/handler"
	"github.com/gabrielleite03/kenjix_core/cmd/api/middleware"
	"github.com/gabrielleite03/kenjix_core/internal/service"
	persist "github.com/gabrielleite03/kenjix_persist/repository"
)

const (
	AuthServiceURL = "http://koto-server01:81"
)

type Router struct {
	authMiddleware  *middleware.AuthMiddleware
	authHandler     *handler.AuthHandler
	productHandler  *handler.ProductHandler
	categoryHandler *handler.CategoryHandler
}

func NewRouter() *Router {
	repo := persist.NewProductDAO()
	svc := service.NewProductService(repo)
	h := handler.NewProductHandler(svc)
	return &Router{
		authHandler:     handler.NewAuthHandler(service.NewAuthService(AuthServiceURL)),
		authMiddleware:  middleware.NewAuthMiddleware(service.NewAuthService(AuthServiceURL)),
		productHandler:  h,
		categoryHandler: handler.NewCategoryHandler(service.NewCategoryService(persist.NewCategoryDAO())),
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
	// Product
	http.HandleFunc("/products", r.handleProducts)
	http.HandleFunc("/products/", r.handleProduct)

	// Category
	http.HandleFunc("/categories", r.handleCategories)
	http.HandleFunc("/categories/", r.handleCategory)

	// Login
	http.HandleFunc("/login", r.authHandler.Login)

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
	case http.MethodPost:
		r.authMiddleware.Middleware(r.categoryHandler.Create)(w, req)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (r *Router) handleCategory(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		r.categoryHandler.Get(w, req)
	case http.MethodPut:
		r.authMiddleware.Middleware(r.categoryHandler.Update)(w, req)
	case http.MethodDelete:
		r.authMiddleware.Middleware(r.categoryHandler.Delete)(w, req)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
