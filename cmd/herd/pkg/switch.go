package pkg

import (
	"net"

	"github.com/project-safari/zebra/network"
)

// generate switch resources.
func GenerateSwitch(numSwitch int) []*network.Switch {
	ip := net.IP(RandData(IPsamples()))
	port := Ports()

	switches := make([]*network.Switch, 0, numSwitch)

	for i := 0; i < numSwitch; i++ {
		arr := []string{Serials(), Models(), Name()}
		labels := CreateLabels()
		sw := network.NewSwitch(arr, port, ip, labels)

		switches = append(switches, sw)
	}

	return switches
}
