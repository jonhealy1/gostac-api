package controllers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/jonhealy1/goapi-stac/es-api/models"
)

func CreateCollectionFromMessage(msg *kafka.Message) {
	// Parse the message to get the collection
	var stac_collection models.StacCollection
	fmt.Println("Message: ", string(msg.Value))
	err := json.Unmarshal(msg.Value, &stac_collection)
	if err != nil {
		log.Printf("Error unmarshalling collection from Kafka message: %v\n", err)
		return
	}

	// Call CreateESCollectionCore with the parsed collection
	collection, err := UpsertESCollectionCore(false, stac_collection.Id, &stac_collection)
	if err != nil {
		log.Printf("Error creating ES collection from Kafka message: %v\n", err)
	} else {
		log.Printf("ES collection created successfully: %s\n", collection.Id)
	}
}

func UpdateCollectionFromMessage(msg *kafka.Message) {
	fmt.Println("UpdateCollectionFromMessage")
	// Parse the message to get the collection
	var stac_collection models.StacCollection
	fmt.Println("Message: ", string(msg.Value))
	err := json.Unmarshal(msg.Value, &stac_collection)
	if err != nil {
		log.Printf("Error1 unmarshalling collection from Kafka message: %v\n", err)
		return
	}

	// Call UpdateESCollectionCore with the parsed collection
	collection, err := UpsertESCollectionCore(true, stac_collection.Id, &stac_collection)
	if err != nil {
		log.Printf("Error updating ES collection from Kafka message: %v\n", err)
	} else {
		log.Printf("ES collection updated successfully: %s\n", collection.Id)
	}
}
