package main

import (
	"log"
	"wasa/service/api/routes"
	database2 "wasa/service/database"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Connect to the database
	database2.InitDB()
	database2.MigrateDB()

	// Create the server
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8081"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Requested-With"},
		AllowCredentials: true, // Allow using cookies
	}))

	// Set up routes
	routes.SetupRoutes(r)

	// Start the server
	log.Println("Server is running on port 8080")
	r.Run(":8080")
}
