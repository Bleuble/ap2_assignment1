package http

import (
	"net/http"
	"order-service/internal/usecase"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	useCase *usecase.OrderUseCase
}

func NewOrderHandler(uc *usecase.OrderUseCase) *OrderHandler {
	return &OrderHandler{useCase: uc}
}

func (h *OrderHandler) RegisterRoutes(router *gin.Engine) {
	router.POST("/orders", h.CreateOrder)
	router.GET("/orders/:id", h.GetOrder)
	router.PATCH("/orders/:id/cancel", h.CancelOrder)
}

type CreateOrderRequest struct {
	CustomerID string `json:"customer_id"`
	ItemName   string `json:"item_name"`
	Amount     int64  `json:"amount"`
}


func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid format"})
		return
	}

	idempotencyKey := c.GetHeader("Idempotency-Key")

	order, err := h.useCase.CreateOrder(req.CustomerID, req.ItemName, req.Amount, idempotencyKey)
	if err != nil {
		if err.Error() == "customer_id is required" || err.Error() == "item_name is required" || err.Error() == "amount must be greater than zero" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

