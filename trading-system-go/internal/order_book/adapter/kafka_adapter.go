package adapter

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
	order "trading-system-go/internal/order/core/domain"
	"trading-system-go/internal/order_book/core/ports"

	"github.com/IBM/sarama"
	"github.com/gofiber/contrib/websocket"
)

type KafkaAdapter struct {
	producer         sarama.SyncProducer
	consumerGroup    sarama.ConsumerGroup
	orderBookService ports.OrderBookService
}

func NewOrderBookKafkaAdpater(orderBookService ports.OrderBookService) *KafkaAdapter {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	kafkaHost := fmt.Sprintf("%s:%s", os.Getenv("KAFKA_HOST"), os.Getenv("KAFKA_PORT"))
	producer, err := sarama.NewSyncProducer([]string{kafkaHost}, config)
	if err != nil {
		fmt.Println("Cannot connect to producer")
		log.Fatal(err)
	}

	config.ClientID = "go-kafka-consumer"
	consumerGroup, err := sarama.NewConsumerGroup([]string{kafkaHost}, "2", config)
	if err != nil {
		fmt.Fprintln(os.Stdout, []any{"Errors is %s", err}...)
		log.Fatal(err)
	}
	return &KafkaAdapter{
		producer:      producer,
		consumerGroup: consumerGroup,
		orderBookService: orderBookService,
	}
}

func (kafkaAdapter KafkaAdapter) CreateOrderBookProducer(topic string, order interface{}) error {

	orderJson, err := json.Marshal(order)

	message := &sarama.ProducerMessage{
		Topic:     topic,
		Key:       nil,
		Value:     sarama.StringEncoder(orderJson),
		Headers:   []sarama.RecordHeader{},
		Metadata:  nil,
		Offset:    0,
		Partition: 0,
		Timestamp: time.Time{},
	}

	_, _, err = kafkaAdapter.producer.SendMessage(message)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func (kafkaAdapter KafkaAdapter) CreateBidConsumer(topic string /* , handler func([]byte) */) error {

	defer func() {
		if err := kafkaAdapter.consumerGroup.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	ctx := context.Background()
	for {
		kafkaAdapter.consumerGroup.Consume(ctx, []string{topic}, &kafkaAdapter)
	}

	/* for {
		select {
		case msg := <-kafkaAdapter.consumerGroup.Consume():
			//valueString := string(msg.Value)
			//fmt.Fprintln(os.Stdout, []any{"Received bid for consumer:", valueString}...)
			//match_order.MatchOrder(valueString)


		case <-signals:
			return
		}
	} */
}

func removeElement(slice []*websocket.Conn, index int) []*websocket.Conn {
	if index < 0 || index >= len(slice) {
		fmt.Println("Invalid index")
		return slice
	}

	return append(slice[:index], slice[index+1:]...)
}

func (ka *KafkaAdapter) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (ka *KafkaAdapter) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (ka *KafkaAdapter) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		fmt.Println("process---------", message.Value)
		var order order.Order
		err := json.Unmarshal(message.Value, &order)
		if err != nil {
			panic(err)
		}
		ka.processMessage(order)
		go ka.orderBookService.WriteConnectedSocketToOrderBookUpdates(message.Value)
	}
	return nil
}

func (ka *KafkaAdapter) processMessage(order order.Order) {
	// Process Kafka message
	// You can call the provided handler function here

	fmt.Println("process message order book got the message--------------------", order)

}
