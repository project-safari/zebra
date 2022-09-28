package model

import (
	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/model/compute"
	"github.com/project-safari/zebra/model/dc"
	"github.com/project-safari/zebra/model/lease"
	"github.com/project-safari/zebra/model/network"
	"github.com/project-safari/zebra/model/user"
)

// Factory returns a resource factory with all the known types.
func Factory() zebra.ResourceFactory {
	factory := zebra.Factory()

	// network resources
	factory.Add(network.SwitchType(), network.EmptySwitch)
	factory.Add(network.IPAddressPoolType(), network.EmptyIPAddressPool)
	factory.Add(network.VLANPoolType(), network.EmptyVLANPool)

	// dc resources
	factory.Add(dc.DataCenterType(), dc.EmptyDataCenter)
	factory.Add(dc.LabType(), dc.EmptyLab)
	factory.Add(dc.RackType(), dc.EmptyRack)

	// compute resources
	factory.Add(compute.ServerType(), compute.EmptyServer)
	factory.Add(compute.ESXType(), compute.EmptyESX)
	factory.Add(compute.VCenterType(), compute.EmptyVCenter)
	factory.Add(compute.VMType(), compute.EmptyVM)

	// zebra server resources
	factory.Add(user.Type(), user.Empty)

	// zebra lease resource
	factory.Add(lease.Type(), lease.Empty)

	// Need to add all the known types here
	return factory
}
