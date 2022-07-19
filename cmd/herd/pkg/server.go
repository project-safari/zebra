package pkg

import (
	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/compute"
)

// generate server resources.
func GenerateServer(numServers int) []zebra.Resource {
	servers := make([]zebra.Resource, 0, numServers)

	for i := 0; i < numServers; i++ {
		arr := []string{Serials(), Models(), Name()}

		labels := CreateLabels()
		ip := RandIP()

		server := compute.NewServer(arr, ip, labels)

		servers = append(servers, server)
	}

	return servers
}
