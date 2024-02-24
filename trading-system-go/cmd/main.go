package main

import (
	"fmt"
	"log"
	"trading-system-go/internal/data/balance"
	"trading-system-go/internal/database"
	"trading-system-go/internal/order/adapter"
	order "trading-system-go/internal/order/core/domain"
	"trading-system-go/internal/order/core/service"
	order_book_adapter "trading-system-go/internal/order_book/adapter"
	order_book_service "trading-system-go/internal/order_book/core/service"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

func main() {
	fmt.Print("Aplication Started...")
	app := fiber.New()
	//connect with db
	postgresDbConn := database.ConnectDatabase()
	redisClient := database.ConnectToRedis()

	//orderbook
	orderBookRepository := order_book_adapter.NewOrderBookepository(redisClient)
	orderBookService := order_book_service.NewOrderBookService(orderBookRepository)
	orderBookHandler := order_book_adapter.NewOrderBookHandler(orderBookService)
	orderBookAdapter := order_book_adapter.NewOrderBookKafkaAdpater(orderBookService)
	orderBookKafkaService := order_book_service.NewOrderBookKakfkaService(orderBookAdapter)
	app.Get("/order-book-update", websocket.New(orderBookHandler.ConnectWithOrderBook))
	
	//order
	orderRepository := adapter.NewOrderRepository(postgresDbConn)
	orderMatchService := service.NewOrderMatchService(orderRepository, orderBookService, orderBookKafkaService)
	kafkaAdapter := adapter.NewKafkaAdpater(orderMatchService)
	kafkaService := service.NewKafkaService(kafkaAdapter, orderMatchService)
	orderService := service.NewOrderService(orderRepository, kafkaService)
	orderHandler := adapter.NewHandler(orderService)
	app.Post("/create-order", orderHandler.CreateOrder)
	app.Get("/get-orders", orderHandler.GetOrders)


	postgresDbConn.AutoMigrate(&order.Order{}, &balance.Balance{})
	//router.AppRoutes(app)
	go kafkaService.OrderMatchConsumer("topic")
	go orderBookKafkaService.CreateBidConsumer("order-book")
	log.Fatal(app.Listen(":8000"))
}
