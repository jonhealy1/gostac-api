package database

import (
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func StartConsumer(topic string, handleMsg func(*kafka.Message)) {
	retryInterval := 10 * time.Second
	maxRetries := 5

	for {
		consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
			"bootstrap.servers": "kafka:9092",
			"group.id":          "es-api-group",
			"auto.offset.reset": "earliest",
		})

		if err != nil {
			log.Printf("Failed to create consumer: %s", err)
			time.Sleep(retryInterval)
			continue
		}

		consumer.Subscribe(topic, nil)

		for retries := 0; retries < maxRetries; {
			msg, err := consumer.ReadMessage(-1)
			if err != nil {
				log.Printf("Error while reading message: %s", err)
				retries++
				time.Sleep(retryInterval)
				continue
			}

			handleMsg(msg)
			retries = 0
		}

		consumer.Close()
		log.Printf("Consumer closed after %d retries, trying to restart", maxRetries)
	}
}
