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

		cent := compute.NewVCenter(name, ip, labels)

		if cent.Labels.Validate() != nil {
			cent.Labels = GroupLabels(cent.Labels, GroupVal(cent))
		}

		centers = append(centers, cent)
	}

	return centers
}
