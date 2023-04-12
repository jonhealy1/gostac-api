package controllers

import (
	"encoding/json"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gofiber/fiber/v2"
	"github.com/jonhealy1/goapi-stac/es-api/models"
)

func CreateCollectionFromMessage(msg *kafka.Message) {
	// Parse the message to get the collection
	var stac_collection models.StacCollection
	err := json.Unmarshal(msg.Value, &stac_collection)
	if err != nil {
		log.Printf("Error unmarshalling collection from Kafka message: %v\n", err)
		return
	}

	// Call CreateESCollection with the parsed collection
	c := &fiber.Ctx{} // Create an empty context
	err = CreateESCollection(c, &stac_collection)
	if err != nil {
		log.Printf("Error creating ES collection from Kafka message: %v\n", err)
	}
}
