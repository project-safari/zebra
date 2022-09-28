package dc

import (
	"fmt"

	"github.com/project-safari/zebra"
)

// Function that generates "mock" dcs as sample data.
//
// It takes in the number of resources to generate and returns a list of zebra resources.
func MockDC(num int) []zebra.Resource {
	rs := make([]zebra.Resource, 0, num)

	for i := 1; i <= num; i++ {
		s := NewDatacenter(
			fmt.Sprintf("101%d Hollywood Blvd, LA, CA", i),
			fmt.Sprintf("mock-dc-%d", i),
			"mocker",
			"dc",
		)

		rs = append(rs, s)
	}

	return rs
}

// Function that generates "mock" labs as sample data.
//
// It takes in the number of resources to generate and returns a list of zebra resources.
func MockLab(num int) []zebra.Resource {
	rs := make([]zebra.Resource, 0, num)

	for i := 1; i <= num; i++ {
		s := NewLab(
			fmt.Sprintf("mock-lab-%d", i),
			"mocker",
			"lab",
		)

		rs = append(rs, s)
	}

	return rs
}

// Function that generates "mock" racks as sample data.
//
// It takes in the number of resources to generate and returns a list of zebra resources.
func MockRack(num int) []zebra.Resource {
	maxRow := 10
	rs := make([]zebra.Resource, 0, num)

	for i := 1; i <= num; i++ {
		s := NewRack(
			fmt.Sprintf("mock-row-%d", i%maxRow),
			fmt.Sprintf("mock-lab-%d", i),
			"mocker",
			"rack",
		)

		rs = append(rs, s)
	}

	return rs
}
