package pkg

import (
	"math"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/network"
)

// generate vlan resources.
func GenerateVlanPool(numVlans int) []zebra.Resource {
	delta := math.MaxUint16 / uint16(numVlans)

	start := uint16(0)

	vlans := make([]zebra.Resource, 0, numVlans)

	for i := 0; i < numVlans; i++ {
		labels := CreateLabels()
		vlan := network.NewVlanPool(start, start+delta-1, labels)

		vlans = append(vlans, vlan)

		start += delta
	}

	return vlans
}
