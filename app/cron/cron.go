package cron

import (
	"context"
	"fmt"
	"log"
	"os"
	"vcs_server/cache"
	"vcs_server/healcheck"
	"vcs_server/mail"

	"github.com/segmentio/kafka-go"
)

func Cron_healcheck() {
	// Implement your cron job here
	fmt.Printf("Starting healcheck cron")
	healcheck.SendHealcheck()
}

func Cron_SendReport() {
	fmt.Printf("Starting send_report cron")
	mail.SendMail(12)

}

func Cron_Invalid_Cache() {
	// Implement your cron job here
	fmt.Printf("Starting invalid cache cron")
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{os.Getenv("KAFKA_BROKER")}, // Địa chỉ Kafka broker
		Topic:       os.Getenv("KAFKA_TOPIC"),
		Partition:   0,
		MinBytes:    10e3, // 10KB
		MaxBytes:    10e6, // 10MB
		StartOffset: kafka.LastOffset,
	})
	serverCache := cache.NewRedisCache("redis:6379", 0, 120)
	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatalf("could not read message %v", err)
		}
		log.Println("received:", string(msg.Value))
		serverCache.Del(string(msg.Value))
	}
}
