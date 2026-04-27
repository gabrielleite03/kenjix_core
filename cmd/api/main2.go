package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gabrielleite03/kenjix_core/internal/config"
	"github.com/gabrielleite03/kenjix_core/internal/database"
	"github.com/gabrielleite03/kenjix_core/internal/mercadolivre"
	"github.com/gabrielleite03/kenjix_core/internal/repository"
	ml_service "github.com/gabrielleite03/kenjix_core/internal/service/mercadolivre"
)

func main2() {
	cfg := config.Load()

	db := database.Connect(cfg.DBUrl)

	repo := &repository.OrderRepository{DB: db}
	service := &ml_service.OrderService{Repo: repo}

	auth := &mercadolivre.AuthService{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Token: &mercadolivre.Token{
			AccessToken:  "INITIAL_TOKEN",
			RefreshToken: "REFRESH_TOKEN",
			ExpiresAt:    time.Now(),
		},
	}

	client := &mercadolivre.Client{
		BaseURL: "https://api.mercadolibre.com",
		Auth:    auth,
	}

	webhook := &mercadolivre.WebhookHandler{
		Client:       client,
		OrderService: service,
	}

	http.HandleFunc("/webhook", webhook.Handle)

	log.Println("rodando na porta 8080")
	http.ListenAndServe(":8080", nil)
}
