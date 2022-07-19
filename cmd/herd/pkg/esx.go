package pkg

import (
	"net"

	"github.com/project-safari/zebra/compute"
)

// ESX generation.
func GenerateESX(numESX int) []*compute.ESX {
	ESXarr := make([]*compute.ESX, 0, numESX)

	for i := 0; i < numESX; i++ {
		name := Name()

		serverID := SelectServer()

		labels := CreateLabels()
		ip := net.IP(RandData(IPsamples()))

		esx := compute.NewESX(name, serverID, ip, labels)

		ESXarr = append(ESXarr, esx)
	}

	return ESXarr
}
