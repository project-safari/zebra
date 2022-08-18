package pkg

import (
	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/dc"
)

// Function to generate datacenter info.
//
// Returns an array of type zebra.Resource.
func GenerateDatacenter(numDC int) []zebra.Resource {
	datacent := make([]zebra.Resource, 0, numDC)

	for i := 0; i < numDC; i++ {
		name := Name()
		labels := CreateLabels()
		location := Addresses()

		data := dc.NewDatacenter(location, name, labels)

		if data.LabelsValidate() != nil {
			data.Labels = GroupLabels(data.Labels, GroupVal(data))
		}

		datacent = append(datacent, data)
	}

	return datacent
}
