package commands

import (
	"errors"
	"fmt"
	"go-todo-cli/tasks"
	"strings"

	"github.com/google/uuid"
)

func Add(title string) error {
	if strings.TrimSpace(title) == "" {
		return errors.New("task title cannot be empty")
	}
	tasksList, err := tasks.LoadTask()
	if err != nil {
		return err
	}

	newTask := tasks.Task{
		ID:        uuid.New(),
		Title:     title,
		Completed: false,
	}

	tasksList = append(tasksList, newTask)
	err = tasks.SaveTask(tasksList)
	if err != nil {
		return err
	}
	fmt.Printf("✅ Added task: \"%s\"\n", title)
	return nil
}
