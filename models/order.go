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
	Phone         string             `bson:"phone" json:"phone"`
	PlacedBy      primitive.ObjectID `bson:"placed_by" json:"placed_by"`
	PlacedByName  string             `bson:"placed_by_name" json:"placed_by_name"`
	AcceptedBy    primitive.ObjectID `bson:"accepted_by" json:"accepted_by"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
}

type OrderModel struct {
	collection     *mongo.Collection
	userCollection *mongo.Collection
	rewardsModel   *RewardsModel // Add this field
}

func NewOrderModel(orderCollection, userCollection *mongo.Collection, rewardsModel *RewardsModel) *OrderModel {
	return &OrderModel{
		collection:     orderCollection,
		userCollection: userCollection,
		rewardsModel:   rewardsModel,
	}
}

func (m *OrderModel) CreateOrder(store, orderDetails string, placedBy primitive.ObjectID) (*Order, error) {

	reward, err := m.rewardsModel.GetRewardsByUserID(placedBy)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve reward info: %v", err)
	}

	if reward.Coins < 10 {
		return nil, fmt.Errorf("not enough coins to place the order")
	}

	status := "NotAccepted"
	otp := generateOTP()
	acceptedBy := primitive.NilObjectID
	customID, err := generateCustomOrderID(m.collection)
	if err != nil {
		return nil, err
	}

	var userDetails struct {
		Name  string `bson:"name"`
		Phone string `bson:"phone"`
	}
	err = m.userCollection.FindOne(context.Background(), bson.M{"_id": placedBy}).Decode(&userDetails)
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %v", err)
	}

	order := &Order{
		CustomOrderID: customID,
		Store:         store,
		OrderDetails:  orderDetails,
		Status:        status,
		OTP:           otp,
		Phone:         userDetails.Phone,
		PlacedBy:      placedBy,
		PlacedByName:  userDetails.Name,
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

func (m *OrderModel) GetOtherIncompleteOrders(userID primitive.ObjectID) ([]Order, error) {
	var orders []Order

	filter := bson.M{
		"placed_by": bson.M{"$ne": userID}, // ne : not equal
		"status":    "NotAccepted",
	}

	cursor, err := m.collection.Find(context.Background(), filter)
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

	return orders, nil
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
			return fmt.Errorf("no order found with this id")
		}
		return err
	}

	if order.PlacedBy != userID {
		return fmt.Errorf("unauthorized: you cannot cancel someone else's order")
	}

	filter := bson.M{
		"_id":       orderID,
		"placed_by": userID,
	}

	result, err := m.collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("no order found with this id")
	}

	return nil
}

func (m *OrderModel) AcceptOrder(userID, orderID primitive.ObjectID) error {

	var order Order
	err := m.collection.FindOne(context.Background(), bson.M{"_id": orderID}).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("order not found")
		}
		return err
	}

	if order.AcceptedBy != primitive.NilObjectID {
		return fmt.Errorf("order already accepted")
	}

	if order.PlacedBy == userID {
		return fmt.Errorf("can't accept own order")
	}

	update := bson.M{
		"$set": bson.M{
			"accepted_by": userID,
			"status":      "Accepted",
		},
	}

	_, err = m.collection.UpdateOne(context.Background(), bson.M{"_id": orderID}, update)
	if err != nil {
		return fmt.Errorf("failed to accept order: %v", err)
	}

	return nil
}

func (m *OrderModel) GetAcceptedOrders(userID primitive.ObjectID) ([]Order, error) {
	var orders []Order

	filter := bson.M{
		"accepted_by": userID,
		"status":      "Accepted",
	}

	cursor, err := m.collection.Find(context.Background(), filter)
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

	return orders, nil
}

// accept order
//requirements : userID, orderID, otp
// check : userID==acceptedby, otp==otp, status=accepted,
// todo : status=completed, reward--, reward++

func (m *OrderModel) CompleteOrder(userID, orderID primitive.ObjectID, otp string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Fetch the order from DB
	var order Order
	err := m.collection.FindOne(ctx, bson.M{"_id": orderID}).Decode(&order)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("order not found")
		}
		return fmt.Errorf("failed to fetch order: %v", err)
	}

	// Verify user is the one who accepted the order
	if order.AcceptedBy != userID {
		return fmt.Errorf("unauthorized: only the user who accepted the order can complete it")
	}

	// Verify order is in accepted state
	if order.Status != "Accepted" {
		return fmt.Errorf("order is not in accepted state")
	}

	// Verify OTP
	if order.OTP != otp {
		return fmt.Errorf("invalid OTP")
	}

	// Update order status to Completed
	update := bson.M{
		"$set": bson.M{"status": "Completed"},
	}
	_, err = m.collection.UpdateOne(ctx, bson.M{"_id": orderID}, update)
	if err != nil {
		return fmt.Errorf("failed to update order status: %v", err)
	}

	// Decrease 10 coins from PlacedBy
	_, err = m.rewardsModel.UpdateCoins(order.PlacedBy, -10)
	if err != nil {
		return fmt.Errorf("failed to deduct coins from order creator: %v", err)
	}

	// Increase 10 coins to AcceptedBy
	_, err = m.rewardsModel.UpdateCoins(order.AcceptedBy, 10)
	if err != nil {
		return fmt.Errorf("failed to add coins to order accepter: %v", err)
	}

	return nil
}
