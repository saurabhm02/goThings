package commands

import (
	"errors"
	"fmt"
	"go-todo-cli/tasks"
)

func List() error {
	taskList, err := tasks.LoadTask()
	if err != nil {
		return err
	}
	if len(taskList) == 0 {
		return errors.New("No task list found cuurently in the db!, you can add task")
	}
	for i, task := range taskList {
		status := "[ ]"
		if task.Completed {
			status = "[✅]"
		}
		fmt.Printf("%d. %s %s (%s)\n", i+1, status, task.Title, task.ID)
	}
	return nil
}
