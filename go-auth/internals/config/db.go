package config

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	MongoClient *mongo.Client
	MongoDB     *mongo.Database
)

func ConnectDb() error {
	url := os.Getenv("MONGO_URL")
	dbName := os.Getenv("DB_NAME")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpt := options.Client().ApplyURI(url)
	client, err := mongo.Connect(ctx, clientOpt)

	if err != nil {
		log.Fatalf("error connecting to MongoDB: %v", err)
		return err
	}
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("error pinging MongoDB: %v", err)
		return err
	}

	MongoClient = client
	MongoDB = client.Database(dbName)

	return nil
}

func GetDB() *mongo.Database {
	return MongoDB
}

func GetMonogClient() *mongo.Client {
	return MongoClient
}
