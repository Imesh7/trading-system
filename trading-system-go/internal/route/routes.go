package routes

import (
	"fmt"
	"trading-system-go/internal/data/order"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

var Conn map[string][]*websocket.Conn

func init() {
	Conn = make(map[string][]*websocket.Conn)
}

func AppRoutes(app *fiber.App) {
	app.Post("/create-order", CreateOrderRoute)
	app.Get("/get-orders", GetOrdersRoute)
	app.Get("/order-book-update", websocket.New(UpdateOrderBook))
}

func CreateOrderRoute(c *fiber.Ctx) error {
	var orderData order.Order
	err := c.BodyParser(&orderData)
	if err != nil {
		return c.Status(402).JSON(&orderData)
	}
	fmt.Println()
	if orderData.UserId == 0 || (orderData.OrderType <= 0 || orderData.OrderType > 4) || orderData.BuyingPair == "" || orderData.SellingPair == "" || orderData.Price == 0 || orderData.Volume == 0 {
		s := fmt.Sprintf("Unprocessable Entity %d", orderData.OrderType)
		return c.Status(402).JSON(s)
	}
	order, err := order.CreateOrder(&orderData)
	if err != nil {
		fmt.Println(err)
		return c.Status(400).JSON("something went wrong on database cannot create a order")
	}
	return c.Status(200).JSON(&order)
}

func GetOrdersRoute(c *fiber.Ctx) error {
	orderList, err := order.GetOrders(c)
	if err != nil {
		return c.Status(400).JSON(err)
	}
	return c.Status(200).JSON(orderList)
}

func GetOrderBook(c *fiber.Ctx) error {
	orderList, err := order.GetOrders(c)
	if err != nil {
		return c.Status(400).JSON(err)
	}
	return c.Status(200).JSON(orderList)
}

func UpdateOrderBook(socket *websocket.Conn) {
	pairConn := socket.Query("pair")

	Conn[pairConn] = append(Conn[pairConn], socket)
	defer socket.Close()
	fmt.Println("received websocket conn req.........................")
	var msg any
	for {
		err := socket.ReadJSON(&msg)
		if err != nil {
			// optional: log the error
			break
		}
	}
}
