package adapter

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	service "trading-system-go/internal/order/core/ports/service"

	"github.com/IBM/sarama"
)

type KafkaAdapter struct {
	producer          sarama.SyncProducer
	consumerGroup     sarama.ConsumerGroup
	orderMatchService service.OrderMatchService
}


func NewKafkaAdpater(orderMatchService service.OrderMatchService) *KafkaAdapter {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	kafkaHost := fmt.Sprintf("%s:%s", os.Getenv("KAFKA_HOST"), os.Getenv("KAFKA_PORT"))
	producer, err := sarama.NewSyncProducer([]string{kafkaHost}, config)
	if err != nil {
		fmt.Println("Cannot connect to producer")
		log.Fatal(err)
	}

	config.ClientID = "go-kafka-consumer"
	consumerGroup, err := sarama.NewConsumerGroup([]string{kafkaHost}, "1", config)
	if err != nil {
		fmt.Fprintln(os.Stdout, []any{"Errors is %s", err}...)
		log.Fatal(err)
	}

	return &KafkaAdapter{
		producer:          producer,
		consumerGroup:     consumerGroup,
		orderMatchService: orderMatchService,
	}
}

func (adapter KafkaAdapter) OrderMatchProducer(topic string, orderId int) {

	/* defer func() {
		if err := adapter.producer.Close(); err != nil {
			log.Fatal(err)
		}
	}() */

	message := &sarama.ProducerMessage{
		Topic:     topic,
		Key:       nil,
		Value:     sarama.StringEncoder(fmt.Sprint(orderId)),
		Headers:   []sarama.RecordHeader{},
		Metadata:  nil,
		Offset:    0,
		Partition: 0,
		Timestamp: time.Time{},
	}

	partition, offset, err := adapter.producer.SendMessage(message)
	if err != nil {
		fmt.Fprintln(os.Stdout, []any{"message passed...................2 %s", err}...)
		log.Fatal(err)
	}
	log.Printf("Produced message to partition %d at offset %d\n", partition, offset)
}

func (adapter KafkaAdapter) OrderMatchConsumer(topic string) {
	defer func() {
		if err := adapter.consumerGroup.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	ctx := context.Background()
	for {
		adapter.consumerGroup.Consume(ctx, []string{topic}, &adapter)
	}

	/* partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt) */

	/* for {
		select {
		case msg := <-partitionConsumer.Messages():
			valueString := string(msg.Value)
			fmt.Fprintln(os.Stdout, []any{"Received integer value: %s", valueString}...)
			match_order.MatchOrder(int64(binary.BigEndian.Uint64(msg.Value)))
		case <-signals:
			return
		}
	} */
}

func (ka *KafkaAdapter) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (ka *KafkaAdapter) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (ka *KafkaAdapter) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		ka.processMessage(message.Value)
		intValue, err := strconv.Atoi(string(message.Value))
		if err != nil {
			panic(err)
		}
		go ka.orderMatchService.MatchOrder(intValue)
	}
	return nil
}

func (ka *KafkaAdapter) processMessage(message []byte) {
	intValue, err := strconv.Atoi(string(message))
	if err != nil {
		panic(err)
	}
	fmt.Println("process message order got the message--------------------", intValue)
}
