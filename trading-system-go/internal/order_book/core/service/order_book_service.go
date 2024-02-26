package service

import (
	"context"
	"fmt"
	"log"
	order "trading-system-go/internal/order/core/domain"
	"trading-system-go/internal/order_book/core/ports"

	"github.com/gofiber/contrib/websocket"
)

type orderBookService struct {
	repository ports.OrderBookRepository
}

var socketConnections map[string]map[string]*websocket.Conn

func NewOrderBookService(repository ports.OrderBookRepository) *orderBookService {
	socketConnections = make(map[string]map[string]*websocket.Conn)
	return &orderBookService{
		repository: repository,
	}
}

func (service *orderBookService) ConnectWithOrderBook(socket *websocket.Conn) {
	pairConn := socket.Query("pair")
	userId := socket.Query("userid")

	if socketConnections[pairConn] == nil {
		socketConnections[pairConn] = make(map[string]*websocket.Conn)
	}

	socketConnections[pairConn][userId] = socket

	defer socket.Close()
	fmt.Println("Received websocket reqest...", len(socketConnections[pairConn]))
	var msg any
	for {
		err := socket.ReadJSON(&msg)
		if err != nil {
			fmt.Println("socket reading error........")
			log.Println(err)
			break
		}
	}
}

func (service *orderBookService) WriteConnectedSocketToOrderBookUpdates(message []byte) {
	connection := socketConnections["btc"]
	for k, v := range connection {
		fmt.Println("loop is going", k, v)
		err := v.WriteMessage(1, message)
		if err != nil {
			fmt.Println("error is -------", err)
			delete(connection, k)
		}
	}
}

func (service *orderBookService) GetOrderBook(ctx context.Context, orderType string) []order.Order {
	return service.repository.GetOrderBook(ctx, orderType)
}
func (service *orderBookService) AddToOrderBook(ctx context.Context, orderKey string, orderId int, value string) error {
	return service.repository.AddToOrderBook(ctx, orderKey, orderId, value)
}
func (service *orderBookService) RemoveFromOrderBook(ctx context.Context, orderType string, orderId int) error {
	return service.repository.RemoveFromOrderBook(ctx, orderType, orderId)
}
func (service *orderBookService) UpdateOrderBook(ctx context.Context, orderType string, orderId int, updatedOrder string) error {
	return service.repository.UpdateOrderBook(ctx, orderType, orderId, updatedOrder)
}
