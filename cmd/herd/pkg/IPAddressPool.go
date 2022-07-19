package pkg

import "github.com/project-safari/zebra/network"

// generate IP pools.
func GenerateIPPool(numAddr int) []*network.IPAddressPool {
	IPpool := make([]*network.IPAddressPool, 0, numAddr)

	for i := 1; i < numAddr; i++ {
		labels := CreateLabels()
		ipNum := RandNum(numAddr)

		ipArr := CreateIPArr(int(ipNum))

		IPaddr := network.NewIPAddressPool(ipArr, labels)

		IPpool = append(IPpool, IPaddr)
	}

	return IPpool
}
