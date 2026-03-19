package main

import (
	"log"
	"net/http"

	"github.com/gabrielleite03/kenjix_core/cmd/api/router"
)

func main() {
	r := router.NewRouter()
	log.Println("Kenjix Core iniciado na porta 7010")
	http.ListenAndServe(":7010", r.RegisterRoutes())
}
