package app_kafka

import (
	"context"
	"fmt"
	"log"
	"os"
	"vcs_server/helper"

	"github.com/segmentio/kafka-go"
)

type ProducerKafka interface {
	ProduceMessage(id string) error
}

type producerKakfa struct {
	writer *kafka.Writer
}

var Producer ProducerKafka

func init() {

	Producer = &producerKakfa{
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers:  []string{os.Getenv("KAFKA_BROKER")},
			Topic:    os.Getenv("KAFKA_TOPIC"),
			Balancer: &kafka.LeastBytes{},
		}),
	}
}

func (producer *producerKakfa) ProduceMessage(id string) error {
	random, _ := helper.RandString(10)
	message := kafka.Message{
		Key:   []byte(id + random),
		Value: []byte(id),
	}
	fmt.Println("Producer to Kafka", id)
	err := producer.writer.WriteMessages(context.Background(), message)
	if err != nil {
		log.Println(
			"Error while writing message to Kafka",
		)
		return err
	}
	return nil

}
