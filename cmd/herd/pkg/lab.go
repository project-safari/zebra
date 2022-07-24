package pkg

import (
	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/dc"
)

func GenerateLab(numLabs int) []zebra.Resource {
	labs := make([]zebra.Resource, 0, numLabs)

	for i := 0; i < numLabs; i++ {
		name := Name()
		labels := CreateLabels()

		lab := dc.NewLab(name, labels)

		if lab.Labels.Validate() != nil {
			lab.Labels = GroupLabels(lab.Labels, GroupVal(lab))
		}

		labs = append(labs, lab)
	}

	return labs
}
