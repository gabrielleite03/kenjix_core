package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gabrielleite03/kenjix_core/internal/config"
	"github.com/gabrielleite03/kenjix_core/internal/mercadolivre"
	ml "github.com/gabrielleite03/kenjix_core/internal/service/mercadolivre"
)

type MercadolivreHandler struct {
	authService  *mercadolivre.AuthService
	OrderService ml.OrderService
}

type WebhookPayload struct {
	Resource string `json:"resource"`
	Topic    string `json:"topic"`
	UserID   int64  `json:"user_id"`
}

func NewMercadolivreHandler(service *mercadolivre.AuthService, orderService ml.OrderService) *MercadolivreHandler {
	return &MercadolivreHandler{authService: service, OrderService: orderService}
}

func (h *MercadolivreHandler) GetTGHandler(w http.ResponseWriter, r *http.Request) {
	mercadolivre.RedirectToMercadoLivreAuth(w, r)
}

func (h *MercadolivreHandler) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "code não encontrado", http.StatusBadRequest)
		return
	}

	cfg := config.Load()

	token, err := mercadolivre.GetToken(
		r.Context(),
		cfg.ClientID,
		cfg.ClientSecret,
		code,
		"https://kenjipet.com.br/redirect_mercadolivre",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.authService.SaveToken(r.Context(), token); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, `{
		"status": "ok",
		"user_id": %d,
		"expires_at": %q
	}`, token.UserID, token.ExpiresAt.Format("2006-01-02T15:04:05Z07:00"))
}

func (h *MercadolivreHandler) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var payload WebhookPayload

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "payload inválido", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	if payload.Topic != "orders" {
		return
	}

	//	go h.processWebhook(payload)

	h.processWebhook(payload)
}

func (h *MercadolivreHandler) processWebhook(payload WebhookPayload) {
	orderID := extractOrderID(payload.Resource)
	if orderID == "" {
		log.Println("id inválido:", payload.Resource)
		return
	}

	client := &mercadolivre.Client{
		BaseURL: "https://api.mercadolibre.com",
		Auth:    h.authService,
		UserID:  payload.UserID,
	}

	order, err := client.GetOrderByID(context.Background(), orderID)
	if err != nil {
		log.Println("erro ao buscar pedido:", err)
		return
	}

	if err := h.OrderService.ProcessOrder(order); err != nil {
		log.Println("erro ao processar pedido:", err)
	}
}

func extractOrderID(resource string) string {
	parts := strings.Split(resource, "/")
	if len(parts) == 0 {
		return ""
	}

	return parts[len(parts)-1]
}
