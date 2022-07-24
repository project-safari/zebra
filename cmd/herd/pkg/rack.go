package pkg

import (
	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/dc"
)

// generate rack.
func GenerateRack(numRacks int) []zebra.Resource {
	racks := make([]zebra.Resource, 0, numRacks)
	rows := Rows()

	for i := 0; i < numRacks; i++ {
		name := Name()
		labels := CreateLabels()

		rack := dc.NewRack(name, rows, labels)

		if rack.Labels.Validate() != nil {
			rack.Labels = GroupLabels(rack.Labels, GroupVal(rack))
		}

		racks = append(racks, rack)
	}

	return racks
}
