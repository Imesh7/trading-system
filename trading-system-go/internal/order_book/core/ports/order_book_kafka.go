package ports

import order "trading-system-go/internal/order/core/domain"

type OrderBookKafka interface {
	CreateOrderBookProducer(topic string, order *order.Order) error
	CreateBidConsumer(topic string/* , handler func([]byte) */) error
}
