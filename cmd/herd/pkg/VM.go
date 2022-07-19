package pkg

import (
	"net"

	"github.com/project-safari/zebra/compute"
)

// generate VM resources.
func GenerateVM(numVM int) []*compute.VM {
	VMarr := make([]*compute.VM, 0, numVM)

	for i := 0; i < numVM; i++ {
		arr := []string{Name(), SelectESX(), SelectVcenter()}

		labels := CreateLabels()
		ip := net.IP(RandData(IPsamples()))

		VM := compute.NewVM(arr, ip, labels)

		VMarr = append(VMarr, VM)
	}

	return VMarr
}
