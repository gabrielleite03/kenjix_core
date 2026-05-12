package mercadolivre

import (
	"log"

	"github.com/gabrielleite03/kenjix_core/internal/mercadolivre"
)

type OrderService interface {
	ProcessOrder(order *mercadolivre.Order) error
}

type orderServiceImpl struct {
	// aqui você pode injetar dependências, como repositórios ou serviços de estoque
}

func NewOrderService() OrderService {
	return &orderServiceImpl{}
}

func (h *orderServiceImpl) ProcessOrder(order *mercadolivre.Order) error {
	log.Println("processando pedido:", order.ID)

	// 🔥 só processa pedidos pagos
	if !(order.IsPaid()) {
		log.Println("pedido não pago:", order.ID)
		return nil
	}
	if order.IsDeliveredOrPickedUpManually() {
		log.Println("pedido retirado em mãos / sem envio:", order.ID)

		// no ERP:
		// delivery_status = "retirado_em_maos"
		// stock_status = "baixado"
	}

	for _, item := range order.OrderItems {
		productID := item.Item.ID
		qty := item.Quantity

		log.Println("baixando estoque:", productID, qty)
		sku := item.Item.SellerSKU
		log.Println("SKU do item:", sku)

		// aqui você integra com seu estoque:
		// h.stockService.DecreaseStock(productID, qty)
	}

	return nil
}
