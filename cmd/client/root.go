package main

import (
	"os"
	"path"
	"path/filepath"

	"github.com/spf13/cobra"
)

var version string

func New() *cobra.Command {
	name := filepath.Base(os.Args[0])

	rootCmd := &cobra.Command{
		Use:     name,
		Short:   "zebra resource reservation client",
		Version: version + "\n",
	}
	rootCmd.SetVersionTemplate(version + "\n")
	rootCmd.PersistentFlags().StringP("config", "c",
		path.Join(os.Getenv("HOME"), ".zebra.yaml"),
		"config file",
	)
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose output")

	rootCmd.AddCommand(NewConfigure())

	return rootCmd
}
