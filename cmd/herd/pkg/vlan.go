package pkg

import (
	"math"

	"github.com/project-safari/zebra/network"
)

const MAX uint16 = 65535

// generate vlan resources.
func GenerateVlanPool(numVlans int) []*network.VLANPool {
	var delta uint16

	var divide uint16

	if numVlans > 0 && numVlans < math.MaxUint16 {
		divide = uint16(numVlans)
	}

	delta = MAX / divide // accepted conversion.

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
