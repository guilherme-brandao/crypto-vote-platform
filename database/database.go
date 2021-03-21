package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var client *mongo.Client
var database *mongo.Database
var mongoCtx context.Context

func GetCollection(collectionName string) *mongo.Collection {
	return database.Collection(collectionName)
}

func CloseDBConnection() {
	fmt.Println("Closing MongoDB connection") 
	client.Disconnect(mongoCtx)
}

func GetContext() context.Context {
	return mongoCtx
}

func Init() {

	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGODB_URL")))
	if err != nil {
		log.Fatal(err)
	}
	mongoCtx = context.Background()
	err = client.Connect(mongoCtx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(mongoCtx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	database = client.Database("crypto-voting")

	fmt.Println("Connected to MongoDB!")

}
