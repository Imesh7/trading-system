package adapter

import (
	"fmt"
	order "trading-system-go/internal/order/core/domain"
	"trading-system-go/internal/order/core/ports/service"

	"github.com/gofiber/fiber/v2"
)

type orderHanlder struct {
	service ports.OrderService
}

func NewHandler(service ports.OrderService) *orderHanlder {
	return &orderHanlder{
		service: service,
	}
}

func (handler *orderHanlder) CreateOrder(c *fiber.Ctx) error {
	var orderData order.Order
	err := c.BodyParser(&orderData)
	if err != nil {
		return c.Status(402).JSON(&orderData)
	}
	if orderData.UserId == 0 || (orderData.OrderType <= 0 || orderData.OrderType > 4) || orderData.BuyingPair == "" || orderData.SellingPair == "" || orderData.Price == 0 || orderData.Volume == 0 {
		s := fmt.Sprintf("Unprocessable Entity %d", orderData.OrderType)
		return c.Status(402).JSON(s)
	}
	order, err := handler.service.CreateOrderService(&orderData)
	if err != nil {
		fmt.Println(err)
		return c.Status(400).JSON("something went wrong on database cannot create a order")
	}
	return c.Status(200).JSON(&order)
}

func (handler *orderHanlder) GetOrders(c *fiber.Ctx) error {
	orderList, err := handler.service.GetOrdersService()
	if err != nil {
		return c.Status(400).JSON(err)
	}
	return c.Status(200).JSON(orderList)
}
