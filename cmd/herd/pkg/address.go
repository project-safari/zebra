package pkg

import (
	"fmt"
	"net"
)

// some random address.
func Addresses() string {
	possibleAdr := []string{
		"NYC", "Dallas", "Seattle", "Ottawa", "Paris",
		"London", "Athens", "Milan", "Philadelphia", "Ann Arbor",
		"DC", "Ankara", "Cape Verde", "LA", "Perth",
	}

	theAdr := RandData(possibleAdr)

	return theAdr
}

// create IP arr to use in generation of data for resources that need arrays of ip addresses.
// returns an array of net.IPNet.
func CreateIPArr(ipNum int) []net.IPNet {
	netArr := make([]net.IPNet, 0, ipNum)

	for i := 0; i < ipNum; i++ {
		_, subnet, _ := net.ParseCIDR(fmt.Sprintf("%d.%d.0.1/24", i, i))
		netArr = append(netArr, *subnet)
	}

	return netArr
}

// this function will help set possible sample IP addresses, returns a string array.
func IPsamples() []string {
	SampleIPAddr := []string{
		"1.1.1.1", "2.2.2.2", "3.3.3.3", "4.4.4.4",
		"5.5.5.5", "6.6.6.6", "7.7.7.7", "8.8.8.8",
	}

	return SampleIPAddr
}

// function to generate random ip addresses.
// returns a net.IP address.
func RandIP() net.IP {
	ipAddr := RandData(IPsamples())

	return net.ParseIP(ipAddr)
}
