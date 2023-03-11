package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Client= nil

func ConnectDB(mongoUrl string){
	if db !=nil {return}
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoUrl))
	if err != nil{ log.Fatal(err)}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err= client.Connect(ctx)
	if err != nil{log.Fatal((err))}

	err= client.Ping(ctx, nil)
	if err != nil{log.Fatal((err))}

	fmt.Println("Connected to MongoDB")
	db= client
}

func GetDatabase(databaseName string) *mongo.Database{
	database := db.Database(databaseName)
	return database
}
