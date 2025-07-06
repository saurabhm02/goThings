package main

import (
	"context"
	"go-notes/internal/config"
	"go-notes/internal/handlers"
	"go-notes/internal/services"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
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

	collection := config.GetDB().Collection("notes")
	service := services.NewNoteService(collection)
	handler := handlers.NewNoteHandler(service)

	r := mux.NewRouter()
	r.HandleFunc("/notes/create", handler.CreateNote).Methods("POST")
	r.HandleFunc("/notes", handler.GetAllNotes).Methods("GET")
	r.HandleFunc("/notes/{id}", handler.GetNoteById).Methods("GET")
	r.HandleFunc("/notes/update/{id}", handler.UpdateNote).Methods("POST")
	r.HandleFunc("/notes/delete/{id}", handler.DeleteNote).Methods("DELETE")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("🚀 Server running at http://localhost:%s", port)
	http.ListenAndServe(":"+port, r)
}
