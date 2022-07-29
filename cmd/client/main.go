package main

import (
	"fmt"
	"os"
)

func main() {
	if e := execRootCmd(); e != nil {
		os.Exit(1)
	}
}

func execRootCmd() error {
	rootCmd := New()

	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
	}

	return err
}
