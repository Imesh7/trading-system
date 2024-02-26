package service

import (
	order "trading-system-go/internal/order/core/domain"
	"trading-system-go/internal/order_book/adapter"
)

type orderBookKafkaService struct {
	adapter *adapter.KafkaAdapter
}

func NewOrderBookKakfkaService(adapter *adapter.KafkaAdapter) *orderBookKafkaService {
	return &orderBookKafkaService{
		adapter: adapter,
	}
}

func (service *orderBookKafkaService) CreateOrderBookProducer(topic string, order *order.Order) error {
	return service.adapter.CreateOrderBookProducer(topic, order)
}
func (service *orderBookKafkaService) CreateBidConsumer(topic string) error {
	service.adapter.CreateBidConsumer(topic)
	return nil
}
