package main

import (
	"context"
	"lock/config"
	"lock/database"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading the config file: %v", err)
	}

	// Connect to MongoDB
	db, err := database.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	// Close the MongoDB client connection when main function exits
	defer db.Client().Disconnect(context.Background())

	// Create a new Gin router
	router := gin.Default()

	// Pass the MongoDB client to route handler

	// Start the server
	err = router.Run("localhost:8080")
	if err != nil {
		log.Fatalf("Local host error: %v", err)
	}
}
