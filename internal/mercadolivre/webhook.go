package mercadolivre

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type WebhookPayload struct {
	Resource string `json:"resource"`
	Topic    string `json:"topic"`
	UserID   int64  `json:"user_id"`
}

type WebhookHandler struct {
	Client       *Client
	OrderService OrderServiceInterface
}

type OrderServiceInterface interface {
	Save(order Order) error
}

func (h *WebhookHandler) Handle(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var payload WebhookPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	// responder rápido pro ML
	w.WriteHeader(http.StatusOK)

	// processar async (ok, mas controlado)
	go func(p WebhookPayload) {
		if err := h.process(r.Context(), p); err != nil {
			log.Println("erro webhook:", err)
		}
	}(payload)
}

func (h *WebhookHandler) process(ctx context.Context, p WebhookPayload) error {
	if p.Topic != "orders" {
		return nil
	}

	id := extractID(p.Resource)
	if id == "" {
		return fmt.Errorf("id inválido: %s", p.Resource)
	}

	order, err := h.Client.GetOrderByID(ctx, id)
	if err != nil {
		return err
	}

	return h.OrderService.Save(*order)
}

func extractID(resource string) string {
	parts := strings.Split(resource, "/")
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}
