package commands

import (
	"errors"
	"fmt"
	"go-todo-cli/tasks"
)

func Delete(id int) error {
	if id <= 0 {
		return errors.New("id required to delete a task!")
	}

	taskList, err := tasks.LoadTask()
	if err != nil {
		return err
	}

	found := false
	var updatedTasks []tasks.Task
	for _, task := range taskList {
		if task.ID == id {
			found = true
			continue
		}
		updatedTasks = append(updatedTasks, task)
	}
	if !found {
		fmt.Println("No task found with the given Id, please provide the correct taskId!")
		return errors.New("No task found with the given ID. Please use a valid task ID")
	}

	err = tasks.SaveTask(updatedTasks)
	if err != nil {
		return err
	}
	fmt.Println("task deleted successfully!")
	return nil
}
