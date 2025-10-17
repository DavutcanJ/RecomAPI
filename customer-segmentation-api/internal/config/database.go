package config

import (
    "context"
    "log"
    "os"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func ConnectDatabase() {
    mongoURI := os.Getenv("MONGODB_URI")
    if mongoURI == "" {
        mongoURI = "mongodb://localhost:27017"
    }

    dbName := os.Getenv("DB_NAME")
    if dbName == "" {
        dbName = "recom_api"
    }

    // MongoDB connection options for performance
    clientOptions := options.Client().
        ApplyURI(mongoURI).
        SetMaxPoolSize(100).                     // Maximum connections in pool
        SetMinPoolSize(5).                       // Minimum connections maintained
        SetMaxConnIdleTime(30 * time.Second).    // Close idle connections after 30s
        SetServerSelectionTimeout(5 * time.Second). // Server selection timeout
        SetConnectTimeout(10 * time.Second)      // Connection timeout

    client, err := mongo.Connect(context.TODO(), clientOptions)
    if err != nil {
        log.Fatal("Failed to connect to MongoDB:", err)
    }

    // Test the connection
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    err = client.Ping(ctx, nil)
    if err != nil {
        log.Fatal("Failed to ping MongoDB:", err)
    }

    DB = client.Database(dbName)
    log.Println("Connected to MongoDB successfully!")
}