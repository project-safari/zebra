package pkg

import (
	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/dc"
)

// Function for lab generation.
//
// Returns an array of type zebra.Resource.
func GenerateLab(numLabs int) []zebra.Resource {
	labs := make([]zebra.Resource, 0, numLabs)

	for i := 0; i < numLabs; i++ {
		name := Name()
		labels := CreateLabels()

		lab := dc.NewLab(name, labels)

		if lab.LabelsValidate() != nil {
			lab.Labels = GroupLabels(lab.Labels, GroupVal(lab))
		}

		labs = append(labs, lab)
	}

	return labs
}
