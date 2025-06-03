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
		log.Fatal("Error loading .env file")
	}

	// Connect to MaongoDB
	client, userCollection := database.Connect()
	defer client.Disconnect(context.Background())

	// Create Models
	userModel := models.NewUserModel(userCollection)

	// Define your JWT secret key (keep it safe and strong)
	jwtSecret := []byte("your-secret-key")

	// Create handlers with JWT-based auth
	authHandler := handlers.NewAuthHandler(userModel, jwtSecret)

	// configure router
	r := mux.NewRouter()
	routes.Setup(r, authHandler)

	// Start server
	log.Println("Server starting at port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
