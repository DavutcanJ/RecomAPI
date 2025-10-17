package services

import (
    "context"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongo() *mongo.Client {
    client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
    if err != nil {
        panic(err)  // Production'da graceful handle
    }
    return client
}