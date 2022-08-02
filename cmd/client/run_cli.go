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

	path := fmt.Sprintf("users/%s", args[0])

	client, err := NewClient(config)
	if err != nil {
		return err
	}

	if len(args) == 0 {
		if _, e := client.Get("users", nil, manyUsr); e != nil {
			return e
		}
	} else {
		if _, e := client.Get(path, nil, usr); e != nil {
			return e
		}

		manyUsr[usr.Name] = usr
	}

	fmt.Println(printUser(manyUsr).Render())

	return nil
}

func ShowReg(cmd *cobra.Command, args []string) error {
	configFile := cmd.Flag("config").Value.String()
	config, e := Load(configFile)

	if e != nil {
		return e
	}

	manyUsr := map[string]*auth.User{}
	usr := new(auth.User)

	path := fmt.Sprintf("registrations/%s", args[0])

	client, err := NewClient(config)
	if err != nil {
		return err
	}

	if len(args) == 0 {
		if _, e := client.Get("registrations", nil, manyUsr); e != nil {
			return e
		}
	} else {
		if _, e := client.Get(path, nil, usr); e != nil {
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

	if e != nil {
		return e
	}

	vlans := map[string]*network.VLANPool{}
	netName := args[0]
	vlan := new(network.VLANPool)

	path := fmt.Sprintf("vlans/%s", args[0])

	client, err := NewClient(config)
	if err != nil {
		return err
	}

	if len(args) == 0 {
		if _, e := client.Get("vlans", nil, vlans); e != nil {
			return e
		}
	} else {
		if _, e := client.Get(path, nil, vlan); e != nil {
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

	if e != nil {
		return e
	}

	swName := args[0]
	sw := new(network.Switch)

	path := fmt.Sprintf("switches/%s", args[0])
	manySw := map[string]*network.Switch{}

	client, err := NewClient(config)
	if err != nil {
		return err
	}

	if len(args) == 0 {
		if _, e := client.Get("switches", nil, manySw); e != nil {
			return e
		}
	} else {
		if _, e := client.Get(path, nil, sw); e != nil {
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

	if e != nil {
		return e
	}

	IPName := args[0]
	addr := new(network.IPAddressPool)

	pools := map[string]*network.IPAddressPool{}
	path := fmt.Sprintf("ip/%s", args[0])

	client, err := NewClient(config)
	if err != nil {
		return err
	}

	if len(args) == 0 {
		if _, e := client.Get("ip", nil, pools); e != nil {
			return e
		}
	} else {
		if _, e := client.Get(path, nil, addr); e != nil {
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

	path := fmt.Sprintf("datacenters/%s", args[0])

	if e != nil {
		return e
	}

	centName := args[0]
	center := new(dc.Datacenter)

	manyCenters := map[string]*dc.Datacenter{}

	client, err := NewClient(config)
	if err != nil {
		return err
	}

	if len(args) == 0 {
		if _, e := client.Get("datacenters", nil, manyCenters); e != nil {
			return e
		}
	} else {
		if _, e := client.Get(path, nil, center); e != nil {
			return e
		}

		manyCenters[centName] = center
	}

	fmt.Println(printDC(manyCenters).Render())

	return nil
}

func ShowLab(cmd *cobra.Command, args []string) error {
	configFile := cmd.Flag("config").Value.String()
	config, e := Load(configFile)

	path := fmt.Sprintf("labs/%s", args[0])

	if e != nil {
		return e
	}

	labName := args[0]
	lab := new(dc.Lab)
	manyLabs := map[string]*dc.Lab{}

	client, err := NewClient(config)
	if err != nil {
		return err
	}

	if len(args) == 0 {
		if _, e := client.Get("labs", nil, manyLabs); e != nil {
			return e
		}
	} else {
		if _, e := client.Get(path, nil, lab); e != nil {
			return e
		}

		manyLabs[labName] = lab
	}

	fmt.Println(printLab(manyLabs).Render())

	return nil
}

func ShowRack(cmd *cobra.Command, args []string) error {
	configFile := cmd.Flag("config").Value.String()
	config, e := Load(configFile)

	path := fmt.Sprintf("racks/%s", args[0])

	if e != nil {
		return e
	}

	vcName := args[0]
	rack := new(dc.Rack)

	manyRacks := map[string]*dc.Rack{}

	client, err := NewClient(config)
	if err != nil {
		return err
	}

	if len(args) == 0 {
		if _, e := client.Get("racks", nil, manyRacks); e != nil {
			return e
		}
	} else {
		if _, e := client.Get(path, nil, rack); e != nil {
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

	path := fmt.Sprintf("servers/%s", args[0])

	if e != nil {
		return e
	}

	srvName := args[0]
	srv := new(compute.Server)

	manySrv := map[string]*compute.Server{}

	client, err := NewClient(config)
	if err != nil {
		return err
	}

	if len(args) == 0 {
		if _, e := client.Get("servers", nil, manySrv); e != nil {
			return e
		}
	} else {
		if _, e := client.Get(path, nil, srv); e != nil {
			return e
		}

		manySrv[srvName] = srv
	}

	fmt.Println(printServer(manySrv).Render())

	return nil
}

func ShowESX(cmd *cobra.Command, args []string) error {
	configFile := cmd.Flag("config").Value.String()
	config, e := Load(configFile)

	path := fmt.Sprintf("esxs/%s", args[0])

	if e != nil {
		return e
	}

	esxName := args[0]
	esx := new(compute.ESX)

	manyESX := map[string]*compute.ESX{}

	client, err := NewClient(config)
	if err != nil {
		return err
	}

	if len(args) == 0 {
		if _, e := client.Get("esxs", nil, manyESX); e != nil {
			return e
		}
	} else {
		if _, e := client.Get(path, nil, esx); e != nil {
			return e
		}

		manyESX[esxName] = esx
	}

	fmt.Println(printESX(manyESX).Render())

	return nil
}

func ShowVC(cmd *cobra.Command, args []string) error {
	configFile := cmd.Flag("config").Value.String()
	config, e := Load(configFile)

	path := fmt.Sprintf("vcenters/%s", args[0])

	if e != nil {
		return e
	}

	vcName := args[0]
	vc := new(compute.VCenter)
	manyVC := map[string]*compute.VCenter{}

	client, err := NewClient(config)
	if err != nil {
		return err
	}

	if len(args) == 0 {
		if _, e := client.Get("vcenters", nil, manyVC); e != nil {
			return e
		}
	} else {
		if _, e := client.Get(path, nil, vc); e != nil {
			return e
		}

		manyVC[vcName] = vc
	}

	fmt.Println(printVC(manyVC).Render())

	return nil
}

func ShowVM(cmd *cobra.Command, args []string) error {
	configFile := cmd.Flag("config").Value.String()
	config, e := Load(configFile)

	path := fmt.Sprintf("vms/%s", args[0])

	if e != nil {
		return e
	}

	vcName := args[0]
	vm := new(compute.VM)
	manyVM := map[string]*compute.VM{}

	client, e := NewClient(config)

	if e != nil {
		return e
	}

	if len(args) == 0 {
		if _, e := client.Get("vms", nil, manyVM); e != nil {
			return e
		}
	} else {
		if _, e := client.Get(path, nil, vm); e != nil {
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
	data.AppendHeader(table.Row{"ID", "Name", "Type", "Labels"})

	for piece, lb := range labs {
		data.AppendRow(table.Row{
			piece,
			fmt.Sprintf(lb.NamedResource.ID),
			fmt.Sprintf(lb.NamedResource.Name),

			fmt.Sprintf(lb.NamedResource.Type),
			fmt.Sprintf("%s", lb.NamedResource.Labels),
		})
	}

	return data
}

func printDC(dcs map[string]*dc.Datacenter) table.Writer {
	data := table.NewWriter()
	data.AppendHeader(table.Row{"ID", "Name", "Type", "Address", "Labels"})

	for piece, dc := range dcs {
		data.AppendRow(table.Row{
			piece,
			fmt.Sprintf(dc.NamedResource.ID),
			fmt.Sprintf(dc.NamedResource.Name),

			fmt.Sprintf(dc.NamedResource.Type),
			fmt.Sprintf(dc.Address),
			fmt.Sprintf("%s", dc.NamedResource.Labels),
		})
	}

	return data
}

func printServer(servers map[string]*compute.Server) table.Writer {
	data := table.NewWriter()
	data.AppendHeader(table.Row{"ID", "Name", "Board IP", "Type", "Model", "Credentials", "Labels"})

	for piece, s := range servers {
		data.AppendRow(table.Row{
			piece,
			fmt.Sprintf(s.NamedResource.ID),
			fmt.Sprintf(s.NamedResource.Name),
			s.BoardIP.String(),

			fmt.Sprintf(s.NamedResource.Type),
			fmt.Sprintf(s.Model),

			fmt.Sprintf("%s", s.Credentials.Keys),
			fmt.Sprintf("%s", s.NamedResource.Labels),
		})
	}

	return data
}

func printESX(manyEsx map[string]*compute.ESX) table.Writer {
	data := table.NewWriter()
	data.AppendHeader(table.Row{"ID", "Name", "Server ID", "IP", "Type", "Credentials", "Labels"})

	for piece, esx := range manyEsx {
		data.AppendRow(table.Row{
			piece,
			fmt.Sprintf(esx.NamedResource.ID),
			fmt.Sprintf(esx.NamedResource.Name),

			fmt.Sprintf(esx.ServerID),
			esx.IP.String(),

			fmt.Sprintf(esx.NamedResource.Type),
			fmt.Sprintf("%s", esx.Credentials.Keys),
			fmt.Sprintf("%s", esx.NamedResource.Labels),
		})
	}

	return data
}

func printVC(manyVC map[string]*compute.VCenter) table.Writer {
	data := table.NewWriter()
	data.AppendHeader(table.Row{"ID", "Name", "IP", "Type", "Credentials", "Labels"})

	for piece, vc := range manyVC {
		data.AppendRow(table.Row{
			piece,
			fmt.Sprintf(vc.NamedResource.ID),

			fmt.Sprintf(vc.NamedResource.Name),
			vc.IP.String(),
			fmt.Sprintf(vc.NamedResource.Type),
			fmt.Sprintf("%s", vc.Credentials.Keys),
			fmt.Sprintf("%s", vc.NamedResource.Labels),
		})
	}

	return data
}

func printVM(manyVM map[string]*compute.VM) table.Writer {
	data := table.NewWriter()
	data.AppendHeader(table.Row{
		"ID", "Name", "IP", "Type", "Credentials",
		"ESXID", "VCenter ID", "Management IP", "Labels",
	})

	for piece, vm := range manyVM {
		data.AppendRow(table.Row{
			piece,
			fmt.Sprintf(vm.NamedResource.ID),

			fmt.Sprintf(vm.NamedResource.Name),
			fmt.Sprintf(vm.NamedResource.Type),
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
	data.AppendHeader(table.Row{"ID", "Name", "Status", "Type", "Row", "Labels"})

	for piece, rack := range racks {
		data.AppendRow(table.Row{
			piece,
			fmt.Sprintf(rack.NamedResource.ID),
			fmt.Sprintf(rack.NamedResource.Name),

			fmt.Sprintf(rack.NamedResource.Type),
			fmt.Sprintf(rack.Row),
			fmt.Sprintf("%s", rack.NamedResource.Labels),
		})
	}

	return data
}

func printUser(users map[string]*auth.User) table.Writer {
	data := table.NewWriter()
	data.AppendHeader(table.Row{"ID", "Name", "Status", "Type", "Password Hash", "Role", "Priviledges", "Labels"})

	for piece, user := range users {
		data.AppendRow(table.Row{
			piece,
			fmt.Sprintf(user.NamedResource.ID),
			fmt.Sprintf(user.NamedResource.Name),

			fmt.Sprintf(user.NamedResource.Type),
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
	data.AppendHeader(table.Row{"ID", "Status", "Type", "Range Start", "Range End", "Labels"})

	for piece, vlan := range vlans {
		data.AppendRow(table.Row{
			piece,
			fmt.Sprintf(vlan.ID),
			fmt.Sprintf(vlan.Status.UsedBy),
			fmt.Sprintf(vlan.Type),
			fmt.Sprintf("%010d", vlan.RangeStart),
			fmt.Sprintf("%010d", vlan.RangeEnd),
			fmt.Sprintf("%s", vlan.Labels),
		})
	}

	return data
}

func printIP(vlans map[string]*network.IPAddressPool) table.Writer {
	data := table.NewWriter()
	data.AppendHeader(table.Row{"ID", "Status", "Type", "Subnets", "Labels"})

	for piece, pool := range vlans {
		data.AppendRow(table.Row{
			piece,
			fmt.Sprintf(pool.ID),
			fmt.Sprintf(pool.Status.UsedBy),
			fmt.Sprintf(pool.Type),
			fmt.Sprintf("%s", pool.Subnets),
			fmt.Sprintf("%s", pool.Labels),
		})
	}

	return data
}
