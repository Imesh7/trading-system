package adapter

import (
	"trading-system-go/internal/order_book/core/ports"

	"github.com/gofiber/contrib/websocket"
)

type OrderBookHandler struct {
	service ports.OrderBookService
}

func NewOrderBookHandler(service ports.OrderBookService) *OrderBookHandler {
	return &OrderBookHandler{
		service: service,
	}
}

func (hanlder OrderBookHandler) ConnectWithOrderBook(socket *websocket.Conn) {
	hanlder.service.ConnectWithOrderBook(socket)
}
