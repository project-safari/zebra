package app

// Contains an activity - a task to be completed in the workflow.
// Sample Temporal by Eva Achim.

import (
	"context"
	"fmt"
)

// The activity itself.
func ComposeGreeting(ctx context.Context, name string) (string, error) {
	greeting := fmt.Sprintf("Hello %s!", name)
	return greeting, nil
}
