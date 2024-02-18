package order_kafka_consumer

import (
	//"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"trading-system-go/api/match_order"
	routes "trading-system-go/internal/route"

	"github.com/IBM/sarama"
	"github.com/gofiber/contrib/websocket"
)

func OrderMatchConsumer(topic string) {
	config := sarama.NewConfig()
	config.ClientID = "go-kafka-consumer"
	kafkaHost := fmt.Sprintf("%s:%s", os.Getenv("KAFKA_HOST"), os.Getenv("KAFKA_PORT"))
	consumer, err := sarama.NewConsumer([]string{kafkaHost}, config)
	if err != nil {
		fmt.Fprintln(os.Stdout, []any{"Errors is %s", err}...)
		log.Fatal(err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			valueString := string(msg.Value)
			fmt.Fprintln(os.Stdout, []any{"Received integer value: %s", valueString}...)
			match_order.MatchOrder(valueString)
		case <-signals:
			return
		}
	}
}

func CreateBidConsumer(topic string) {
	config := sarama.NewConfig()
	config.ClientID = "go-kafka-consumer"
	kafkaHost := fmt.Sprintf("%s:%s", os.Getenv("KAFKA_HOST"), os.Getenv("KAFKA_PORT"))
	consumer, err := sarama.NewConsumer([]string{kafkaHost}, config)
	if err != nil {
		fmt.Fprintln(os.Stdout, []any{"Errors is %s", err}...)
		log.Fatal(err)
	}
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	partitionConsumer, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			//valueString := string(msg.Value)
			//fmt.Fprintln(os.Stdout, []any{"Received bid for consumer:", valueString}...)
			//match_order.MatchOrder(valueString)
			con := routes.Conn["btc"]
			//jsonData, err := json.Marshal(msg.Value)
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
			}
		case <-signals:
			return
		}
	}
}

func removeElement(slice []*websocket.Conn, index int) []*websocket.Conn {
	if index < 0 || index >= len(slice) {
		fmt.Println("Invalid index")
		return slice
	}

	return append(slice[:index], slice[index+1:]...)
}
