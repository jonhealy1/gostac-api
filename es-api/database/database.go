package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/olivere/elastic/v7"
)

type ESInstance struct {
	Client *elastic.Client
}

var ES ESInstance

func getEnv(key string) (string, error) {
	err := godotenv.Load()
	return os.Getenv(key), err
}

func getEnvWithDefault(key, defaultValue string) string {
	value, err := getEnv(key)
	if err != nil || value == "" {
		return defaultValue
	}
	return value
}

func ConnectES() {
	host := getEnvWithDefault("ES_HOST", "localhost")
	port := getEnvWithDefault("ES_PORT", "9200")
	user := getEnvWithDefault("ES_USER", "username")
	pass := getEnvWithDefault("ES_PASS", "password")

	log.Println("host: ", host)
	log.Println("port: ", port)

	dsn := fmt.Sprintf(
		"http://%s:%s",
		host, port,
	)

	log.Println("dsn: ", dsn)

	// es, err := elastic.NewClient(
	// 	elastic.SetURL(dsn),
	// 	elastic.SetSniff(false),
	// 	elastic.SetBasicAuth(user, pass), // Add this line to set basic authentication
	// )

	es, err := elastic.NewClient(
		elastic.SetURL(dsn),
		elastic.SetSniff(false),
		elastic.SetBasicAuth(user, pass),
		elastic.SetHealthcheckTimeoutStartup(10*time.Second), // Increase the timeout
	)

	if err != nil {
		log.Fatal("Failed to connect to elastic search. \n", err)
	}

	log.Println("connected to elastic search")
	ES = ESInstance{
		Client: es,
	}
	createCollectionsIndex(ES)
	createItemsIndex(ES)
}

func createCollectionsIndex(database ESInstance) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create the collections index
	indexName := "collections"
	exists, err := database.Client.IndexExists(indexName).Do(ctx)
	if err != nil {
		log.Fatalf("Could not contact Elasticsearch: %v", err)
	}
	if !exists {
		// Define the mapping for the index
		mapping := `{
			"mappings": {
				"properties": {
				"data": {
					"properties": {
					"extent": {
						"properties": {
						"temporal": {
							"properties": {
							"interval": {
								"type": "text"
							}
							}
						}
						}
					}
					}
				}
				}
			}
		}`

		_, err := database.Client.CreateIndex(indexName).BodyString(mapping).Do(ctx)
		if err != nil {
			log.Fatalf("Could not create Elasticsearch index: %v", err)
		}
	}

	// Create other indices as needed
}

func createItemsIndex(database ESInstance) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create the items index
	indexName := "items"
	exists, err := database.Client.IndexExists(indexName).Do(ctx)
	if err != nil {
		log.Fatalf("Could not contact Elasticsearch: %v", err)
	}
	if !exists {
		// Define the mapping for the index
		mapping := `{
            "mappings": {
                "properties": {
                    "geometry": {
                        "type": "geo_shape"
                	},
					"collection": {
						"type": "keyword"
					},
					"properties": {
						"properties": {
							"datetime": {
								"type": "date"
							}
						}
					}
            	}
			}
        }`

		_, err := database.Client.CreateIndex(indexName).BodyString(mapping).Do(ctx)
		if err != nil {
			log.Fatalf("Could not create Elasticsearch index: %v", err)
		}
	}

	// Create other indices as needed
}
