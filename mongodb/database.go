package mongodb

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var client *mongo.Client
var database *mongo.Database

func ConnectDB(uri, dbName string) error {
	clientOptions := options.Client().ApplyURI(uri)
	ctx := context.TODO()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
		return err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
		return err
	}

	database = client.Database(dbName)
	fmt.Println("[DB]: Connected to MongoDB!")
	return nil
}

func GetCollection(collectionName string) *mongo.Collection {
	return database.Collection(collectionName)
}
