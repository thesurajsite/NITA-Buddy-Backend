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

func (m *OrderModel) CreateOrder(store, orderDetails string, placedBy primitive.ObjectID) (*Order, error) {

	status := "NotAccepted"
	otp := generateOTP()
	acceptedBy := primitive.NilObjectID
	customID, err := generateCustomOrderID(m.collection)
	if err != nil {
		return nil, err
	}

	order := &Order{
		CustomOrderID: customID,
		Store:         store,
		OrderDetails:  orderDetails,
		Status:        status,
		OTP:           otp,
		PlacedBy:      placedBy,
		AcceptedBy:    acceptedBy,
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

func generateOTP() string {
	rand.Seed(time.Now().UnixNano())
	code := rand.Intn(9000) + 1000 // 4-digit otp
	return fmt.Sprintf("%04d", code)
}

func (m *OrderModel) GetOrdersByUserID(userID primitive.ObjectID) ([]Order, error) {
	var orders []Order // slice of name orders, type Order

	cursor, err := m.collection.Find(context.Background(), bson.M{"placed_by": userID})
	if err != nil {
		return []Order{}, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var order Order
		if err := cursor.Decode(&order); err != nil {
			return []Order{}, err
		}
		orders = append(orders, order)
	}

	if orders == nil {
		return []Order{}, nil
	}

	return orders, nil
}

func (m *OrderModel) CancelOrder(userID, orderID primitive.ObjectID) error {

	var order Order
	err := m.collection.FindOne(context.Background(), bson.M{"_id": orderID}).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("No order found with this ID")
		}
		return err
	}

	if order.PlacedBy != userID {
		return fmt.Errorf("Unauthorized: You cannot cancel someone else's order")
	}

	filter := bson.M{
		"_id":       orderID,
		"placed_by": userID,
	}

	update := bson.M{
		"$set": bson.M{
			"status": "Cancelled",
		},
	}

	result, err := m.collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("No order found with this ID")
	}

	return nil
}
