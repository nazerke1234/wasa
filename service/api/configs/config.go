package configs

import (
	"log"

	"github.com/joho/godotenv"
)

// InitConfig loads environment variables
func InitConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error downloading .env file")
	}
}
