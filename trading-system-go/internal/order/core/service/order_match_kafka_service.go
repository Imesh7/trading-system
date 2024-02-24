package service

import (
	//"encoding/binary"
	"trading-system-go/internal/order/adapter"
)

type OrderMatchKafkaService struct {
	adapter           *adapter.KafkaAdapter
	orderMatchService *orderMatchService
}

func NewKafkaService(adapter *adapter.KafkaAdapter, orderMatchService *orderMatchService) *OrderMatchKafkaService {
	return &OrderMatchKafkaService{
		adapter:           adapter,
		orderMatchService: orderMatchService,
	}
}

func (service OrderMatchKafkaService) OrderMatchProducer(topic string, orderId int) {
	service.adapter.OrderMatchProducer(topic, orderId)
}

func (service OrderMatchKafkaService) OrderMatchConsumer(topic string) {
	service.adapter.OrderMatchConsumer(topic)
}
