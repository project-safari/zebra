package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/cmd/herd/pkg"
	"github.com/project-safari/zebra/filestore"
	"github.com/project-safari/zebra/store"
	"github.com/spf13/cobra"
)

const (
	version             = "unknown"
	DefaultResourceSize = 100
	DefaultUserSize     = 10
	Max                 = 1500
)

func main() {
	name := filepath.Base(os.Args[0])
	rootCmd := &cobra.Command{
		Use:          name,
		Short:        "herd",
		Version:      version + "\n",
		RunE:         run,
		SilenceUsage: true,
	}

	rootCmd.SetVersionTemplate(version + "\n")
	rootCmd.Flags().String("store", path.Join(
		func() string {
			s, _ := os.Getwd()

			return s
		}(), "herd_store"),
		"root directory of the store",
	)

	rootCmd.Flags().Int16("vlan-pool", DefaultResourceSize, "number of vlan pools")
	rootCmd.Flags().Int16("switch", DefaultResourceSize, "number of switches")

	rootCmd.Flags().Int16("ip-address-pool", DefaultResourceSize, "number of ip address pools")
	rootCmd.Flags().Int16("dc", DefaultResourceSize, "number of data centers")

	rootCmd.Flags().Int16("server", DefaultResourceSize, "number of servers")
	rootCmd.Flags().Int16("vm", DefaultResourceSize, "number of vms")

	rootCmd.Flags().Int16("rack", DefaultResourceSize, "number of racks")
	rootCmd.Flags().Int16("vcenter", DefaultResourceSize, "number of vcenters")

	rootCmd.Flags().Int16("esx", DefaultResourceSize, "number of esx servers")
	rootCmd.Flags().Int16("lab", DefaultResourceSize, "number of labs")

	rootCmd.Flags().Int16("user", DefaultUserSize, "number of users")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func storeResources(resources []zebra.Resource, fs *filestore.FileStore) error {
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
	rootDir := cmd.Flag("store").Value.String()
	fs := initStore(rootDir)
	resources := make([]zebra.Resource, 0, Max)

	// Generate all the resources
	resources = genResources(cmd, "vlan-pool", pkg.GenerateVlanPool, resources)
	resources = genResources(cmd, "switch", pkg.GenerateSwitch, resources)
	resources = genResources(cmd, "ip-address-pool", pkg.GenerateIPPool, resources)
	resources = genResources(cmd, "dc", pkg.GenerateDatacenter, resources)
	resources = genResources(cmd, "server", pkg.GenerateServer, resources)
	resources = genResources(cmd, "vm", pkg.GenerateVM, resources)
	resources = genResources(cmd, "rack", pkg.GenerateRack, resources)
	resources = genResources(cmd, "lab", pkg.GenerateLab, resources)
	resources = genResources(cmd, "user", pkg.GenerateUser, resources)

	return storeResources(resources, fs)
}

func intVal(cmd *cobra.Command, flag string) int {
	v := cmd.Flag(flag).Value.String()
	i, _ := strconv.Atoi(v)

	return i
}

func initStore(rootDir string) *filestore.FileStore {
	fs := filestore.NewFileStore(rootDir, store.DefaultFactory())
	if e := fs.Initialize(); e != nil {
		fmt.Println("Error initializing store")
		panic(e)
	}

	return fs
}
