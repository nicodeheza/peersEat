package models

import (
	"context"

	"github.com/nicodeheza/peersEat/config"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Center struct{
	Long float64
	Lat float64
}

type Peer struct {
	Id  				primitive.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Url 				string 				 `bson:"url" json:"url"`
	Center 				Center
	City 				string  			 `bson:"city,omitempty" json:"city,omitempty"`
	Country 			string  			 `bson:"country,omitempty" json:"country,omitempty"`
	InfluenceRadius 	float64 			 `bson:"influence_radius,omitempty" json:"influence_radius,omitempty"`
	DeliveryRadius 		float64 			 `bson:"delivery_radius,omitempty" json:"delivery_radius,omitempty"`
	InAreaPeers         []primitive.ObjectID `bson:"in_area_peers,omitempty" json:"in_area_peers,omitempty"`
	InDeliveryAreaPeers []primitive.ObjectID `bson:"in_area_delivery_peers,omitempty" json:"in_area_delivery_peers,omitempty"`
}

func GetPeerColl() *mongo.Collection{
	collection:= config.GetDatabase().Collection("peers")
	return collection
}

func InitPeerModel(){
	GetPeerColl().Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{{Key: "url", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
}