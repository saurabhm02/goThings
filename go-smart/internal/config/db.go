package config

import (
	"fmt"
	"go-smart/internal/utils"
	"os"

	"github.com/gofiber/fiber/v2/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	MongoClient *mongo.Client
	MongoDb     *mongo.Database
)

func ConnectDb() error {
	ctx, cancel := utils.GetContext()
	defer cancel()

	url := os.Getenv("MONGO_URL")
	dbName := os.Getenv("DB_NAME")

	if url == "" || dbName == "" {
		return fmt.Errorf("MONGO_URL or DB_NAME environment variable not set")
	}

	clientOptions := options.Client().ApplyURI(url)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Errorf("Error connecting to MongoDB: %v", err)
		return err
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Errorf("MongoDB ping failed: %v", err)
		return err
	}

	MongoClient = client
	MongoDb = client.Database(dbName)

	log.Info("Successfully connected to the database!")
	return nil
}

func GetDB() *mongo.Database {
	return MongoDb
}

func GetMongoClient() *mongo.Client {
	return MongoClient
}
