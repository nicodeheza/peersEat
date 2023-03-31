package models

import (
	"context"

	"github.com/nicodeheza/peersEat/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GeoCords struct {
	Long float64 `validate:"required"`
	Lat  float64 `validate:"required"`
}

type Peer struct {
	Id                  primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Url                 string               `bson:"url" json:"url" validate:"required,url"`
	Center              GeoCords             `validate:"dive"`
	City                string               `bson:"city,omitempty" json:"city,omitempty" validate:"required"`
	Country             string               `bson:"country,omitempty" json:"country,omitempty" validate:"required"`
	InfluenceRadius     float64              `bson:"influence_radius,omitempty" json:"influence_radius,omitempty"`
	DeliveryRadius      float64              `bson:"delivery_radius,omitempty" json:"delivery_radius,omitempty"`
	InAreaPeers         []primitive.ObjectID `bson:"in_area_peers,omitempty" json:"in_area_peers,omitempty"`
	InDeliveryAreaPeers []primitive.ObjectID `bson:"in_area_delivery_peers,omitempty" json:"in_area_delivery_peers,omitempty"`
}

func GetPeerColl(databaseName string) *mongo.Collection {
	collection := config.GetDatabase(databaseName).Collection("peers")
	return collection
}

func InitPeerModel(databaseName string) {
	GetPeerColl(databaseName).Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "url", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
}
