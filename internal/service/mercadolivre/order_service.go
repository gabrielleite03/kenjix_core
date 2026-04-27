package mercadolivre

import (
	"log"

	"github.com/gabrielleite03/kenjix_core/internal/mercadolivre"
)

type OrderRepositoryInterface interface {
	Exists(id int64) (bool, error)
	Create(order mercadolivre.Order) error
}

type OrderService struct {
	Repo OrderRepositoryInterface
}

func (s *OrderService) Save(order mercadolivre.Order) error {
	exists, _ := s.Repo.Exists(order.ID)
	if exists {
		return nil
	}

	log.Println("salvando pedido:", order.ID)
	return s.Repo.Create(order)
}
