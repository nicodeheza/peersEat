package models

import (
	"context"

	"github.com/nicodeheza/peersEat/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Rate struct {
	Stars uint
	Votes uint
}

type Restaurant struct {
	Id                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name              string
	Address           string
	City              string
	Country           string
	Coord             GeoCoords
	IsConnected       bool    `bson:"isConnected,omitempty" json:"isConnected,omitempty"`
	Menu              Menu    `bson:"menu,omitempty" json:"menu,omitempty"`
	OpenTime          string  `bson:"openTime,omitempty" json:"openTime,omitempty"`
	CloseTime         string  `bson:"closeTime,omitempty" json:"closeTime,omitempty"`
	Rate              Rate    `bson:"rate,omitempty" json:"rate,omitempty"`
	Phone             string  `bson:"Phone,omitempty" json:"phone,omitempty"`
	DeliveryCost      float32 `bson:"deliveryCost,omitempty" json:"deliveryCost,omitempty"`
	IsDeliveryFixCost bool    `bson:"isDeliveryFixCost,omitempty" json:"isDeliveryFixCost,omitempty"`
	MinDeliveryTime   uint    `bson:"minDeliveryTime,omitempty" json:"minDeliveryTime,omitempty"`
	MaxDeliveryTime   uint    `bson:"maxDeliveryTime,omitempty" json:"maxDeliveryTime,omitempty"`
	DeliveryRadius    float64 `bson:"deliveryRadius,omitempty" json:"deliveryRadius,omitempty"`
	UserName          string  `bson:"userName,omitempty" json:"userName,omitempty"`
	Password          string  `bson:"password,omitempty" json:"password,omitempty"`
	IsFinalPassword   bool    `bson:"isFinalPassword,omitempty" json:"isFinalPassword,omitempty"`
}

func GetRestaurantColl(databaseName string) *mongo.Collection {
	return config.GetDatabase(databaseName).Collection("restaurants")
}

func InitRestaurantModel(databaseName string) {
	GetRestaurantColl(databaseName).Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "userName", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
}
