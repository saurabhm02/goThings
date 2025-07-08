package config

import (
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

	clientOptions := options.Client().ApplyURI(url)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Errorf("Error while connecting to db! %v", err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		log.Errorf("Error while pinging with the db!", err)
	}
	MongoClient = client
	MongoDb = client.Database(dbName)
	log.Info("successfuly database connected!")
	return nil
}

func GetDB() *mongo.Database {
	return MongoDb
}

func GetMongoClient() *mongo.Client {
	return MongoClient
}
