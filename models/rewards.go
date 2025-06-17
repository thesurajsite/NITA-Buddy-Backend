package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Rewards struct {
	ID    primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Coins int                `bson:"coins" json:"coins"`
}

type RewardsModel struct {
	collection *mongo.Collection
}

func NewRewardsModel(collection *mongo.Collection) *RewardsModel {
	return &RewardsModel{collection: collection}
}

func (r *RewardsModel) CreateRewardsOnSignup(userID primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	reward := Rewards{
		ID:    userID,
		Coins: 50,
	}

	_, err := r.collection.InsertOne(ctx, reward)
	return err
}

func (r *RewardsModel) GetRewardsByUserID(userID primitive.ObjectID) (*Rewards, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var reward Rewards
	err := r.collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&reward)
	if err != nil {
		return nil, err
	}

	return &reward, nil
}
