package database

import (
	"fmt"
	"log"
	"os"

	"github.com/jonhealy1/goapi-stac/models"

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

func ConnectES() {
	host, port, user, pass := "", "", "", ""
	host, err := getEnv("ES_HOST")

	// this is done for CI, not ideal ....
	if err != nil {
		host = "localhost"
		port = "9200"
		user = ""
		pass = ""
	} else {
		port, _ = getEnv("ES_PORT")
		user, _ = getEnv("ES_USER")
		pass, _ = getEnv("ES_PASS")
	}

	dsn := fmt.Sprintf(
		"http://%s:%s",
		host, port,
	)

	es, err := elastic.NewClient(
		elastic.SetURL(dsn),
		elastic.SetSniff(false),
		elastic.SetBasicAuth(user, pass), // Add this line to set basic authentication
	)

	if err != nil {
		log.Fatal("Failed to connect to elastic search. \n", err)
	}

	log.Println("connected to elastic search")
	ES = ESInstance{
		Client: es,
	}
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
