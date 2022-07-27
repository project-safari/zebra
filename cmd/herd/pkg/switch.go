package pkg

import (
	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/network"
)

// generate switch resources.
func GenerateSwitch(numSwitch int) []zebra.Resource {
	port := Ports()

	switches := make([]zebra.Resource, 0, numSwitch)

	for i := 0; i < numSwitch; i++ {
		arr := []string{Serials(), Models(), Name()}
		labels := CreateLabels()
		ip := RandIP()
		sw := network.NewSwitch(arr, port, ip, labels)

		if sw.LabelsValidate() != nil {
			sw.Labels = GroupLabels(sw.Labels, GroupVal(sw))
		}

		switches = append(switches, sw)
	}

	return switches
}
