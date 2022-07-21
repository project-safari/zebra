package pkg

import (
	"net"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/compute"
)

// gererate VCenter.
func GenerateVCenter(numVC int) []zebra.Resource {
	centers := make([]zebra.Resource, 0, numVC)

	for i := 0; i < numVC; i++ {
		name := Name()
		labels := CreateLabels()
		ip := net.IP(RandData(IPsamples()))

		center := compute.NewVCenter(name, ip, labels)

		centers = append(centers, center)
	}

	return centers
}
