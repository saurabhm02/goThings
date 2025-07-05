package commands

import (
	"errors"
	"fmt"
	"go-todo-cli/tasks"
	"strings"
)

func Add(title string) error {
	if strings.TrimSpace(title) == "" {
		return errors.New("task title cannot be empty")
	}
	tasksList, err := tasks.LoadTask()
	if err != nil {
		return err
	}
	newID := 1
	if len(tasksList) > 0 {
		newID = tasksList[len(tasksList)-1].ID + 1
	}

	newTask := tasks.Task{
		ID:        newID,
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
