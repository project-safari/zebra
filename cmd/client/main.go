package main

import (
	"fmt"
	"os"
)

// main puts together the execution of the commands.
func main() {
	if e := execRootCmd(); e != nil {
		os.Exit(1)
	}
}

// getting the root command to execute in main.
func execRootCmd() error {
	rootCmd := New()

	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
	}

	return err
}
