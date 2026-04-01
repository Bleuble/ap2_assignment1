package usecase

import (
	"fmt"
	"order-service/internal/domain"
	"time"
	"github.com/google/uuid"
)
type OrderUseCase struct {
	repo          domain.OrderRepository
	paymentClient domain.PaymentClient
}
func NewOrderUseCase(repo domain.OrderRepository, paymentClient domain.PaymentClient) *OrderUseCase {
	return &OrderUseCase{
		repo:          repo,
		paymentClient: paymentClient,
	}
}

func (uc *OrderUseCase) CreateOrder(customerID, itemName string, amount int64, idempotencyKey string) (*domain.Order, error) {
	if idempotencyKey != "" {
		existingOrder, err := uc.repo.GetByIdempotencyKey(idempotencyKey)
		if err == nil && existingOrder != nil {
			return existingOrder, nil
		}
	}
	order := &domain.Order{
		ID:             uuid.New().String(),
		CustomerID:     customerID,
		ItemName:       itemName,
		Amount:         amount,
		Status:         "Pending",
		IdempotencyKey: idempotencyKey,
		CreatedAt:      time.Now(),
	}
	if err := order.Validate(); err != nil {
		return nil, err
	}
	if err := uc.repo.Create(order); err != nil {
		return nil, fmt.Errorf("failed to save order: %v", err)
	}
	_, err := uc.paymentClient.AuthorizePayment(order.ID, order.Amount)
	if err != nil {
		order.DanaFailed()
		uc.repo.UpdateStatus(order.ID, order.Status)
		return order, fmt.Errorf("payment failed: %v", err)
	}
	order.DanaPaid()
	uc.repo.UpdateStatus(order.ID, order.Status)
	return order, nil
}
