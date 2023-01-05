package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/model"
	"github.com/project-safari/zebra/script/migration"
	"github.com/project-safari/zebra/store"
	"github.com/spf13/cobra"
)

const version = "unknown"

var Max = 200 //nolint:gochecknoglobals
// Max should ideally be: len(migration.Do()) //nolint:gochecknoglobals.

func migCmd() *cobra.Command {
	name := filepath.Base(os.Args[0])
	rootCmd := &cobra.Command{
		Use:          name,
		Short:        "mig",
		RunE:         run,
		SilenceUsage: true,
	}

	rootCmd.SetVersionTemplate(version + "\n")
	rootCmd.Flags().String("migrate", path.Join(
		func() string {
			s, _ := os.Getwd()

			return s
		}(), "mig_store"),
		"root directory of the store",
	)

	rootCmd.Flags().Int16("db-res", int16(len(migration.Do())), "number of db resources")

	return rootCmd
}

func execRootCmd() error {
	rootCmd := migCmd()

	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
	}

	return err
}

func main() { //nolint:wsl
	// migration.Post() is to execute the API post - github cannot access the DB and it affects the test check.
	// migration.Post()

	if e := execRootCmd(); e != nil {
		os.Exit(1)
	}
}

func storeResources(resources []zebra.Resource, fs *store.FileStore) error {
	for _, res := range resources {
		if e := fs.Create(res); e != nil {
			return e
		}
	}

	return nil
}

func genResources(cmd *cobra.Command,
	flag string,
	factory func(int) []zebra.Resource,
	resources []zebra.Resource,
) []zebra.Resource {
	n := intVal(cmd, flag)
	r := factory(n)

	fmt.Printf("generated %s: %d\n", flag, n)

	resources = append(resources, r...)

	return resources
}

// run for each resource.
func run(cmd *cobra.Command, _ []string) error {
	rootDir := cmd.Flag("migrate").Value.String()
	fs := initStore(rootDir)
	resources := make([]zebra.Resource, 0, Max)

	resources = genResources(cmd, "db-res", migration.DBData, resources)

	return storeResources(resources, fs)
}

func intVal(cmd *cobra.Command, flag string) int {
	v := cmd.Flag(flag).Value.String()
	i, _ := strconv.Atoi(v)

	return i
}

func initStore(rootDir string) *store.FileStore {
	fs := store.NewFileStore(rootDir, model.Factory())
	if e := fs.Initialize(); e != nil {
		fmt.Println("Error initializing store")
		panic(e)
	}

	return fs
}
