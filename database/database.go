package database

import (
	"fmt"
	"go-stac-api-postgres/models"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Dbinstance struct {
	Db *gorm.DB
}

var DB Dbinstance

func getEnv(key string) (string, error) {
	err := godotenv.Load()
	return os.Getenv(key), err
}

func ConnectDb() {
	host, port, user, pass, dbname := "", "", "", "", ""
	host, err := getEnv("POSTGRES_HOST")

	// this is done for CI, not ideal ....
	if err != nil {
		host = "database"
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
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, pass, dbname, "disable",
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
