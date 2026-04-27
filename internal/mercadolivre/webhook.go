package mercadolivre

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type WebhookPayload struct {
	Resource string `json:"resource"`
	Topic    string `json:"topic"`
}

type WebhookHandler struct {
	Client       *Client
	OrderService OrderServiceInterface
}

type OrderServiceInterface interface {
	Save(order Order) error
}

func (h *WebhookHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var payload WebhookPayload
	json.NewDecoder(r.Body).Decode(&payload)

	go h.process(payload)

	w.WriteHeader(http.StatusOK)
}

func (h *WebhookHandler) process(p WebhookPayload) {
	if p.Topic != "orders" {
		return
	}

	id := extractID(p.Resource)

	order, err := h.Client.GetOrderByID(context.Background(), id)
	if err != nil {
		log.Println("erro ao buscar pedido:", err)
		return
	}

	h.OrderService.Save(*order)
}

func extractID(resource string) string {
	// "/orders/123456"
	return resource[len("/orders/"):]
}

func rContext() *http.Request {
	return &http.Request{}
}
