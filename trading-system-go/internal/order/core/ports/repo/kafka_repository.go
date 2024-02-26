package ports

type Kafka interface {
	OrderMatchProducer(topic string, orderId int)
	OrderMatchConsumer(topic string, handler func([]byte))
}
