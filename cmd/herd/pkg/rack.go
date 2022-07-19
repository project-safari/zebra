package pkg

import "github.com/project-safari/zebra/dc"

// generate rack.
func GenerateRack(numRacks int) []*dc.Rack {
	racks := make([]*dc.Rack, 0, numRacks)
	rows := Rows()

	for i := 0; i < numRacks; i++ {
		name := Name()
		labels := CreateLabels()

		rack := dc.NewRack(name, rows, labels)

		racks = append(racks, rack)
	}

	return racks
}
