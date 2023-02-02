package main

import (
	"context"
	"fmt"

	"log"
	"sync"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"github.com/temporalio/samples-go/helloworld"
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

func getName() string {
	fmt.Println("Enter Your First Name: ")

	// first name
	var first string

	// Taking input from user.
	fmt.Scanln(&first)
	fmt.Println("Enter Your Last Name: ")

	// second name
	var second string

	// Taking input from the user.
	fmt.Scanln(&second)

	return first + " " + second
}

func firstClient(wg *sync.WaitGroup) {
	// The client and worker are heavyweight objects that should be created once per process.
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "hello-world", worker.Options{})

	w.RegisterWorkflow(helloworld.Workflow)
	w.RegisterActivity(helloworld.Activity)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}

func secondClient(wg *sync.WaitGroup) {
	// The client is a heavyweight object that should be created once per process.
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	workflowOptions := client.StartWorkflowOptions{
		ID:        "hello_world_workflowID",
		TaskQueue: "hello-world",
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, helloworld.Workflow, "Temporal")
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}

	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

	// Synchronously wait for the workflow completion.
	var result = getName()
	err = we.Get(context.Background(), &result)
	if err != nil {
		log.Fatalln("Unable get workflow result", err)
	}
	log.Println("Workflow result:", result)

}
