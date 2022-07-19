package pkg

import (
	"net"

	"github.com/project-safari/zebra/compute"
)

// generate server resources.
func GenerateServer(numServers int) []*compute.Server {
	servers := make([]*compute.Server, 0, numServers)

	for i := 0; i < numServers; i++ {
		arr := []string{Serials(), Models(), Name()}

		labels := CreateLabels()
		ip := net.IP(RandData(IPsamples()))

		server := compute.NewServer(arr, ip, labels)

		servers = append(servers, server)
	}

	return servers
}
