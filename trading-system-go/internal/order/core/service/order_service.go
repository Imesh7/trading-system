package service

import (
	order "trading-system-go/internal/order/core/domain"
	ports "trading-system-go/internal/order/core/ports/repo"
)

type orderService struct {
	repo                  ports.OrderRepository
	kafkaService          *OrderMatchKafkaService
	
}

func NewOrderService(repo ports.OrderRepository, kafkaService *OrderMatchKafkaService) *orderService {
	return &orderService{
		repo:                  repo,
		kafkaService:          kafkaService,
	}
}

func (orderService *orderService) CreateOrderService(orderData *order.Order) (*order.Order, error) {
	orderData.OrderStatus = order.NotFilled
	order, err := orderService.repo.CreateOrder(orderData)
	if err != nil {
		return nil, err
	}

	go orderService.kafkaService.OrderMatchProducer("topic", orderData.OrderId)

	return order, nil
}

func (orderService *orderService) GetOrdersService() (*[]order.Order, error) {
	orderList, err := orderService.repo.GetOrders()
	if err != nil {
		return nil, err
	}
	return orderList, nil
}
