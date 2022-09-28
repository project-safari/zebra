package main

import (
	"errors"
	"fmt"
	"net/http"
	"path"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/model"
	"github.com/project-safari/zebra/model/compute"
	"github.com/project-safari/zebra/model/dc"
	"github.com/project-safari/zebra/model/lease"
	"github.com/project-safari/zebra/model/network"
	"github.com/project-safari/zebra/model/user"
	"github.com/spf13/cobra"
)

var ErrQuery = errors.New("server query failed")

type QueryRequest struct {
	IDs        []string      `json:"ids,omitempty"`
	Types      []string      `json:"types,omitempty"`
	Labels     []zebra.Query `json:"labels,omitempty"`
	Properties []zebra.Query `json:"properties,omitempty"`
}

func NewShow() *cobra.Command { //nolint:funlen
	showCmd := &cobra.Command{
		Use:   "show",
		Short: "show resources",
	}

	// server resource types - server, esx, vcenter, vm.
	showCmd.AddCommand(&cobra.Command{
		Use:          "server",
		Short:        "show the specified servers",
		RunE:         showServers,
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
	})

	showCmd.AddCommand(&cobra.Command{
		Use:          "esx",
		Short:        "show the specified esxes",
		RunE:         showESX,
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
	})

	showCmd.AddCommand(&cobra.Command{
		Use:          "vcenter",
		Short:        "show the specified vcenters",
		RunE:         showVCenters,
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
	})

	showCmd.AddCommand(&cobra.Command{
		Use:          "vm",
		Short:        "show the specified vms",
		RunE:         showVM,
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
	})

	// dc resource types - datacenter, lab, rack.
	showCmd.AddCommand(&cobra.Command{
		Use:          "datacenter",
		Short:        "show the specified datacenters",
		RunE:         showDatacenters,
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
	})

	showCmd.AddCommand(&cobra.Command{
		Use:          "lab",
		Short:        "show the specified labs",
		RunE:         showLabs,
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
	})

	showCmd.AddCommand(&cobra.Command{
		Use:          "rack",
		Short:        "show the specified racks",
		RunE:         showRacks,
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
	})

	// network resource types: vlan, switch, IPAddressPool.
	showCmd.AddCommand(&cobra.Command{
		Use:          "vlan",
		Short:        "show the specified vlans",
		RunE:         showVlans,
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
	})

	showCmd.AddCommand(&cobra.Command{
		Use:          "switch",
		Short:        "show the specified switches",
		RunE:         showSwitches,
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
	})

	showCmd.AddCommand(&cobra.Command{
		Use:          "ip",
		Short:        "show the specified IPAddressPools",
		RunE:         showIPs,
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
	})

	showCmd.AddCommand(&cobra.Command{
		Use:          "lease",
		Short:        "show datacenter lease information",
		RunE:         showLeases,
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
	})

	showCmd.AddCommand(&cobra.Command{
		Use:          "user",
		Short:        "show users",
		RunE:         showUsers,
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
	})

	showCmd.AddCommand(&cobra.Command{
		Use:          "registration",
		Short:        "show registrations",
		RunE:         showRegistrations,
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
	})

	showCmd.AddCommand(&cobra.Command{
		Use:          "public-key",
		Short:        "show the public key of the user",
		RunE:         showPublicKey,
		Args:         cobra.ExactArgs(0),
		SilenceUsage: true,
	})

	return showCmd
}

func justGet(cmd *cobra.Command, p string, resTypes ...string) (int, *zebra.ResourceMap, error) {
	cfgFile := cmd.Flag("config").Value.String()

	cfg, e := Load(cfgFile)
	if e != nil {
		return 0, nil, e
	}

	c, e := NewClient(cfg)
	if e != nil {
		return 0, nil, e
	}

	in := &QueryRequest{Types: resTypes}
	resMap := zebra.NewResourceMap(model.Factory())
	status, err := c.Get(path.Join("api", "v1", p), in, resMap)

	return status, resMap, err
}

func showResources(cmd *cobra.Command, args []string) error {
	code, resMap, err := justGet(cmd, "resources")
	if err != nil {
		return err
	}

	if code != http.StatusOK {
		return ErrQuery
	}

	printResources(resMap)

	return nil
}

// Show server resource types - server, esx, vcenter, vm.
// Function for servers.
func showServers(cmd *cobra.Command, args []string) error {
	code, resMap, err := justGet(cmd, "resources", "Server")
	if err != nil {
		return err
	}

	if code != http.StatusOK {
		return ErrQuery
	}

	if l, ok := resMap.Resources["Server"]; ok {
		printServers(l.Resources)
	}

	return nil
}

// Show server resource types - server, esx, vcenter, vm.
// Function for esx servers.
func showESX(cmd *cobra.Command, args []string) error {
	code, resMap, err := justGet(cmd, "resources", "ESX")
	if err != nil {
		return err
	}

	if code != http.StatusOK {
		return ErrQuery
	}

	if l, ok := resMap.Resources["ESX"]; ok {
		printESX(l.Resources)
	}

	return nil
}

// Show server resource types - server, esx, vcenter, vm.
// Function for VCenters.
func showVCenters(cmd *cobra.Command, args []string) error {
	code, resMap, err := justGet(cmd, "resources", "VCenter")
	if err != nil {
		return err
	}

	if code != http.StatusOK {
		return ErrQuery
	}

	if l, ok := resMap.Resources["VCenter"]; ok {
		printVCenters(l.Resources)
	}

	return nil
}

// Show server resource types - server, esx, vcenter, vm.
// Function for virtual machines.
func showVM(cmd *cobra.Command, args []string) error {
	code, resMap, err := justGet(cmd, "resources", "VM")
	if err != nil {
		return err
	}

	if code != http.StatusOK {
		return ErrQuery
	}

	if l, ok := resMap.Resources["VM"]; ok {
		printVM(l.Resources)
	}

	return nil
}

// Show dc resource types - datacenter, lab, rack.
// Function for datacenters.
func showDatacenters(cmd *cobra.Command, args []string) error {
	code, resMap, err := justGet(cmd, "resources", "Datacenter")
	if err != nil {
		return err
	}

	if code != http.StatusOK {
		return ErrQuery
	}

	if l, ok := resMap.Resources["Datacenter"]; ok {
		printDatacenters(l.Resources)
	}

	return nil
}

// Show dc resource types - datacenter, lab, rack.
// Function for labs.
func showLabs(cmd *cobra.Command, args []string) error {
	code, resMap, err := justGet(cmd, "resources", "Lab")
	if err != nil {
		return err
	}

	if code != http.StatusOK {
		return ErrQuery
	}

	if l, ok := resMap.Resources["Lab"]; ok {
		printLabs(l.Resources)
	}

	return nil
}

// Show dc resource types - datacenter, lab, rack.
// Function for racks.
func showRacks(cmd *cobra.Command, args []string) error {
	code, resMap, err := justGet(cmd, "resources", "Rack")
	if err != nil {
		return err
	}

	if code != http.StatusOK {
		return ErrQuery
	}

	if l, ok := resMap.Resources["Rack"]; ok {
		printRacks(l.Resources)
	}

	return nil
}

// Show network resource types: vlan, switch, IPAddressPool.
// Function for vlan pools.
func showVlans(cmd *cobra.Command, args []string) error {
	code, resMap, err := justGet(cmd, "resources", "VLANPool")
	if err != nil {
		return err
	}

	if code != http.StatusOK {
		return ErrQuery
	}

	if l, ok := resMap.Resources["VLANPool"]; ok {
		printVlans(l.Resources)
	}

	return nil
}

// Show network resource types: vlan, switch, IPAddressPool.
// Function for switches.
func showSwitches(cmd *cobra.Command, args []string) error {
	code, resMap, err := justGet(cmd, "resources", "Switch")
	if err != nil {
		return err
	}

	if code != http.StatusOK {
		return ErrQuery
	}

	if l, ok := resMap.Resources["Switch"]; ok {
		printSwitches(l.Resources)
	}

	return nil
}

// Show network resource types: vlan, switch, IPAddressPool.
// Function for IP address pools.
func showIPs(cmd *cobra.Command, args []string) error {
	code, resMap, err := justGet(cmd, "resources", "IPAddressPool")
	if err != nil {
		return err
	}

	if code != http.StatusOK {
		return ErrQuery
	}

	if l, ok := resMap.Resources["IPAddressPool"]; ok {
		printIPs(l.Resources)
	}

	return nil
}

// Function to show leases.
func showLeases(cmd *cobra.Command, args []string) error {
	code, resMap, err := justGet(cmd, "resources", "Lease")
	if err != nil {
		return err
	}

	if code != http.StatusOK {
		return ErrQuery
	}

	if l, ok := resMap.Resources["Lease"]; ok {
		printLeases(l.Resources)
	}

	return nil
}

// Function to show users.
func showUsers(cmd *cobra.Command, args []string) error {
	code, resMap, err := justGet(cmd, "resources", "User")
	if err != nil {
		return err
	}

	if code != http.StatusOK {
		return ErrQuery
	}

	if l, ok := resMap.Resources["User"]; ok {
		printUsers(l.Resources)
	}

	return nil
}

// Function to show registrations.
func showRegistrations(cmd *cobra.Command, args []string) error {
	code, resMap, err := justGet(cmd, "resources", "Registration")
	if err != nil {
		return err
	}

	if code != http.StatusOK {
		return ErrQuery
	}

	if l, ok := resMap.Resources["Registration"]; ok {
		printUsers(l.Resources)
	}

	return nil
}

// Function to show public keys.
func showPublicKey(cmd *cobra.Command, args []string) error {
	cfgFile := cmd.Flag("config").Value.String()

	cfg, e := Load(cfgFile)
	if e != nil {
		return e
	}

	fmt.Println(cfg.Key.Public())

	return nil
}

// Function to show the state of a resource.
func state(r zebra.Resource) string {
	return r.GetStatus().State.String()
}

// Function to show who uses a resource.
func usedBy(r zebra.Resource) string {
	s := r.GetStatus().UsedBy
	if s == "" {
		s = "--"
	}

	return s
}

func printResources(resources *zebra.ResourceMap) {
	tw := table.NewWriter()
	tw.AppendHeader(table.Row{"Name", "Type", "Status"})

	for t, l := range resources.Resources {
		for _, resource := range l.Resources {
			tw.AppendRow(table.Row{
				resource.GetMeta().Name,
				t,
				state(resource),
			})
		}
	}

	fmt.Println(tw.Render())
}

// Print server resource types - servers, esx, vcenters, vms.
// Function for servers.
func printServers(servers []zebra.Resource) {
	tw := table.NewWriter()
	tw.AppendHeader(table.Row{"Name", "Board IP", "Model", "Serialnumber", "User", "Status"})

	for _, s := range servers {
		if server, ok := s.(*compute.Server); ok {
			tw.AppendRow(table.Row{
				server.GetMeta().Name,
				server.BoardIP,
				server.Model,
				server.SerialNumber,
				usedBy(server),
				state(server),
			})
		}
	}

	fmt.Println(tw.Render())
}

// Print server resource types - servers, esx, vcenters, vms.
// Function for esx.
func printESX(manyESX []zebra.Resource) {
	data := table.NewWriter()
	data.AppendHeader(table.Row{"Name", "Server ID", "IP", "Credentials", "User", "Status"})

	for _, e := range manyESX {
		if esx, ok := e.(*compute.ESX); ok {
			data.AppendRow(table.Row{
				esx.GetMeta().Name,
				esx.ServerID,
				esx.IP.String(),
				esx.Credentials.Keys,
				usedBy(esx),
				state(esx),
			})
		}
	}

	fmt.Println(data.Render())
}

// Print server resource types - servers, esx, vcenters, vms.
// Function for vcenters.
func printVCenters(manyVC []zebra.Resource) {
	data := table.NewWriter()
	data.AppendHeader(table.Row{"Name", "IP", "Credentials", "User", "Status"})

	for _, vc := range manyVC {
		if vcenter, ok := vc.(*compute.VCenter); ok {
			data.AppendRow(table.Row{
				vcenter.GetMeta().Name,
				vcenter.IP.String(),
				vcenter.Credentials.Keys,
				usedBy(vcenter),
				state(vcenter),
			})
		}
	}

	fmt.Println(data.Render())
}

// Print server resource types - servers, esx, vcenters, vms.
// Function for virtual machines.
func printVM(manyVM []zebra.Resource) {
	data := table.NewWriter()
	data.AppendHeader(table.Row{"Name", "IP", "Credentials", "ESXID", "VCID", "User", "Status"})

	for _, vm := range manyVM {
		if machine, ok := vm.(*compute.VM); ok {
			data.AppendRow(table.Row{
				machine.GetMeta().Name,
				machine.ManagementIP.String(),
				machine.Credentials.Keys,
				machine.ESXID,
				machine.VCenterID,
				usedBy(machine),
				state(machine),
			})
		}
	}

	fmt.Println(data.Render())
}

// Print dc resource types - datacenter, lab, rack.
// Function for datacenters.
func printDatacenters(dcs []zebra.Resource) {
	data := table.NewWriter()
	data.AppendHeader(table.Row{"Name", "Address", "User", "Status"})

	for _, d := range dcs {
		if dc, ok := d.(*dc.Datacenter); ok {
			data.AppendRow(table.Row{
				dc.GetMeta().Name,
				dc.Address,
				usedBy(dc),
				state(dc),
			})
		}
	}

	fmt.Println(data.Render())
}

// Print dc resource types - datacenter, lab, rack.
// Function for labs.
func printLabs(labs []zebra.Resource) {
	data := table.NewWriter()
	data.AppendHeader(table.Row{"Name", "User", "Status"})

	for _, lb := range labs {
		if lab, ok := lb.(*dc.Lab); ok {
			data.AppendRow(table.Row{
				lab.GetMeta().Name,
				usedBy(lab),
				state(lab),
			})
		}
	}

	fmt.Println(data.Render())
}

// Print dc resource types - datacenter, lab, rack.
// Function for racks.
func printRacks(racks []zebra.Resource) {
	data := table.NewWriter()
	data.AppendHeader(table.Row{"Name", "Row", "User", "Status"})

	for _, r := range racks {
		if rack, ok := r.(*dc.Rack); ok {
			data.AppendRow(table.Row{
				rack.GetMeta().Name,
				rack.Row,
				usedBy(rack),
				state(rack),
			})
		}
	}

	fmt.Println(data.Render())
}

// Print network resource types: vlan, switch, IPAddressPool.
// Function for vlan pools.
func printVlans(vlans []zebra.Resource) {
	data := table.NewWriter()
	data.AppendHeader(table.Row{"VLanPool", "User", "Status"})

	for _, v := range vlans {
		vlan, ok := v.(*network.VLANPool)

		if ok {
			data.AppendRow(table.Row{
				vlan.String(),
				usedBy(vlan),
				state(vlan),
			})
		}
	}

	fmt.Println(data.Render())
}

// Print network resource types: vlan, switch, IPAddressPool.
// Function for switches.
func printSwitches(switches []zebra.Resource) {
	data := table.NewWriter()
	data.AppendHeader(table.Row{
		"Name", "Management IP", "Credentials",
		"Serial Number", "Model", "Ports", "User", "Status",
	})

	for _, s := range switches {
		if sw, ok := s.(*network.Switch); ok {
			data.AppendRow(table.Row{
				sw.GetMeta().Name,
				sw.ManagementIP.String(),
				sw.Credentials.Keys,
				sw.SerialNumber,
				sw.Model,
				sw.NumPorts,
				usedBy(sw),
				state(sw),
			})
		}
	}

	fmt.Println(data.Render())
}

// Print network resource types: vlan, switch, IPAddressPool.
// Function for IP address pools.
func printIPs(ips []zebra.Resource) {
	data := table.NewWriter()
	data.AppendHeader(table.Row{"Subnets", "User", "Status"})

	for _, addr := range ips {
		if pool, ok := addr.(*network.IPAddressPool); ok {
			data.AppendRow(table.Row{
				pool.Subnets,
				usedBy(pool),
				state(pool),
			})
		}
	}

	fmt.Println(data.Render())
}

// Function to print leases.
func printLeases(leases []zebra.Resource) {
	tw := table.NewWriter()
	tw.AppendHeader(table.Row{
		"Owner", "Requested Time", "Start Time",
		"Time Left", "Active",
	})

	// print the table here.
	for _, s := range leases {
		if l, ok := s.(*lease.Lease); ok {
			status := l.GetStatus()
			if status.State == zebra.Active {
				tw.AppendRow(table.Row{
					l.Status.UsedBy,
					l.Duration,
					l.ActivationTime,
					time.Until(l.ActivationTime.Add(l.Duration)),
					state(l),
				})
			} else {
				tw.AppendRow(table.Row{
					l.Status.UsedBy,
					l.Duration,
					"--",
					"--",
					state(l),
				})
			}
		}
	}

	fmt.Println(tw.Render())
}

// Function to print the users.
func printUsers(users []zebra.Resource) {
	data := table.NewWriter()
	data.AppendHeader(table.Row{"Name", "Role", "Privileges", "Status"})

	for _, u := range users {
		user, ok := u.(*user.User)

		if ok {
			data.AppendRow(table.Row{
				user.GetMeta().Name,
				user.Role.Name,
				user.Role.Privileges,
				state(user),
			})
		}
	}

	fmt.Println(data.Render())
}
