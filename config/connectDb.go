package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(GetEnv("MONGO_URI")))
	if err != nil{ log.Fatal(err)}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err= client.Connect(ctx)
	if err != nil{log.Fatal((err))}

	err= client.Ping(ctx, nil)
	if err != nil{log.Fatal((err))}

	fmt.Println("Connected to MongoDB")
	return client
}

var DB *mongo.Client = ConnectDB()

func GetCollection (client *mongo.Client, collectionName string) *mongo.Collection{
	collection := client.Database("peersEatDB").Collection(collectionName)
	return collection
}