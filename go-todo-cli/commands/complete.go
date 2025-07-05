package commands

import (
	"errors"
	"fmt"
	"go-todo-cli/tasks"
)

func Complete(id int) error {
	if id <= 0 {
		return errors.New("ID must be a positive number")
	}
	taskList, err := tasks.LoadTask()
	if err != nil {
		return err
	}

	found := false
	for i, task := range taskList {
		if task.ID == id {
			taskList[i].Completed = true
			found = true
			break
		}
	}

	if !found {
		fmt.Println("No task found with the given Id, please provide the correct taskId!")
		return errors.New("No task found with the given ID. Please use a valid task ID")
	}

	err = tasks.SaveTask(taskList)
	if err != nil {
		return err
	}

	fmt.Println("Task marked as completed.")
	return nil
}
