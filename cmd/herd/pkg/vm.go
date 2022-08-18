package pkg

import (
	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/compute"
)

// Function to generate VM resources.
//
// Returns an array of type zebra.Resource.
func GenerateVM(numVM int) []zebra.Resource {
	VMarr := make([]zebra.Resource, 0, numVM)

	for i := 0; i < numVM; i++ {
		arr := []string{Name(), SelectESX(), SelectVcenter()}

		labels := CreateLabels()
		ip := RandIP()

		VM := compute.NewVM(arr, ip, labels)

		if VM.LabelsValidate() != nil {
			VM.Labels = GroupLabels(VM.Labels, GroupVal(VM))
		}

		VMarr = append(VMarr, VM)
	}

	return VMarr
}
