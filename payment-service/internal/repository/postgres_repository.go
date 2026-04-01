package repository

import (
	"database/sql"
	"order-service/internal/domain"
)

type PostgresOrderRepository struct {
	db *sql.DB
}

func NewPostgresOrderRepository(db *sql.DB) *PostgresOrderRepository {
	return &PostgresOrderRepository{db: db}
}

func (r *PostgresOrderRepository) Create(order *domain.Order) error {
	query := `
		INSERT INTO orders (id, customer_id, item_name, amount, status, idempotency_key, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`
	
	_, err := r.db.Exec(query,
		order.ID, order.CustomerID, order.ItemName, order.Amount,
		order.Status, order.IdempotencyKey, order.CreatedAt,
	)
	return err
}

func (r *PostgresOrderRepository) GetByID(id string) (*domain.Order, error) {
	query := `SELECT id, customer_id, item_name, amount, status, idempotency_key, created_at FROM orders WHERE id = $1`
	
	var o domain.Order
	err := r.db.QueryRow(query, id).Scan(
		&o.ID, &o.CustomerID, &o.ItemName, &o.Amount,
		&o.Status, &o.IdempotencyKey, &o.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	
	return &o, nil
}

func (r *PostgresOrderRepository) GetByIdempotencyKey(key string) (*domain.Order, error) {
	query := `SELECT id, customer_id, item_name, amount, status, idempotency_key, created_at FROM orders WHERE idempotency_key = $1`
	
	var o domain.Order
	err := r.db.QueryRow(query, key).Scan(
		&o.ID, &o.CustomerID, &o.ItemName, &o.Amount,
		&o.Status, &o.IdempotencyKey, &o.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	
	return &o, nil
}

func (r *PostgresOrderRepository) UpdateStatus(id string, status string) error {
	query := `UPDATE orders SET status = $1 WHERE id = $2`
	_, err := r.db.Exec(query, status, id)
	return err
}

