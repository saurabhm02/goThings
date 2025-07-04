package commands

import (
	"errors"
	"go-todo-cli/tasks"

	"github.com/google/uuid"
)

func Complete(id uuid.UUID) error {
	if id == nil {
		return errors.New("Id cannnot be empty")
	}
	taskList, err := tasks.LoadTask()
	if err != nil {
		return err
	}
	found := false

}
