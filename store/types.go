package store

import (
	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/auth"
	"github.com/project-safari/zebra/compute"
	"github.com/project-safari/zebra/dc"
	"github.com/project-safari/zebra/lease"
	"github.com/project-safari/zebra/network"
)

// DefaultFactory returns a resource factory with all the known types.
func DefaultFactory() zebra.ResourceFactory {
	factory := zebra.Factory()

	// Network resources.
	factory.Add(network.SwitchType())
	factory.Add(network.IPAddressPoolType())
	factory.Add(network.VLANPoolType())

	// DC resources.
	factory.Add(dc.DataCenterType())
	factory.Add(dc.LabType())
	factory.Add(dc.RackType())

	// Compute resources.
	factory.Add(compute.ServerType())
	factory.Add(compute.ESXType())
	factory.Add(compute.VCenterType())
	factory.Add(compute.VMType())

	// Zebra server resources.
	factory.Add(auth.UserType())

	// Zebra lease resource.
	factory.Add(lease.Type())

	// Need to add all the known types here
	return factory
}
