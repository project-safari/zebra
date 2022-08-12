package pkg

import (
	// some imports from zebra.
	"github.com/project-safari/zebra"
	// some imports from compute.
	"github.com/project-safari/zebra/compute"
)

// function to generate server resources.
//
// returns an array of type zebra.Resource.
func GenerateServer(numServers int) []zebra.Resource {
	servers := make([]zebra.Resource, 0, numServers)

	for i := 0; i < numServers; i++ {
		arr := []string{Serials(), Models(), Name()}

		labels := CreateLabels()
		ip := RandIP()

		server := compute.NewServer(arr, ip, labels)

		if server.LabelsValidate() != nil {
			server.Labels = GroupLabels(server.Labels, GroupVal(server))
		}

		servers = append(servers, server)
	}

	return servers
}
