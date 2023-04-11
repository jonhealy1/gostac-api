package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jonhealy1/goapi-stac/pg-api/models"

	"github.com/joho/godotenv"
	"github.com/olivere/elastic/v7"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Dbinstance struct {
	Db *gorm.DB
}

type ESInstance struct {
	Client *elastic.Client
}

var DB Dbinstance
var ES ESInstance

func getEnv(key string) (string, error) {
	err := godotenv.Load()
	return os.Getenv(key), err
}

func ConnectDb() {
	host, port, user, pass, dbname := "", "", "", "", ""
	host, err := getEnv("POSTGRES_HOST")

	// this is done for CI, not ideal ....
	if err != nil {
		host = "localhost"
		port = "5433"
		user = "username"
		pass = "password"
		dbname = "postgis"
	} else {
		port, _ = getEnv("POSTGRES_PORT")
		user, _ = getEnv("POSTGRES_USER")
		pass, _ = getEnv("POSTGRES_PASS")
		dbname, _ = getEnv("POSTGRES_DBNAME")
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s",
		host, port, user, pass, dbname,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
	}

	log.Println("connected")
	db.Logger = logger.Default.LogMode(logger.Info)
	log.Println("running migrations")
	db.AutoMigrate(&models.Collection{})

	db.Exec(`CREATE TABLE IF NOT EXISTS items (
		id TEXT PRIMARY KEY NOT NULL,
		collection TEXT,
		data JSONB,
		geometry geometry(POLYGON, 4326) NOT NULL
	);`)

	DB = Dbinstance{
		Db: db,
	}
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

// func ConnectRedis() {
// 	host, port, pass := "", "", ""
// 	host, err := getEnv("REDIS_HOST")

// 	// this is done for CI, not ideal ....
// 	if err != nil {
// 		host = "localhost"
// 		port = "6379"
// 		pass = ""
// 	} else {
// 		port, _ = getEnv("REDIS_PORT")
// 		pass, _ = getEnv("REDIS_PASS")
// 	}

// 	dsn := fmt.Sprintf(
// 		"%s:%s",
// 		host, port,
// 	)

// 	client := redis.NewClient(&redis.Options{
// 		Addr:     dsn,
// 		Password: pass, // no password set
// 		DB:       0,    // use default DB
// 	})

// 	_, err = client.Ping().Result()
// 	if err != nil {
// 		log.Fatal("Failed to connect to redis. \n", err)
// 	}

// 	log.Println("connected to redis")
// 	Redis = client
// }
