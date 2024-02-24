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

func NewOrderBookService(repository ports.OrderBookRepository) *orderBookService {
	return &orderBookService{
		repository: repository,
	}
}

var socketConnections map[string][]*websocket.Conn

func init() {
	socketConnections = make(map[string][]*websocket.Conn)
}

func (service *orderBookService) ConnectWithOrderBook(socket *websocket.Conn) {
	pairConn := socket.Query("pair")

	socketConnections[pairConn] = append(socketConnections[pairConn], socket)
	defer socket.Close()
	fmt.Println("received websocket conn req.........................", len(socketConnections[pairConn]))
	var msg any
	for {
		err := socket.ReadJSON(&msg)
		if err != nil {
			// optional: log the error
			break
		}
	}
}

func (service *orderBookService) WriteConnectedSocketToOrderBookUpdates(message []byte) {
	con := socketConnections["btc"]

	for i, v := range con {
		err := v.WriteMessage(1, message)
		if err != nil {
			fmt.Println("error is -------", err)
			log.Print(err)
			con = removeElement(con, i)
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

func removeElement(slice []*websocket.Conn, index int) []*websocket.Conn {
	if index < 0 || index >= len(slice) {
		fmt.Println("Invalid index")
		return slice
	}

	return append(slice[:index], slice[index+1:]...)
}
