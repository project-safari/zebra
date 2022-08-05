package main

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/project-safari/zebra/auth"
	"github.com/project-safari/zebra/compute"
	"github.com/project-safari/zebra/dc"
	"github.com/project-safari/zebra/network"
	"github.com/spf13/cobra"
)

// user info.
func ShowUsr(cmd *cobra.Command, args []string) error {
	configFile := cmd.Flag("config").Value.String()
	config, e := Load(configFile)

	if e != nil {
		return e
	}

	manyUsr := map[string]*auth.User{}
	usr := new(auth.User)

	p := fmt.Sprintf("/login/%s", args[0])

	if len(args) == 0 {
		if _, e := GetPath(config, "/login", GetType("User")); e != nil {
			return e
		}
	} else {
		if _, e := GetPath(config, p, usr); e != nil {
			return e
		}

		manyUsr[usr.Name] = usr
	}

	// cannot use manyUsr (variable of type *zebra.ResourceMap)
	// as map[string]*auth.User value in argument to
	fmt.Println(printUser(manyUsr).Render())

	return nil
}

func ShowReg(cmd *cobra.Command, args []string) error {
	configFile := cmd.Flag("config").Value.String()
	config, e := Load(configFile)

	manyUsr := map[string]*auth.User{}
	usr := new(auth.User)

	p := fmt.Sprintf("/register/%s", args[0])

	if e != nil {
		return e
	}

	if len(args) == 0 {
		if _, e := GetPath(config, "/register", GetType("User")); e != nil {
			return e
		}
	} else {
		if _, e := GetPath(config, p, usr); e != nil {
			return e
		}

		manyUsr[usr.Name] = usr
	}

	fmt.Println(printUser(manyUsr).Render())

	return nil
}

// network resources.
func ShowVlan(cmd *cobra.Command, args []string) error {
	configFile := cmd.Flag("config").Value.String()
	config, e := Load(configFile)

	p := fmt.Sprintf("/refresh/%s", args[0])

	vlans := map[string]*network.VLANPool{}
	netName := args[0]

	vlan := new(network.VLANPool)

	if e != nil {
		return e
	}

	if len(args) == 0 {
		if _, e := GetPath(config, "/refresh", GetType("VLANPool")); e != nil {
			return e
		}
	} else {
		if _, e := GetPath(config, p, vlan); e != nil {
			return e
		}

		vlans[netName] = vlan
	}

	fmt.Println(printNets(vlans).Render())

	return nil
}

func ShowSw(cmd *cobra.Command, args []string) error {
	configFile := cmd.Flag("config").Value.String()
	config, e := Load(configFile)

	swName := args[0]
	sw := new(network.Switch)

	p := fmt.Sprintf("/refresh/%s", args[0])
	manySw := map[string]*network.Switch{}

	if e != nil {
		return e
	}

	if len(args) == 0 {
		if _, e := GetPath(config, "/refresh", GetType("Switch")); e != nil {
			return e
		}
	} else {
		if _, e := GetPath(config, p, sw); e != nil {
			return e
		}

		manySw[swName] = sw
	}

	fmt.Println(printSwitch(manySw).Render())

	return nil
}

func ShowIP(cmd *cobra.Command, args []string) error {
	configFile := cmd.Flag("config").Value.String()
	config, e := Load(configFile)

	p := fmt.Sprintf("/refresh/%s", args[0])
	IPName := args[0]

	pools := map[string]*network.IPAddressPool{}
	addr := new(network.IPAddressPool)

	if e != nil {
		return e
	}

	if len(args) == 0 {
		if _, e := GetPath(config, "/refresh", GetType("IPAddressPool")); e != nil {
			return e
		}
	} else {
		if _, e := GetPath(config, p, addr); e != nil {
			return e
		}

		pools[IPName] = addr
	}

	fmt.Println(printIP(pools).Render())

	return nil
}

// datacenter.
func ShowDC(cmd *cobra.Command, args []string) error {
	configFile := cmd.Flag("config").Value.String()
	config, e := Load(configFile)

	p := fmt.Sprintf("/refresh/%s", args[0])

	center := new(dc.Datacenter)
	centName := args[0]

	manyCenters := map[string]*dc.Datacenter{}

	if e != nil {
		return e
	}

	if len(args) == 0 {
		if _, e := GetPath(config, "/refresh", GetType("Datacenter")); e != nil {
			return e
		}
	} else {
		if _, e := GetPath(config, p, center); e != nil {
			return e
		}

		manyCenters[centName] = center
	}

	fmt.Println(printDC(manyCenters).Render())

	return nil
}

func ShowLab(cmd *cobra.Command, args []string) error {
	p := fmt.Sprintf("/refresh/%s", args[0])
	labName := args[0]

	configFile := cmd.Flag("config").Value.String()
	config, e := Load(configFile)

	manyLabs := map[string]*dc.Lab{}

	lab := new(dc.Lab)

	if e != nil {
		return e
	}

	if len(args) == 0 {
		if _, e := GetPath(config, "/refresh", GetType("Lab")); e != nil {
			return e
		}
	} else {
		if _, e := GetPath(config, p, lab); e != nil {
			return e
		}

		manyLabs[labName] = lab
	}

	fmt.Println(printLab(manyLabs).Render())

	return nil
}

func ShowRack(cmd *cobra.Command, args []string) error {
	configFile := cmd.Flag("config").Value.String()
	p := fmt.Sprintf("/refresh/%s", args[0])

	config, e := Load(configFile)

	if e != nil {
		return e
	}

	vcName := args[0]

	manyRacks := map[string]*dc.Rack{}

	rack := new(dc.Rack)

	if len(args) == 0 {
		if _, e := GetPath(config, "/refresh", GetType("Rack")); e != nil {
			return e
		}
	} else {
		if _, e := GetPath(config, p, rack); e != nil {
			return e
		}

		manyRacks[vcName] = rack
	}

	fmt.Println(printRack(manyRacks).Render())

	return nil
}

// server.
func ShowServ(cmd *cobra.Command, args []string) error {
	configFile := cmd.Flag("config").Value.String()
	config, e := Load(configFile)

	p := fmt.Sprintf("/refresh/%s", args[0])

	if e != nil {
		return e
	}

	srvName := args[0]
	srv := new(compute.Server)

	manySrv := map[string]*compute.Server{}

	if len(args) == 0 {
		if _, e := GetPath(config, "/refresh", GetType("Server")); e != nil {
			return e
		}
	} else {
		if _, e := GetPath(config, p, srv); e != nil {
			return e
		}

		manySrv[srvName] = srv
	}

	fmt.Println(printServer(manySrv).Render())

	return nil
}

func ShowESX(cmd *cobra.Command, args []string) error {
	config, e := Load(cmd.Flag("config").Value.String())

	p := fmt.Sprintf("/refresh/%s", args[0])

	if e != nil {
		return e
	}

	esxName := args[0]
	manyESX := map[string]*compute.ESX{}

	esx := new(compute.ESX)

	if len(args) == 0 {
		if _, e := GetPath(config, "/refresh", GetType("ESX")); e != nil {
			return e
		}
	} else {
		if _, e := GetPath(config, p, esx); e != nil {
			return e
		}

		manyESX[esxName] = esx
	}

	fmt.Println(printESX(manyESX).Render())

	return nil
}

func ShowVC(cmd *cobra.Command, args []string) error {
	p := fmt.Sprintf("/refresh/%s", args[0])

	config, e := Load(cmd.Flag("config").Value.String())

	if e != nil {
		return e
	}

	vc := new(compute.VCenter)

	vcName := args[0]
	manyVC := map[string]*compute.VCenter{}

	if len(args) == 0 {
		if _, e := GetPath(config, "/refresh", GetType("VCenter")); e != nil {
			return e
		}
	} else {
		if _, e := GetPath(config, p, vc); e != nil {
			return e
		}

		manyVC[vcName] = vc
	}

	fmt.Println(printVC(manyVC).Render())

	return nil
}

func ShowVM(cmd *cobra.Command, args []string) error {
	config, e := Load(cmd.Flag("config").Value.String())
	vcName := args[0]

	vm := new(compute.VM)

	manyVM := map[string]*compute.VM{}

	p := fmt.Sprintf("/refresh/%s", args[0])

	if e != nil {
		return e
	}

	if len(args) == 0 {
		if _, e := GetPath(config, "/refresh", GetType("VM")); e != nil {
			return e
		}
	} else {
		if _, e := GetPath(config, p, vm); e != nil {
			return e
		}

		manyVM[vcName] = vm
	}

	fmt.Println(printVM(manyVM).Render())

	return nil
}

func printSwitch(srv map[string]*network.Switch) table.Writer {
	data := table.NewWriter()
	data.AppendHeader(table.Row{"ID", "Management IP", "Credentials", "Serial Number", "Model", "Ports", "Labels"})

	for piece, sw := range srv {
		data.AppendRow(table.Row{
			piece,
			fmt.Sprintf(sw.ID),
			sw.ManagementIP.String(),

			fmt.Sprintf("%s", sw.Credentials.Keys),
			fmt.Sprintf(sw.SerialNumber),
			fmt.Sprintf(sw.Model),

			fmt.Sprintf("%010d", sw.NumPorts),
			fmt.Sprintf("%s", sw.Labels),
		})
	}

	return data
}

func printLab(labs map[string]*dc.Lab) table.Writer {
	data := table.NewWriter()
	data.AppendHeader(table.Row{"ID", "Name", "Labels"})

	for piece, lb := range labs {
		data.AppendRow(table.Row{
			piece,
			fmt.Sprintf(lb.NamedResource.ID),
			fmt.Sprintf(lb.NamedResource.Name),

			fmt.Sprintf("%s", lb.NamedResource.Labels),
		})
	}

	return data
}

func printDC(dcs map[string]*dc.Datacenter) table.Writer {
	data := table.NewWriter()
	data.AppendHeader(table.Row{"ID", "Name", "Address", "Labels"})

	for piece, dc := range dcs {
		data.AppendRow(table.Row{
			piece,
			fmt.Sprintf(dc.NamedResource.ID),
			fmt.Sprintf(dc.NamedResource.Name),

			fmt.Sprintf(dc.Address),
			fmt.Sprintf("%s", dc.NamedResource.Labels),
		})
	}

	return data
}

func printServer(servers map[string]*compute.Server) table.Writer {
	data := table.NewWriter()
	data.AppendHeader(table.Row{"ID", "Name", "Board IP", "Model", "Credentials", "Labels"})

	for piece, s := range servers {
		data.AppendRow(table.Row{
			piece,
			fmt.Sprintf(s.NamedResource.ID),
			fmt.Sprintf(s.NamedResource.Name),
			s.BoardIP.String(),

			fmt.Sprintf(s.Model),

			fmt.Sprintf("%s", s.Credentials.Keys),
			fmt.Sprintf("%s", s.NamedResource.Labels),
		})
	}

	return data
}

func printESX(manyEsx map[string]*compute.ESX) table.Writer {
	data := table.NewWriter()
	data.AppendHeader(table.Row{"ID", "Name", "Server ID", "IP", "Credentials", "Labels"})

	for piece, esx := range manyEsx {
		data.AppendRow(table.Row{
			piece,
			fmt.Sprintf(esx.NamedResource.ID),
			fmt.Sprintf(esx.NamedResource.Name),

			fmt.Sprintf(esx.ServerID),
			esx.IP.String(),

			fmt.Sprintf("%s", esx.Credentials.Keys),
			fmt.Sprintf("%s", esx.NamedResource.Labels),
		})
	}

	return data
}

func printVC(manyVC map[string]*compute.VCenter) table.Writer {
	data := table.NewWriter()
	data.AppendHeader(table.Row{"ID", "Name", "IP", "Credentials", "Labels"})

	for piece, vc := range manyVC {
		data.AppendRow(table.Row{
			piece,
			fmt.Sprintf(vc.NamedResource.ID),

			fmt.Sprintf(vc.NamedResource.Name),
			vc.IP.String(),

			fmt.Sprintf("%s", vc.Credentials.Keys),
			fmt.Sprintf("%s", vc.NamedResource.Labels),
		})
	}

	return data
}

func printVM(manyVM map[string]*compute.VM) table.Writer {
	data := table.NewWriter()
	data.AppendHeader(table.Row{
		"ID", "Name", "IP", "Credentials",
		"ESXID", "VCenter ID", "Management IP", "Labels",
	})

	for piece, vm := range manyVM {
		data.AppendRow(table.Row{
			piece,
			fmt.Sprintf(vm.NamedResource.ID),

			fmt.Sprintf(vm.NamedResource.Name),
			fmt.Sprintf("%s", vm.Credentials.Keys),
			fmt.Sprintf(vm.ESXID),
			fmt.Sprintf(vm.VCenterID),
			vm.ManagementIP.String(),
			fmt.Sprintf("%s", vm.NamedResource.Labels),
		})
	}

	return data
}

func printRack(racks map[string]*dc.Rack) table.Writer {
	data := table.NewWriter()
	data.AppendHeader(table.Row{"ID", "Name", "Status", "Row", "Labels"})

	for piece, rack := range racks {
		data.AppendRow(table.Row{
			piece,
			fmt.Sprintf(rack.NamedResource.ID),
			fmt.Sprintf(rack.NamedResource.Name),

			fmt.Sprintf(rack.Row),
			fmt.Sprintf("%s", rack.NamedResource.Labels),
		})
	}

	return data
}

func printUser(users map[string]*auth.User) table.Writer {
	data := table.NewWriter()
	data.AppendHeader(table.Row{"ID", "Name", "Status", "Password Hash", "Role", "Priviledges", "Labels"})

	for piece, user := range users {
		data.AppendRow(table.Row{
			piece,
			fmt.Sprintf(user.NamedResource.ID),
			fmt.Sprintf(user.NamedResource.Name),

			fmt.Sprintf(user.PasswordHash),
			fmt.Sprintf(user.Role.Name),

			fmt.Sprintf("%s", user.Role.Privileges),
			fmt.Sprintf("%s", user.NamedResource.Labels),
		})
	}

	return data
}

func printNets(vlans map[string]*network.VLANPool) table.Writer {
	data := table.NewWriter()
	data.AppendHeader(table.Row{"ID", "Status", "Range Start", "Range End", "Labels"})

	for piece, vlan := range vlans {
		data.AppendRow(table.Row{
			piece,
			fmt.Sprintf(vlan.ID),
			fmt.Sprintf(vlan.Status.UsedBy),

			fmt.Sprintf("%010d", vlan.RangeStart),
			fmt.Sprintf("%010d", vlan.RangeEnd),
			fmt.Sprintf("%s", vlan.Labels),
		})
	}

	return data
}

func printIP(vlans map[string]*network.IPAddressPool) table.Writer {
	data := table.NewWriter()
	data.AppendHeader(table.Row{"ID", "Status", "Subnets", "Labels"})

	for piece, pool := range vlans {
		data.AppendRow(table.Row{
			piece,
			fmt.Sprintf(pool.ID),
			fmt.Sprintf(pool.Status.UsedBy),

			fmt.Sprintf("%s", pool.Subnets),
			fmt.Sprintf("%s", pool.Labels),
		})
	}

	return data
}
