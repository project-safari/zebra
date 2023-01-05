//nolint:gomnd, goconst
package migration

import (
	"fmt"
	"net"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/model/compute"
	"github.com/project-safari/zebra/model/network"
)

func checkIPAddress(ip string) error {
	if ip == "" || net.ParseIP(ip) == nil || ip == "<nil>" || ip == "N/A" {
		return network.ErrIPEmpty
	}

	return nil
}

func serverFiller(rt Racktables) zebra.Resource {
	theIP := net.IP("1.1.1.1")

	user := ""

	if rt.Owner == "" {
		user = "admin"
	} else {
		user = rt.Owner
	}

	if checkIPAddress(rt.IP) == nil {
		resIP := net.ParseIP(rt.IP)

		if resIP != nil {
			theIP = resIP
		} else {
			theIP = net.ParseIP("127.0.0.1") // some default IP since db is empty.
		}
	}

	s := compute.NewServer(
		"DB-SERVER",
		"db-model",
		rt.Name,
		"server",
		"system.group-server"+"-"+rt.Group,
	)

	s.BoardIP = theIP
	s.Credentials = zebra.NewCredentials(user)
	_ = s.Credentials.Add("password", fmt.Sprintf(user, 123))

	return s
}

func esxFiller(rt Racktables) zebra.Resource {
	user := ""

	theIP := net.IP("1.1.1.1")

	if checkIPAddress(rt.IP) == nil {
		resIP := net.ParseIP(rt.IP)

		if resIP != nil {
			theIP = resIP
		} else {
			theIP = net.ParseIP("127.0.0.1")
		}
	}

	if rt.Owner == "" {
		user = "admin"
	} else {
		user = rt.Owner
	}

	s := compute.NewESX(
		rt.ID,
		rt.Name,
		user,
		"system.group-server-esx"+"-"+rt.Group,
	)

	s.IP = theIP
	s.Credentials = zebra.NewCredentials(user)
	_ = s.Credentials.Add("password", fmt.Sprintf(user, 123))

	return s
}

func vcenterFiller(rt Racktables) zebra.Resource {
	if rt.Owner == "" {
		rt.Owner = "admin"
	}

	theIP := net.IP("1.1.1.1")

	user := rt.Owner

	if checkIPAddress(rt.IP) == nil {
		resIP := net.ParseIP(rt.IP)

		if resIP != nil {
			theIP = resIP
		} else {
			theIP = net.ParseIP("127.0.0.1")
		}
	}

	s := compute.NewVCenter(
		fmt.Sprintf(rt.Name),
		user,
		"system.group-server-vcenter",
	)

	s.IP = theIP
	s.Credentials = zebra.NewCredentials(user)
	_ = s.Credentials.Add("password", fmt.Sprintf(user, 123))

	return s
}

func vmFiller(rt Racktables) zebra.Resource {
	res := compute.NewVM("DB-esx", rt.Name, rt.Owner, "system.group-server-vcenter-vm"+"-"+rt.Group)

	if rt.Owner == "" {
		rt.Owner = "admin"
	}

	user := rt.Owner

	if checkIPAddress(rt.IP) == nil {
		resIP := net.ParseIP(rt.IP)

		if resIP != nil {
			res.ManagementIP = resIP
		} else {
			res.ManagementIP = net.ParseIP("127.0. 0.1")
		}
	}

	res.Credentials = zebra.NewCredentials(user)
	_ = res.Credentials.Add("password", fmt.Sprintf(user, 123))

	return res
}

func addressPoolFiller(rt Racktables) zebra.Resource {
	var resIP net.IP

	if checkIPAddress(rt.IP) == nil {
		resIP = net.ParseIP(rt.IP)

		if resIP == nil && net.IPMask(resIP) == nil {
			resIP = net.ParseIP("127.0. 0.1")
		}
	}

	if rt.Owner == "" {
		rt.Owner = "admin"
	}

	s := network.NewIPAddressPool(
		rt.Name,
		rt.Owner,
		"system.group-vlan-ipaddrpool"+"-"+rt.Group,
	)

	s.Subnets = []net.IPNet{{
		IP:   resIP,
		Mask: net.IPMask(resIP),
	}}

	return s
}

func switchFiller(rt Racktables) zebra.Resource {
	user := rt.Owner

	if user == "" { // if user is unspecified, the default user is admin.
		user = "admin"
	}

	iRes := network.NewSwitch(
		rt.Name,
		user,
		"system.group-vlan-switch"+"-"+rt.Group,
	)

	iRes.Model = "N/A" // curently no provision for this in the db.
	iRes.Meta.Type.Name = "network.switch"

	iRes.SerialNumber = "N/A" // curently no provision for this in the db.
	iRes.Credentials = zebra.NewCredentials(user)
	iRes.NumPorts = uint32(rt.Port)

	if checkIPAddress(rt.IP) == nil {
		resIP := net.ParseIP(rt.IP)

		if resIP != nil {
			iRes.ManagementIP = net.ParseIP(rt.IP)
		} else {
			iRes.ManagementIP = net.ParseIP("127.0. 0.1")
		}
	}

	return iRes
}

func vlanFiller(rt Racktables) zebra.Resource {
	r := network.NewVLANPool(
		rt.Name,
		rt.Owner,
		"system.group-ip-vlan"+"-"+rt.Group,
	)

	r.RangeStart = 0 // curently no provision for this in the db.
	someEnd := 100   // curently no provision for this in the db.
	r.RangeEnd = uint16(someEnd)

	return r
}

func DBData(n int) []zebra.Resource {
	RackArr := Do()

	resources := []zebra.Resource{}

	for i := 0; i < n; i++ {
		res := RackArr[i]

		eachRes, _, _ := createResFromData(res)

		resources = append(resources, eachRes)
	}

	return resources
}
