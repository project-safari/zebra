package pkg

import "net"

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

// create IP arr.
func CreateIPArr(ipNum int) []net.IPNet {
	var nets net.IPNet // masked below.

	netArr := []net.IPNet{}

	SampleIPAddr := IPsamples()

	eight := 8

	for i := 0; i < ipNum; i++ {
		ip := RandData(SampleIPAddr)
		nets.IP = net.IP(ip)

		masked := net.IPv6len * eight
		what := net.IPv4len * eight

		nets.Mask = net.CIDRMask(masked, what)
		netArr = append(netArr, nets)
	}

	return netArr
}

// this array will help set possible sample IP addresses.
func IPsamples() []string {
	SampleIPAddr := []string{
		"192.332.11.05", "192.232.11.37", "192.232.22.05", "192.225.11.05",
		"192.0.0.0", "192.192.192.192", "225.225.225.225", "192.192.64.08",
	}

	return SampleIPAddr
}
