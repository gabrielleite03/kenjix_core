package repository

import (
	"database/sql"

	"github.com/gabrielleite03/kenjix_core/internal/mercadolivre"
)

type OrderRepository struct {
	DB *sql.DB
}

func (r *OrderRepository) Exists(id int64) (bool, error) {
	var exists bool
	err := r.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM orders WHERE id=$1)", id).Scan(&exists)
	return exists, err
}

func (r *OrderRepository) Create(order mercadolivre.Order) error {
	_, err := r.DB.Exec(
		"INSERT INTO orders (id, status, total) VALUES ($1,$2,$3)",
		order.ID,
		order.Status,
		order.Total,
	)
	return err
}
