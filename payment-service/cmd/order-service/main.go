package main

import (
	"database/sql"
	"log"

	"order-service/internal/app"
	"order-service/internal/repository"
	"order-service/internal/transport/http"
	"order-service/internal/usecase"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	connStr := "postgres://user:password@localhost:5432/order_db?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}

	orderRepo := repository.NewPostgresOrderRepository(db)
	paymentClient := app.NewHttpPaymentClient("http://localhost:8081")
	orderUseCase := usecase.NewOrderUseCase(orderRepo, paymentClient)

	router := gin.Default()
	orderHandler := http.NewOrderHandler(orderUseCase)
	orderHandler.RegisterRoutes(router)

	log.Println("Order Service is running on port 8080...")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}


