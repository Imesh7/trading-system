package main

import (
	"fmt"
	"log"
	order_kafka_consumer "trading-system-go/api/kafka_config/consumer"
	"trading-system-go/database"
	"trading-system-go/internal/data/balance"
	"trading-system-go/internal/data/order"
	router "trading-system-go/internal/route"

	"github.com/gofiber/fiber/v2"
)

func main() {
	fmt.Print("aplication started.........................")
	app := fiber.New()
	database.ConnectDatabase()
	database.ConnectToRedis()
	database.DB.DataBase.AutoMigrate(&order.Order{}, &balance.Balance{})
	router.AppRoutes(app)
	go order_kafka_consumer.OrderMatchConsumer("topic")
	go order_kafka_consumer.CreateBidConsumer("order-book")
	log.Fatal(app.Listen(":8000"))
}
