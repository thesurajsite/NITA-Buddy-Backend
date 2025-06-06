package models

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Order struct {
	OrderID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	CustomOrderID string             `bson:"custom_order_id" json:"custom_order_id"`
	Store         string             `bson:"store" json:"store"`
	OrderDetails  string             `bson:"order_details" json:"order_details"`
	Status        string             `bson:"status" json:"status"`
	OTP           string             `bson:"otp" json:"otp"`
	PlacedBy      primitive.ObjectID `bson:"placed_by" json:"placed_by"`
	AcceptedBy    primitive.ObjectID `bson:"accepted_by" json:"accepted_by"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
}

type OrderModel struct {
	collection *mongo.Collection
}

func NewOrderModel(collection *mongo.Collection) *OrderModel {
	return &OrderModel{collection: collection}
}

func (m *OrderModel) CreateOrder(input Order) (*Order, error) {

	customID, err := generateCustomOrderID(m.collection)
	if err != nil {
		return nil, err
	}

	order := &Order{
		CustomOrderID: customID,
		Store:         input.Store,
		OrderDetails:  input.OrderDetails,
		Status:        input.Status,
		OTP:           input.OTP,
		PlacedBy:      input.PlacedBy,
		AcceptedBy:    input.AcceptedBy,
		CreatedAt:     time.Now(),
	}

	result, err := m.collection.InsertOne(context.Background(), order)
	if err != nil {
		return nil, err
	}

	order.OrderID = result.InsertedID.(primitive.ObjectID)
	return order, nil
}

func generateCustomOrderID(collection *mongo.Collection) (string, error) {
	rand.Seed(time.Now().UnixNano())
	maxAttempts := 5

	for i := 0; i < maxAttempts; i++ {
		suffix := rand.Intn(90000) + 10000
		id := fmt.Sprintf("#NBO%d", suffix)

		count, err := collection.CountDocuments(context.Background(), bson.M{"custom_order_id": id})
		if err != nil {
			return "", err
		}
		if count == 0 {
			return id, nil // unique id found
		}
	}

	return "", fmt.Errorf("failed to generate unique custom order id")
}
