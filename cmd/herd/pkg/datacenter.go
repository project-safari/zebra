package pkg

import "github.com/project-safari/zebra/dc"

// generate datacenter info.
func GenerateDatacenter(numDC int) []*dc.Datacenter {
	datacent := make([]*dc.Datacenter, 0, numDC)

	for i := 0; i < numDC; i++ {
		name := Name()
		labels := CreateLabels()
		location := Addresses()

		data := dc.NewDatacenter(location, name, labels)

		datacent = append(datacent, data)
	}

	return datacent
}
