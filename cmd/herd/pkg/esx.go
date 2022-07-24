package pkg

import (
	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/compute"
)

// ESX generation.
func GenerateESX(numESX int) []zebra.Resource {
	ESXarr := make([]zebra.Resource, 0, numESX)

	for i := 0; i < numESX; i++ {
		name := Name()

		serverID := SelectServer()

		labels := CreateLabels()
		ip := RandIP()

		esx := compute.NewESX(name, serverID, ip, labels)

		if esx.Labels.Validate() != nil {
			esx.Labels = GroupLabels(esx.Labels, GroupVal(esx))
		}

		ESXarr = append(ESXarr, esx)
	}

	return ESXarr
}
