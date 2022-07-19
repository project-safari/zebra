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
	rootCmd := &cobra.Command{ //nolint // exhaustruct,exhaustivestruct
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

//nolint // need all return statements.
func run(cmd *cobra.Command, _ []string) error {
	// Initialize the store.
	rootDir := cmd.Flag("store").Value.String()
	fs := initStore(rootDir)

	// vlans.
	n := intVal(cmd, "vlan-pool")
	vlans := pkg.GenerateVlanPool(n)

	for _, vlan := range vlans {
		if e := fs.Create(vlan); e != nil {
			return e
		}
	}

	fmt.Println("generated vlans:", len(vlans))

	// switches.
	s := intVal(cmd, "switch")
	switches := pkg.GenerateSwitch(s)

	for _, sw := range switches {
		if e := fs.Create(sw); e != nil {
			return e
		}
	}

	fmt.Println("generated switches:", len(switches))

	// ip-address-pool.
	addr := intVal(cmd, "ip-address-pool")
	addresses := pkg.GenerateSwitch(addr)

	for _, a := range addresses {
		if e := fs.Create(a); e != nil {
			return e
		}
	}

	fmt.Println("generated ip-address-pools:", len(addresses))

	// dc.
	datac := intVal(cmd, "dc")
	centers := pkg.GenerateIPPool(datac)

	for _, dc := range centers {
		if e := fs.Create(dc); e != nil {
			return e
		}
	}

	fmt.Println("generated data centers:", len(centers))

	// server.
	srv := intVal(cmd, "server")
	servers := pkg.GenerateServer(srv)

	for _, sv := range servers {
		if e := fs.Create(sv); e != nil {
			return e
		}
	}

	fmt.Println("generated servers:", len(servers))

	// vm.
	machine := intVal(cmd, "vm")
	vms := pkg.GenerateVM(machine)

	for _, vm := range vms {
		if e := fs.Create(vm); e != nil {
			return e
		}
	}

	fmt.Println("generated virtual machines:", len(vms))

	// racks.
	rack := intVal(cmd, "rack")
	racks := pkg.GenerateRack(rack)

	for _, r := range racks {
		if e := fs.Create(r); e != nil {
			return e
		}
	}

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

//nolint // need init for each type.
func initTypes() zebra.ResourceFactory {
	factory := zebra.Factory()

	// network resources.
	vln := network.VLANPoolType()
	factory.Add(vln)

	// IP pool init.
	addP := network.IPAddressPoolType()
	factory.Add(addP)

	// switch init.
	sw := network.SwitchType()
	factory.Add(sw)

	// dc resources.
	// lab init.
	l := dc.LabType()
	factory.Add(l)

	// rack init.
	ra := dc.RackType()
	factory.Add(ra)

	// dc init.
	dc := dc.DataCenterType()
	factory.Add(dc)

	// compute resources.
	esx := compute.ESXType()
	factory.Add(esx)

	srv := compute.ServerType()
	factory.Add(srv)

	// VC init.
	VC := compute.VCenterType()
	factory.Add(VC)

	// vm init.
	VM := compute.VMType()
	factory.Add(VM)

	// other resources.
	base := new(zebra.Type)

	base.Name = "Base Resource"
	base.Description = "some base resources"
	base.Constructor = func() zebra.Resource { return new(zebra.BaseResource) }

	factory.Add(*base)

	// named resource init.
	hasName := new(zebra.Type)

	hasName.Name = "Base Resource"
	hasName.Description = "some named resources"
	hasName.Constructor = func() zebra.Resource { return new(zebra.NamedResource) }

	factory.Add(*hasName)

	// credentials init.
	crd := new(zebra.Type)

	crd.Name = "Credentials"
	crd.Description = "some credentials"
	crd.Constructor = func() zebra.Resource { return new(zebra.Credentials) }

	factory.Add(*crd)

	// users init.
	u := auth.UserType()
	factory.Add(u)

	// Need to add all the known types here
	return factory
}
