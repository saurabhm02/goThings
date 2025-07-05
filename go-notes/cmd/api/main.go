package main

import (
	"context"
	"go-notes/internal/config"
	"log"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := godotenv.Load(); err != nil {
		log.Fatalf("❌ Error while loading .env file: %v", err)
	}
	if err := config.ConnectDb(); err != nil {
		log.Fatalf("Error while connectin db, %v", err)
	}
	db := config.GetMongoClient()
	defer func() {
		if err := db.Disconnect(ctx); err != nil {
			log.Fatalf("Failed to disconnect MongoDB: %v", err)
		} else {
			log.Fatalf("MongoDB connection closed successfully..")
		}
	}()

}
