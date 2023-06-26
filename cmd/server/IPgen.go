package main

import (
	"crypto/rand"
	"net"
)

var ipInUse []net.IP

func contains(all []net.IP, one net.IP) bool {
	for _, each := range all {
		eachStr := string(each)
		str := string(one)
		if eachStr == str {
			return true
		}
	}

	return false
}

func generateIP(cidr string) (net.IP, []net.IP, error) {
GENERATE:

	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, nil, err
	}

	ones, _ := ipnet.Mask.Size()
	quotient := ones / 8
	remainder := ones % 8

	r := make([]byte, 4)
	rand.Read(r)

	for i := 0; i <= quotient; i++ {
		if i == quotient {
			shifted := byte(r[i]) >> remainder
			r[i] = ^ipnet.IP[i] & shifted
		} else {
			r[i] = ipnet.IP[i]
		}
	}
	ip = net.IPv4(r[0], r[1], r[2], r[3])

	if contains(ipInUse, ip) {
		goto GENERATE
	} else {
		ipInUse = append(ipInUse, ip)
	}

	if ip.Equal(ipnet.IP) {
		goto GENERATE
	}

	return ip, ipInUse, nil
}
