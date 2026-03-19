// internal/router/router.go
package router

import (
	"fmt"
	"net/http"

	"github.com/gabrielleite03/kenjix_core/cmd/api/handler"
	"github.com/gabrielleite03/kenjix_core/internal/service"
	persist "github.com/gabrielleite03/kenjix_persist/repository"
)

type Router struct {
	productHandler *handler.ProductHandler
}

func NewRouter() *Router {
	repo := persist.NewProductRepository()
	svc := service.NewProductService(repo)
	h := handler.NewProductHandler(svc)
	return &Router{
		productHandler: h,
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
