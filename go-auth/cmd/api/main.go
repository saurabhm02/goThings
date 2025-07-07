package main

import (
	"context"
	"go-auth/internals/config"
	"go-auth/internals/handlers"
	"go-auth/internals/services"
	"log"
	"net/http"
	"os"
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
			log.Println("MongoDB connection closed successfully!...")
		}

	}()

	collection := config.GetDB().Collection("Users")
	service := services.NewUserService(collection)
	handlers := handlers.NewHandler(service)
	routes := handlers.RegisterRoutes()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("🚀 Server running at http://localhost:%s\n", port)
	http.ListenAndServe(":"+port, routes)

}
