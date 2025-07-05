package config

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	MongoClient *mongo.Client
	MongoDb     *mongo.Database
)

func ConnectDb() error {
	url := os.Getenv("MONGO_URL")
	dbName := os.Getenv("DB_NAME")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(url)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("error connecting to MongoDB: %w", err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("error pinging MongoDB: %w", err)
	}
	MongoClient = client
	MongoDb = client.Database(dbName)

	fmt.Println("✅ Successfully connected to MongoDB")
	return nil
}

func GetDB() *mongo.Database {
	return MongoDb
}

func GetMongoClient() *mongo.Client {
	return MongoClient
}
