package database

import (
	"fmt"
	"go-stac-api-postgres/models"
	"log"
	"os"

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
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}

// ConnectDb creates a database connection and sets the DB global variable.
// The function gets the database connection parameters from environment variables POSTGRES_HOST, POSTGRES_PORT, POSTGRES_USER, POSTGRES_PASS, and POSTGRES_DBNAME.
// It uses the gorm library to create the database connection and configures it to log database queries in the "info" mode.
// The function also runs database migrations and creates the "items" table if it does not exist.
func ConnectDb() {
	host := getEnv("POSTGRES_HOST")
	port := getEnv("POSTGRES_PORT")
	user := getEnv("POSTGRES_USER")
	pass := getEnv("POSTGRES_PASS")
	dbname := getEnv("POSTGRES_DBNAME")

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
