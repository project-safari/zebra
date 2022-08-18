package pkg

import (
	// Some imports from zebra.
	"github.com/project-safari/zebra"
	// Some imports from compute.
	"github.com/project-safari/zebra/compute"
)

// Function to generate server resources.
//
// Returns an array of type zebra.Resource.
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
