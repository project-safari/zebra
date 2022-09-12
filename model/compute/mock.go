package compute

import (
	"fmt"
	"net"

	"github.com/project-safari/zebra"
)

func MockServer(num int) []zebra.Resource {
	models := []string{"SERVER-MODEL-1", "SERVER-MODEL-2", "SERVER-MODEL-3"}
	rs := make([]zebra.Resource, 0, num)

	for i := 1; i <= num; i++ {
		s := NewServer(
			fmt.Sprintf("SERVER-SERIAL-%d", i),
			models[i%3],
			fmt.Sprintf("mock-server-%d", i),
			"mocker",
			"server",
		)
		s.BoardIP = net.IP{10, 10, 10, byte(i)}
		s.Credentials = zebra.NewCredentials("admin")
		_ = s.Credentials.Add("password", fmt.Sprintf("ServerSecret%d!!!", i))

		rs = append(rs, s)
	}

	return rs
}

func MockESX(num int) []zebra.Resource {
	rs := make([]zebra.Resource, 0, num)

	for i := 1; i <= num; i++ {
		s := NewESX(
			fmt.Sprintf("mock-server-%d", i),
			fmt.Sprintf("mock-esx-%d", i),
			"mocker",
			"esx",
		)
		s.IP = net.IP{11, 11, 11, byte(i)}
		s.Credentials = zebra.NewCredentials("admin")
		_ = s.Credentials.Add("password", fmt.Sprintf("ESXSecret%d!!!", i))

		rs = append(rs, s)
	}

	return rs
}

func MockVCenter(num int) []zebra.Resource {
	rs := make([]zebra.Resource, 0, num)

	for i := 1; i <= num; i++ {
		s := NewVCenter(
			fmt.Sprintf("mock-vcenter-%d", i),
			"mocker",
			"vcenter",
		)
		s.IP = net.IP{12, 12, 12, byte(i)}
		s.Credentials = zebra.NewCredentials("admin")
		_ = s.Credentials.Add("password", fmt.Sprintf("VCenterSecret%d!!!", i))

		rs = append(rs, s)
	}

	return rs
}

func MockVM(num int) []zebra.Resource {
	rs := make([]zebra.Resource, 0, num)

	for i := 1; i <= num; i++ {
		s := NewVM(
			fmt.Sprintf("mock-esx-%d", i),
			fmt.Sprintf("mock-vm-%d", i),
			"mocker",
			"vm",
		)
		s.ManagementIP = net.IP{13, 13, 13, byte(i)}
		s.Credentials = zebra.NewCredentials("admin")
		_ = s.Credentials.Add("password", fmt.Sprintf("VMSecret%d!!!", i))

		rs = append(rs, s)
	}

	return rs
}
