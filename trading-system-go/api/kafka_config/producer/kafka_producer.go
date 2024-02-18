package order_kafka_producer

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/IBM/sarama"
)

func OrderMatchProducer(topic string, orderId int) {
	con := sarama.NewConfig()
	con.Producer.Return.Successes = true
	kafkaHost := fmt.Sprintf("%s:%s", os.Getenv("KAFKA_HOST"), os.Getenv("KAFKA_PORT"))
	producer, err := sarama.NewSyncProducer([]string{kafkaHost}, con)
	if err != nil {
		fmt.Println("Cannot connect to producer")
		log.Fatal(err)
	}

	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatal(err)
		}
	}()

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

	partition, offset, err := producer.SendMessage(message)
	if err != nil {
		fmt.Fprintln(os.Stdout, []any{"message passed...................2 %s", err}...)
		log.Fatal(err)
	}
	log.Printf("Produced message to partition %d at offset %d\n", partition, offset)
}

func CreateOrderBookProducer(topic string, order interface{}) {
	con := sarama.NewConfig()
	con.Producer.Return.Successes = true
	kafkaHost := fmt.Sprintf("%s:%s", os.Getenv("KAFKA_HOST"), os.Getenv("KAFKA_PORT"))
	producer, err := sarama.NewSyncProducer([]string{kafkaHost}, con)
	if err != nil {
		fmt.Println("Cannot connect to producer")
		log.Fatal(err)
	}

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

	_, _, err = producer.SendMessage(message)
	if err != nil {
		log.Fatal(err)
	}

}
