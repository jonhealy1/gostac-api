package database

import (
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func StartConsumer(topic string, handleMsg func(msg *kafka.Message)) {
	retryInterval := 10 * time.Second
	maxRetries := 10

	for {
		consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
			"bootstrap.servers": "kafka:9092",
			"group.id":          "es-api-group",
			"auto.offset.reset": "earliest",
			"socket.timeout.ms": 60000, // Increase the timeout to 60 seconds
		})

		if err != nil {
			log.Printf("Failed to create consumer: %s", err)
			time.Sleep(retryInterval)
			continue
		}

		consumer.Subscribe(topic, nil)

		for retries := 0; retries < maxRetries; {
			ev := consumer.Poll(100) // Changed from Poll(0) to Poll(100)
			if ev == nil {
				retries++
				time.Sleep(retryInterval)
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				handleMsg(e) // Removed 'topic' from the function call
				retries = 0
			case kafka.Error:
				log.Printf("Error in Kafka consumer: %v\n", e)
			default:
				// Do nothing
			}
		}

		consumer.Close()
		log.Printf("Consumer closed after %d retries, trying to restart", maxRetries)
	}
}
