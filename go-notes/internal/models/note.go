package models

import "time"

type Note struct {
	NoteId    int       `json:"noteId" bson:"noteId"`
	Title     string    `json:"title" bson:"title"`
	Content   string    `json:"content" bson:"content"`
	CreatedAt time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
}

type Notes interface {
	CreateNote(note Note) (Response, error)
	GetAllNotes() (Response, error)
	GetNoteById(NoteId int) (Response, error)
	UpdateNote(NoteId int, updatedNote Note) (Response, error)
	DeleteNote(NoteId int) (Response, error)
}
