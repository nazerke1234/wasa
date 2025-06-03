package database

import (
	"fmt"
	"log"
	"os"
	"wasa/service/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB — global variable for database access
var DB *gorm.DB

// Initialize the database
func InitDB() {
	// Load environment variables
	if err := godotenv.Load("/app/.env"); err != nil {
		log.Println("⚠.env file not found, default values are being used")
	}

	// Build the connection string
	dsn := fmt.Sprintf(
		"host=db user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)

	// Connect to the database
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to the database:", err)
	}

	// Perform database migrations
	err = DB.AutoMigrate(&models.User{}, &models.Message{})
	if err != nil {
		log.Fatal("Error migrating the database:", err)
	}

	log.Println("Database connected!")
}
