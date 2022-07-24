package pkg

import (
	// some imports from zebra.
	"github.com/project-safari/zebra"
	// some imports from compute.
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

		if server.Labels.Validate() != nil {
			server.Labels = GroupLabels(server.Labels, GroupVal(server))
		}

		servers = append(servers, server)
	}

	return servers
}
