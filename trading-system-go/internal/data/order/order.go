package order

import (
	"fmt"
	order_kafka_producer "trading-system-go/api/kafka_config/producer"
	"trading-system-go/database"

	"github.com/gofiber/fiber/v2"
)

type Order struct {
	OrderId     int         `gorm:"AUTO_INCREMENT;PRIMARY_KEY;not null" json:"-"`
	UserId      int32       `json:"user_id"`
	OrderType   OrderType   `json:"order_type"`
	Type        string      `json:"type"`
	Price       float64     `json:"price"`
	Volume      float64     `json:"volume"`
	BuyingPair  string      `json:"buying_pair"`
	SellingPair string      `json:"selling_pair"`
	OrderStatus OrderStatus `json:"-"`
	CreatedAt   int64       `gorm:"autoCreateTime" json:"-"`
}

type OrderType int

const (
	_ OrderType = iota
	MarketOrderBuy
	LimitOrderBuy
	MarketOrderSell
	LimitOrderSell
)

type OrderStatus string

const (
	PartialyFilled OrderStatus = "PartialyFilled"
	Filled         OrderStatus = "Filled"
	NotFilled      OrderStatus = "NotFilled"
)

// create a order
func CreateOrder(order *Order) (*Order, error) {
	order.OrderStatus = NotFilled
	//first create the order in DB
	err := database.DB.DataBase.Create(&order).Error
	if err != nil {
		fmt.Println("error is.......")
		fmt.Println(err)
		return nil, err
	}

	go order_kafka_producer.OrderMatchProducer("topic", order.OrderId)

	return order, nil
}

// get all the order details
func GetOrders(c *fiber.Ctx) (*[]Order, error) {
	var orderList []Order

	err := database.DB.DataBase.Find(&orderList).Error
	if err != nil {
		return nil, err
	}
	return &orderList, nil
}
