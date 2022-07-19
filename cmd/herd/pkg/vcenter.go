package pkg

import (
	"net"

	"github.com/project-safari/zebra/compute"
)

// gererate VCenter.
func GenerateVcenter(numVC int) []*compute.VCenter {
	centers := make([]*compute.VCenter, 0, numVC)

	for i := 0; i < numVC; i++ {
		name := Name()
		labels := CreateLabels()
		ip := net.IP(RandData(IPsamples()))

		center := compute.NewVCenter(name, ip, labels)

		centers = append(centers, center)
	}

	return centers
}
