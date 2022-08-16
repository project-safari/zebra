package main //nolint:testpackage

import (
	"os"
	"path"
	"testing"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/cmd/herd/pkg"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func mockCmd() *cobra.Command {
	cmd := &cobra.Command{ //nolint:exhaustruct,exhaustivestruct
		Use:          "Name",
		Short:        "herd",
		Version:      version + "\n",
		RunE:         run,
		SilenceUsage: true,
	}

	cmd.Flags().String("store", path.Join(
		func() string {
			s, _ := os.Getwd()

			return s
		}(), "herd_store"),
		"root directory of the store",
	)

	cmd.Flags().Int16("user", DefaultUserSize, "number of users")
	cmd.Flags().Int16("vlan-pool", DefaultResourceSize, "number of vlan pools")
	cmd.Flags().Int16("switch", DefaultResourceSize, "number of switches")

	cmd.Flags().Int16("ip-address-pool", DefaultResourceSize, "number of ip address pools")
	cmd.Flags().Int16("dc", DefaultResourceSize, "number of data centers")

	cmd.Flags().Int16("server", DefaultResourceSize, "number of servers")
	cmd.Flags().Int16("vm", DefaultResourceSize, "number of vms")

	cmd.Flags().Int16("rack", DefaultResourceSize, "number of racks")
	cmd.Flags().Int16("vcenter", DefaultResourceSize, "number of vcenters")

	cmd.Flags().Int16("esx", DefaultResourceSize, "number of esx servers")
	cmd.Flags().Int16("lab", DefaultResourceSize, "number of labs")

	return cmd
}

func TestInit(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	stored := initStore("test_store")
	assert.NotNil(stored)
	assert.Nil(os.RemoveAll("test_store"))
}

func TestFileStore(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	fs := initStore("herd_file_store")

	defer os.RemoveAll("herd_file_store")

	resources := make([]zebra.Resource, 0, 2)
	res := storeResources(resources, fs)

	assert.Nil(res)

	labels := pkg.CreateLabels()
	resources = append(resources, zebra.NewBaseResource(zebra.DefaultType(), labels))
	ret := storeResources(resources, fs)

	assert.Nil(ret)
}

func TestInitVal(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	vals := intVal(mockCmd(), "user")

	assert.NotNil(vals)

	assert.True(vals > 0)
}

func TestGenerateStore(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	resources := make([]zebra.Resource, 0, Max)

	res := genResources(mockCmd(), "vlan-pool", pkg.GenerateVlanPool, resources)

	assert.NotNil(res)

	assert.True(len(res) > 0)
}

func TestRun(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	arr := []string{}

	exec := run(mockCmd(), arr)

	assert.Nil(exec)
}
