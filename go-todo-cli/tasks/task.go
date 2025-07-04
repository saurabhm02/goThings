package tasks

import "github.com/google/uuid"

type Task struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
}
