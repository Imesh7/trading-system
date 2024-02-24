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
func (service *orderBookKafkaService) CreateBidConsumer(topic string/* , handler func([]byte) */) error {
	service.adapter.CreateBidConsumer(topic, /* handler */)
	/* con := routes.Conn["btc"]

	jsonData, err := json.Marshal(handler.Value)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(msg.Value)
	for i, v := range con {
		err := v.WriteMessage(1, msg.Value)
		if err != nil {
			log.Print(err)
			con = removeElement(con, i)
		}
	} */

	return nil
}
