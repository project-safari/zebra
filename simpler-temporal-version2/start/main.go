package main

import (
	"context"
	"fmt"
	"hello-world-temporal/app"
	"log"
        "go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/client"
        "time"
	"sync"
)

func main() {
	var count = 2
	var wg sync.WaitGroup

	wg.Add(count)


go firstClient(&wg)

   defer wg.Done()

   go secondClient(&wg)

wg.Wait()
}

func printResults(greeting string, workflowID, runID string) {
	fmt.Printf("\nWorkflowID: %s RunID: %s\n", workflowID, runID)
	fmt.Printf("\n%s\n\n", greeting)
}

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


func secondClient(wg *sync.WaitGroup){
		defer wg.Done()

		time.Sleep(time.Second * 2)

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

	// The  Workflow
	name := getName()
	we, err := c.ExecuteWorkflow(context.Background(), options, app.GreetingWorkflow, name)
	if err != nil {
		log.Fatalln("unable to complete Workflow", err)
	}

	// Get the results
	var greeting string
	err = we.Get(context.Background(), &greeting)
	if err != nil {
		log.Fatalln("unable to get Workflow result", err)
	}

		printResults(greeting, we.GetID(), we.GetRunID())
               
	}


func firstClient(wg *sync.WaitGroup) {
		defer wg.Done()
	// Create the client object just once per process
    c, err := client.Dial(client.Options{})
    if err != nil {
        log.Fatalln("unable to create Temporal client", err)
    }
    defer c.Close()

    // This worker hosts both Workflow and Activity functions
    w := worker.New(c, app.GreetingTaskQueue, worker.Options{})
    w.RegisterWorkflow(app.GreetingWorkflow)
    w.RegisterActivity(app.ComposeGreeting)

    // Start listening to the Task Queue
    err = w.Run(worker.InterruptCh())
    if err != nil {
        log.Fatalln("unable to start Worker", err)
    }

wg.Done()
}






