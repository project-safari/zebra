package pkg

import (
	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/network"
)

// function to generate IP pools.
//
// returns an array of type zebra.Resource.
func GenerateIPPool(numAddr int) []zebra.Resource {
	IPpool := make([]zebra.Resource, 0, numAddr)

	for i := 0; i < numAddr; i++ {
		labels := CreateLabels()

		ipArr := CreateIPArr(numAddr)

		IPaddr := network.NewIPAddressPool(ipArr, labels)

		if IPaddr.LabelsValidate() != nil {
			IPaddr.Labels = GroupLabels(IPaddr.Labels, GroupVal(IPaddr))
		}

		IPpool = append(IPpool, IPaddr)
	}

	return IPpool
}
