package main

// Contains the workflow that helps complete one/more activities.
// Execution.
// Sample Temporal by Eva Achim.

import (
	"context"
	"fmt"
	"hello-world-temporal/app"
	"log"

	"go.temporal.io/sdk/client"
)

func main() {

	// Create the client object just once per process
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	defer c.Close()

	options := client.StartWorkflowOptions{
		ID:        "greeting-workflow",
		TaskQueue: app.GreetingTaskQueue,
	}

	// The  Workflow.
	name := getName()
	we, err := c.ExecuteWorkflow(context.Background(), options, app.GreetingWorkflow, name)
	if err != nil {
		log.Fatalln("unable to complete Workflow", err)
	}

	// Get the results.
	var greeting string
	err = we.Get(context.Background(), &greeting)
	if err != nil {
		log.Fatalln("unable to get Workflow result", err)
	}

	printResults(greeting, we.GetID(), we.GetRunID())
}

// Print the results.
func printResults(greeting string, workflowID, runID string) {
	fmt.Printf("\nWorkflowID: %s RunID: %s\n", workflowID, runID)
	fmt.Printf("\n%s\n\n", greeting)
}

// Get user imput.
func getName() string {
	fmt.Println("Enter Your First Name: ")

	// first name
	var first string

	// Taking input from user.
	fmt.Scanln(&first)
	fmt.Println("Enter Second Last Name: ")

	// second name
	var second string

	// Taking input from the user.
	fmt.Scanln(&second)

	return first + " " + second
}
