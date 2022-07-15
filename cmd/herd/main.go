package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/auth"
	"github.com/project-safari/zebra/cmd/herd/pkg"
	"github.com/project-safari/zebra/compute"
	"github.com/project-safari/zebra/dc"
	"github.com/project-safari/zebra/filestore"
	"github.com/project-safari/zebra/network"
	"github.com/spf13/cobra"
)

const (
	version             = "unknown"
	DefaultResourceSize = 100
)

func main() {
	name := filepath.Base(os.Args[0])
	rootCmd := &cobra.Command{ // nolint:exhaustruct,exhaustivestruct
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
		}(), "store"),
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

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	// Initialize the store
	rootDir := cmd.Flag("store").Value.String()
	fs := initStore(rootDir)

	n := intVal(cmd, "vlan-pool")
	vlans := pkg.GenerateVlanPool(n)

	for _, vlan := range vlans {
		if e := fs.Create(vlan); e != nil {
			return e
		}
	}

	fmt.Println("generated vlans:", len(vlans))

	return nil
}

func intVal(cmd *cobra.Command, flag string) int {
	v := cmd.Flag(flag).Value.String()
	i, _ := strconv.Atoi(v)

	return i
}

func initStore(rootDir string) *filestore.FileStore {
	fs := filestore.NewFileStore(rootDir, initTypes())
	if e := fs.Initialize(); e != nil {
		fmt.Println("Error initializing store")
		panic(e)
	}

	return fs
}

func initTypes() zebra.ResourceFactory {
	factory := zebra.Factory()

	// network resources
	factory.Add("Switch", func() zebra.Resource {
		return new(network.Switch)
	})
	factory.Add("IPAddressPool", func() zebra.Resource {
		return new(network.IPAddressPool)
	})
	factory.Add("VLANPool", func() zebra.Resource {
		return new(network.VLANPool)
	})

	// dc resources
	factory.Add("Datacenter", func() zebra.Resource {
		return new(dc.Datacenter)
	})
	factory.Add("Lab", func() zebra.Resource {
		return new(dc.Lab)
	})
	factory.Add("Rack", func() zebra.Resource {
		return new(dc.Rack)
	})

	// compute resources
	factory.Add("Server", func() zebra.Resource {
		return new(compute.Server)
	})
	factory.Add("ESX", func() zebra.Resource {
		return new(compute.ESX)
	})
	factory.Add("VCenter", func() zebra.Resource {
		return new(compute.VCenter)
	})
	factory.Add("VM", func() zebra.Resource {
		return new(compute.VM)
	})

	// other resources
	factory.Add("BaseResource", func() zebra.Resource {
		return new(zebra.BaseResource)
	})

	factory.Add("NamedResource", func() zebra.Resource {
		return new(zebra.NamedResource)
	})

	factory.Add("Credentials", func() zebra.Resource {
		return new(zebra.Credentials)
	})

	factory.Add("User", func() zebra.Resource {
		return new(auth.User)
	})

	// Need to add all the known types here
	return factory
}
