package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func main() {
	// Set Gin mode to release
	gin.SetMode(gin.ReleaseMode)

	var MONGO_URL = "<your_connection_string>"

	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI(MONGO_URL)
	var err error
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	defer client.Disconnect(context.Background())

	// Create a new router
	router := gin.Default()

	// Task 7 solution starts
	// Task 7 solution ends

	// Start the server
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
