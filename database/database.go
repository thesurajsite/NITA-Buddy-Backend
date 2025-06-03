package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// connect() establishes a connection to MongoDB  and returns the client and collections
func Connect() (*mongo.Client, *mongo.Collection) {

	// connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// MongoDB connection string
	connectionString := os.Getenv("MONGO_URI")
	log.Println("MONGO_URI:", connectionString)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatal(err)
	}

	// Initialize Collections
	db := client.Database("nita_buddy")
	userCollection := db.Collection("users")

	// check connection by running a query
	err = userCollection.FindOne(ctx, bson.M{}).Err()
	if err != nil && err != mongo.ErrNoDocuments {
		log.Fatalf("Failer to query user Collection: %v", err)
	}

	log.Println("Successfully connected to NITA Buddy Database")
	return client, userCollection
}
