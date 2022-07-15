package pkg

import (
	"github.com/project-safari/zebra/network"
)

const MAX = uint16(0xFFFF)

func GenerateVlanPool(numVlans int) []*network.VLANPool {
	delta := MAX / uint16(numVlans)
	start := uint16(0)
	vlans := make([]*network.VLANPool, 0, numVlans)

	for i := 0; i < numVlans; i++ {
		labels := CreateLabels()
		vlan := network.NewVlanPool(start, start+delta-1, labels)

		vlans = append(vlans, vlan)
		start += delta
	}

	return vlans
}
