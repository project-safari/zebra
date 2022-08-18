package main

import (
	"fmt"
	"os"
)

// Main puts together the execution of the commands.
func main() {
	if e := execRootCmd(); e != nil {
		os.Exit(1)
	}
}

// Getting the root command to execute in main.
func execRootCmd() error {
	rootCmd := New()

	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
	}

	return err
}
