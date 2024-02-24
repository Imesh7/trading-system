package ports

import order "trading-system-go/internal/order/core/domain"

type OrderRepository interface {
	CreateOrder(order *order.Order) (*order.Order, error)
	GetOrders() (*[]order.Order, error)
	FindOrder(orderId int) (*order.Order, error)
}
