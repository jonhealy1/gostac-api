package database

import (
	"fmt"
	"log"
	"os"

	"go-stac-api-postgres/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/joho/godotenv"
)

type Dbinstance struct {
	Db *gorm.DB
}

var DB Dbinstance

func getEnv(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

func ConnectDb() {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		getEnv("POSTGRES_HOST"), getEnv("POSTGRES_PORT"), getEnv("POSTGRES_USER"),
		getEnv("POSTGRES_PASS"), getEnv("POSTGRES_DBNAME"), "disable",
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
	db.AutoMigrate(&models.Item{})

	DB = Dbinstance{
		Db: db,
	}
}
