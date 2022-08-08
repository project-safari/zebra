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
	usr := new(auth.User)

	configFile := cmd.Flag("config").Value.String()

	manyUsr := map[string]*auth.User{}

	manyUsr[usr.Name] = usr

	config, e := Load(configFile)

	if e != nil {
		return e
	}

	if _, e := GetPath(config, GetType(args), manyUsr); e != nil {
		return e
	}

	fmt.Println(printUser(manyUsr).Render())

	return nil
}

func ShowReg(cmd *cobra.Command, args []string) error {
	manyUsr := map[string]*auth.User{}
	usr := new(auth.User)

	configFile := cmd.Flag("config").Value.String()
	config, e := Load(configFile)

	manyUsr[usr.Name] = usr

	if e != nil {
		return e
	}

	if _, e := GetPath(config, GetType(args), manyUsr); e != nil {
		return e
	}

	fmt.Println(printUser(manyUsr).Render())

	return nil
}

// network resources.
func ShowVlan(cmd *cobra.Command, args []string) error {
	vlans := map[string]*network.VLANPool{}
	netName := args[0]

	vlan := new(network.VLANPool)

	vlans[netName] = vlan

	configFile := cmd.Flag("config").Value.String()
	config, e := Load(configFile)

	if e != nil {
		return e
	}

	if _, e := GetPath(config, GetType(args), vlans); e != nil {
		return e
	}

	fmt.Println(printNets(vlans).Render())

	return nil
}

func ShowSw(cmd *cobra.Command, args []string) error {
	manySw := map[string]*network.Switch{}

	configFile := cmd.Flag("config").Value.String()

	swName := args[0]
	sw := new(network.Switch)

	config, e := Load(configFile)

	if e != nil {
		return e
	}

	manySw[swName] = sw

	if _, e := GetPath(config, GetType(args), manySw); e != nil {
		return e
	}

	fmt.Println(printSwitch(manySw).Render())

	return nil
}

func ShowIP(cmd *cobra.Command, args []string) error {
	configFile := cmd.Flag("config").Value.String()

	pools := map[string]*network.IPAddressPool{}
	addr := new(network.IPAddressPool)

	config, e := Load(configFile)

	IPName := args[0]

	pools[IPName] = addr

	if e != nil {
		return e
	}

	// pass args
	if _, e := GetPath(config, GetType(args), pools); e != nil {
		return e
	}

	fmt.Println(printIP(pools).Render())

	return nil
}

// datacenter.
func ShowDC(cmd *cobra.Command, args []string) error {
	center := new(dc.Datacenter)

	manyCenters := map[string]*dc.Datacenter{}

	configFile := cmd.Flag("config").Value.String()
	config, e := Load(configFile)

	centName := args[0]

	manyCenters[centName] = center

	if e != nil {
		return e
	}

	if _, e := GetPath(config, GetType(args), manyCenters); e != nil {
		return e
	}

	fmt.Println(printDC(manyCenters).Render())

	return nil
}

func ShowLab(cmd *cobra.Command, args []string) error {
	manyLabs := map[string]*dc.Lab{}

	labName := args[0]

	configFile := cmd.Flag("config").Value.String()
	config, e := Load(configFile)

	lab := new(dc.Lab)

	if e != nil {
		return e
	}

	manyLabs[labName] = lab

	if _, e := GetPath(config, GetType(args), manyLabs); e != nil {
		return e
	}

	fmt.Println(printLab(manyLabs).Render())

	return nil
}

func ShowRack(cmd *cobra.Command, args []string) error {
	configFile := cmd.Flag("config").Value.String()

	config, e := Load(configFile)

	if e != nil {
		return e
	}

	vcName := args[0]

	manyRacks := map[string]*dc.Rack{}

	rack := new(dc.Rack)

	manyRacks[vcName] = rack

	if _, e := GetPath(config, GetType(args), manyRacks); e != nil {
		return e
	}

	fmt.Println(printRack(manyRacks).Render())

	return nil
}

// server.
func ShowServ(cmd *cobra.Command, args []string) error {
	srvName := args[0]

	configFile := cmd.Flag("config").Value.String()
	config, e := Load(configFile)

	if e != nil {
		return e
	}

	srv := new(compute.Server)

	manySrv := map[string]*compute.Server{}

	manySrv[srvName] = srv

	if _, e := GetPath(config, GetType(args), manySrv); e != nil {
		return e
	}

	fmt.Println(printServer(manySrv).Render())

	return nil
}

func ShowESX(cmd *cobra.Command, args []string) error {
	config, e := Load(cmd.Flag("config").Value.String())

	if e != nil {
		return e
	}

	esxName := args[0]
	manyESX := map[string]*compute.ESX{}

	esx := new(compute.ESX)

	manyESX[esxName] = esx

	if _, e := GetPath(config, GetType(args), manyESX); e != nil {
		return e
	}

	fmt.Println(printESX(manyESX).Render())

	return nil
}

func ShowVC(cmd *cobra.Command, args []string) error {
	config, e := Load(cmd.Flag("config").Value.String())

	if e != nil {
		return e
	}

	vc := new(compute.VCenter)

	manyVC := map[string]*compute.VCenter{}

	vcName := args[0]

	manyVC[vcName] = vc

	if _, e := GetPath(config, GetType(args), manyVC); e != nil {
		return e
	}

	fmt.Println(printVC(manyVC).Render())

	return nil
}

func ShowVM(cmd *cobra.Command, args []string) error {
	config, e := Load(cmd.Flag("config").Value.String())
	vcName := args[0]

	vm := new(compute.VM)

	manyVM := map[string]*compute.VM{}

	if e != nil {
		return e
	}

	manyVM[vcName] = vm

	if _, e := GetPath(config, GetType(args), manyVM); e != nil {
		return e
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
