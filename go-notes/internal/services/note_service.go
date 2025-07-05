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
	db  *mongo.Collection
	ctx context.Context
}

func NewNoteService(db *mongo.Collection, ctx context.Context) *NoteService {
	return &NoteService{
		db:  db,
		ctx: ctx,
	}
}

func (s *NoteService) CreateNote(note models.Note) (models.Response, error) {
	log.Printf("Start creating new note!")
	if note.Title == "" {
		return models.Response{
			Status:  "error",
			Message: "Title is required!",
			Data:    note,
		}, errors.New("title is required")
	}
	now := time.Now()
	note.CreatedAt = now
	note.UpdatedAt = now
	data, err := s.db.InsertOne(s.ctx, note)
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

	cursor, err := s.db.Find(s.ctx, primitive.D{})
	if err != nil {
		return models.Response{
			Status:  "Error",
			Message: "Error while getting data from db",
		}, err
	}
	defer cursor.Close(s.ctx)
	var notes []models.Note
	for cursor.Next(s.ctx) {
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
	query := primitive.M{"noteId": NoteId}

	var note models.Note
	err := s.db.FindOne(s.ctx, query).Decode(&note)
	if err != nil {
		log.Printf("Note not found: %v", err)
		return models.Response{
			Status:  "error",
			Message: "Note not found",
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

	query := primitive.M{"noteId": NoteId}
	var note models.Note

	err := s.db.FindOne(s.ctx, query).Decode(&note)
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
	res, err := s.db.UpdateOne(s.ctx, query, update)
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
		Data: map[string]interface{}{
			"matchedCount":  res.MatchedCount,
			"modifiedCount": res.ModifiedCount,
		},
	}, nil
}
