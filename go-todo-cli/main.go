package main

import (
	"fmt"
	"go-todo-cli/commands"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide any of the command: add, complete, delete, list")
		return
	}
	command := os.Args[1]

	switch command {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("please provide title to add any task to you list!")
			return
		}
		title := strings.Join(os.Args[2:], " ")
		err := commands.Add(title)
		if err != nil {
			fmt.Println("Error while adding task:", err)
		} else {
			fmt.Println("Task added successfully.")
		}
	case "complete":
		if len(os.Args) < 3 {
			fmt.Println("please provide id to complete the task!")
			return
		}
		idStr := strings.Join(os.Args[2:], " ")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Println("Invalid task ID. Please provide a number.")
			return
		}
		err = commands.Complete(id)
		if err != nil {
			fmt.Println("Error while complete a task:", err)
		} else {
			fmt.Println("task completed successfully!")
		}
	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("please provide id to complete the task!")
			return
		}
		idStr := strings.Join(os.Args[2:], " ")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			fmt.Println("Invalid task ID. Please provide a number.")
			return
		}
		err = commands.Delete(id)
		if err != nil {
			fmt.Println("Error while deleting a task:", err)
		}
	case "list":
		err := commands.List()
		if err != nil {
			fmt.Println("❌", err)
		}
	default:
		fmt.Printf("❌ Unknown command: %s\n", command)
		fmt.Println("Available commands: add, list, complete, delete")
	}
}
