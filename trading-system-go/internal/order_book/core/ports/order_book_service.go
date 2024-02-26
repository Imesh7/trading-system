package ports

import (
	"context"
	order "trading-system-go/internal/order/core/domain"

	"github.com/gofiber/contrib/websocket"
)

type OrderBookService interface {
	ConnectWithOrderBook(socket *websocket.Conn)
	WriteConnectedSocketToOrderBookUpdates(message []byte)
	GetOrderBook(ctx context.Context, orderType string) []order.Order
	AddToOrderBook(ctx context.Context, orderKey string, orderId int, value string) error
	RemoveFromOrderBook(ctx context.Context, orderType string, orderId int) error
	UpdateOrderBook(ctx context.Context, orderType string, orderId int, updatedOrder string) error
}
