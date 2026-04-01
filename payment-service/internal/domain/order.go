package domain

import (
	"fmt"
	"time"
)

type Order struct {
	ID             string
	CustomerID     string
	ItemName       string
	Amount         int64
	Status         string
	IdempotencyKey string
	CreatedAt      time.Time
}

func (o *Order) Validate() error {
	if o.CustomerID == "" {
		return fmt.Errorf("customer_id is required")
	}
	if o.ItemName == "" {
		return fmt.Errorf("item_name is required")
	}
	if o.Amount <= 0 {
		return fmt.Errorf("amount must be greater than zero")
	}
	return nil
}