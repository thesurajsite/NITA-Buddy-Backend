package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/suraj/nitabuddy/database"
	"github.com/suraj/nitabuddy/handlers"
	"github.com/suraj/nitabuddy/models"
	"github.com/suraj/nitabuddy/routes"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found (likely running in production):", err)
	}

	// Connect to MongoDB
	client, userCollection, orderCollection, rewardsCollection := database.Connect() // returns collection reference
	defer client.Disconnect(context.Background())

	// Create Models
	rewardsModel := models.NewRewardsModel(rewardsCollection)
	userModel := models.NewUserModel(userCollection, rewardsModel)
	orderModel := models.NewOrderModel(orderCollection, userCollection, rewardsModel)

	// Define your JWT secret key (keep it safe and strong)
	jwtSecret := []byte("my-secret-key")

	// Create handlers with JWT-based auth
	authHandler := handlers.NewAuthHandler(userModel, jwtSecret)
	orderHandler := handlers.NewOrderHandler(orderModel, authHandler)       // Pass authHandler
	rewardsHandler := handlers.NewRewardsHandler(rewardsModel, authHandler) // Pass authHandler

	// configure router
	r := mux.NewRouter()
	routes.Setup(r, authHandler, orderHandler, rewardsHandler)

	// Start server
	log.Println("Server starting at port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
