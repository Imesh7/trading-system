package ports

type kafkaService interface {
	OrderMatchConsumer()
	OrderMatchProducer()
}


