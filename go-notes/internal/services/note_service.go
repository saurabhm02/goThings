package services

import (
	"context"
	"errors"
	"go-notes/internal/models"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type NoteService struct {
	db *mongo.Collection
}

func NewNoteService(db *mongo.Collection) *NoteService {
	return &NoteService{
		db: db,
	}
}

func (s *NoteService) CreateNote(note models.Note) (models.Response, error) {
	log.Printf("Start creating new note!")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if note.Title == "" {
		return models.Response{
			Status:  "error",
			Message: "Title is required!",
			Data:    note,
		}, errors.New("title is required")
	}
	var existingNote models.Note
	err := s.db.FindOne(ctx, primitive.M{"noteId": note.NoteId}).Decode(&existingNote)
	if err == nil {
		return models.Response{}, errors.New("noteId already exists")
	}

	count, err := s.db.CountDocuments(ctx, primitive.M{})
	if err != nil {
		log.Printf("Failed to count notes: %v", err)
		return models.Response{
			Status:  "error",
			Message: "Failed to generate noteId",
		}, err
	}
	note.NoteId = int(count) + 1
	now := time.Now()
	note.CreatedAt = now
	note.UpdatedAt = now
	data, err := s.db.InsertOne(ctx, note)
	if err != nil {
		log.Printf("Failed to insert note: %v", err)
		return models.Response{
			Status:  "error",
			Message: "Failed to insert note",
		}, err
	}
	id, ok := data.InsertedID.(primitive.ObjectID)
	if !ok {
		return models.Response{
			Status:  "error",
			Message: "Failed to parse inserted ID",
		}, errors.New("failed to convert inserted ID to ObjectID")
	}
	log.Printf("note created successfully: %s", id.Hex())
	return models.Response{
		Status:  "success",
		Message: "note created successfully!",
		Data: map[string]interface{}{
			"id": id.Hex(),
		},
	}, nil
}

func (s *NoteService) GetAllNotes() (models.Response, error) {
	log.Println("Getting all the notes!")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cursor, err := s.db.Find(ctx, primitive.D{})
	if err != nil {
		return models.Response{
			Status:  "Error",
			Message: "Error while getting data from db",
		}, err
	}
	defer cursor.Close(ctx)
	var notes []models.Note
	for cursor.Next(ctx) {
		var note models.Note
		if err := cursor.Decode(&note); err != nil {
			log.Println("Decode error:", err)
			continue
		}
		notes = append(notes, note)
	}
	if err := cursor.Err(); err != nil {
		return models.Response{
			Status:  "error",
			Message: "Cursor error while fetching notes",
		}, err
	}
	return models.Response{
		Status:  "success",
		Message: "Successfully fetched all the notes",
		Data:    notes,
	}, nil
}

func (s *NoteService) GetNoteById(NoteId int) (models.Response, error) {
	log.Println("Getting note by noteId:", NoteId)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	query := primitive.M{"noteId": NoteId}

	var note models.Note
	err := s.db.FindOne(ctx, query).Decode(&note)
	if err != nil {
		log.Printf("no note foun in db: %v", err)
		return models.Response{
			Status:  "success",
			Message: "no note foun in db",
		}, err
	}
	return models.Response{
		Status:  "success",
		Message: "Note fetched successfully",
		Data:    note,
	}, nil
}

func (s *NoteService) UpdateNote(NoteId int, updatedNote models.Note) (models.Response, error) {
	log.Println("Updating existing note!")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := primitive.M{"noteId": NoteId}
	var note models.Note

	err := s.db.FindOne(ctx, query).Decode(&note)
	if err != nil {
		log.Printf("Note not found: %v", err)
		return models.Response{
			Status:  "error",
			Message: "Note not found",
		}, err
	}
	updatedNote.UpdatedAt = time.Now()
	update := primitive.M{
		"$set": updatedNote,
	}
	res, err := s.db.UpdateOne(ctx, query, update)
	if err != nil {
		log.Printf("Failed to update note: %v", err)
		return models.Response{
			Status:  "error",
			Message: "Failed to update note",
		}, err
	}

	return models.Response{
		Status:  "success",
		Message: "Note updated successfully",
		Data:    res,
	}, nil
}

func (s *NoteService) DeleteNote(noteId int) (models.Response, error) {
	log.Println("Deleting a note with given noteID", noteId)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := primitive.M{"noteId": noteId}
	_, err := s.db.DeleteOne(ctx, query)
	if err != nil {
		log.Printf("Failed to delete note: %v", err)
		return models.Response{
			Status:  "Error",
			Message: "Failed to delete note",
		}, err
	}
	log.Println("note deleted successfully!")
	return models.Response{
		Status:  "Success",
		Message: "note deleted successfully",
	}, nil
}
