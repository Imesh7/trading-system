package ports

import order "trading-system-go/internal/order/core/domain"

type OrderService interface {
	CreateOrderService(order *order.Order) (*order.Order, error)
	GetOrdersService() (*[]order.Order, error)
}
