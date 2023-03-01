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

	showCmd.AddCommand(&cobra.Command{
		Use:          "resources",
		Short:        "show all the resources",
		RunE:         showResources,
		Args:         cobra.MaximumNArgs(0),
		SilenceUsage: true,
	})

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

	leasecmd := &cobra.Command{
		Use:          "lease",
		Short:        "show datacenter lease information",
		RunE:         showLeases,
		Args:         cobra.MaximumNArgs(1),
		SilenceUsage: true,
	}
	leasecmd.Flags().BoolP("request", "r", true, "show the resource requests of the users lease(s)")
	showCmd.AddCommand(leasecmd)

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

// show server resource types - server, esx, vcenter, vm.
func showServers(cmd *cobra.Command, args []string) error {
	code, resMap, err := justGet(cmd, "resources", "compute.server")
	if err != nil {
		return err
	}

	if code != http.StatusOK {
		return ErrQuery
	}

	if l, ok := resMap.Resources["compute.server"]; ok {
		printServers(l.Resources)
	}

	return nil
}

func showESX(cmd *cobra.Command, args []string) error {
	code, resMap, err := justGet(cmd, "resources", "compute.esx")
	if err != nil {
		return err
	}

	if code != http.StatusOK {
		return ErrQuery
	}

	if l, ok := resMap.Resources["compute.esx"]; ok {
		printESX(l.Resources)
	}

	return nil
}

func showVCenters(cmd *cobra.Command, args []string) error {
	code, resMap, err := justGet(cmd, "resources", "compute.vcenter")
	if err != nil {
		return err
	}

	if code != http.StatusOK {
		return ErrQuery
	}

	if l, ok := resMap.Resources["compute.vcenter"]; ok {
		printVCenters(l.Resources)
	}

	return nil
}

func showVM(cmd *cobra.Command, args []string) error {
	code, resMap, err := justGet(cmd, "resources", "compute.vm")
	if err != nil {
		return err
	}

	if code != http.StatusOK {
		return ErrQuery
	}

	if l, ok := resMap.Resources["compute.vm"]; ok {
		printVM(l.Resources)
	}

	return nil
}

// show dc resource types - datacenter, lab, rack.
func showDatacenters(cmd *cobra.Command, args []string) error {
	code, resMap, err := justGet(cmd, "resources", "dc.datacenter")
	if err != nil {
		return err
	}

	if code != http.StatusOK {
		return ErrQuery
	}

	if l, ok := resMap.Resources["dc.datacenter"]; ok {
		printDatacenters(l.Resources)
	}

	return nil
}

func showLabs(cmd *cobra.Command, args []string) error {
	code, resMap, err := justGet(cmd, "resources", "dc.lab")
	if err != nil {
		return err
	}

	if code != http.StatusOK {
		return ErrQuery
	}

	if l, ok := resMap.Resources["dc.lab"]; ok {
		printLabs(l.Resources)
	}

	return nil
}

func showRacks(cmd *cobra.Command, args []string) error {
	code, resMap, err := justGet(cmd, "resources", "dc.rack")
	if err != nil {
		return err
	}

	if code != http.StatusOK {
		return ErrQuery
	}

	if l, ok := resMap.Resources["dc.rack"]; ok {
		printRacks(l.Resources)
	}

	return nil
}

// show network resource types: vlan, switch, IPAddressPool.
func showVlans(cmd *cobra.Command, args []string) error {
	code, resMap, err := justGet(cmd, "resources", "network.vlanPool")
	if err != nil {
		return err
	}

	if code != http.StatusOK {
		return ErrQuery
	}

	if l, ok := resMap.Resources["network.vlanPool"]; ok {
		printVlans(l.Resources)
	}

	return nil
}

func showSwitches(cmd *cobra.Command, args []string) error {
	code, resMap, err := justGet(cmd, "resources", "network.switch")
	if err != nil {
		return err
	}

	if code != http.StatusOK {
		return ErrQuery
	}

	if l, ok := resMap.Resources["network.switch"]; ok {
		printSwitches(l.Resources)
	}

	return nil
}

func showIPs(cmd *cobra.Command, args []string) error {
	code, resMap, err := justGet(cmd, "resources", "network.ipAddressPool")
	if err != nil {
		return err
	}

	if code != http.StatusOK {
		return ErrQuery
	}

	if l, ok := resMap.Resources["network.ipAddressPool"]; ok {
		printIPs(l.Resources)
	}

	return nil
}

func showLeases(cmd *cobra.Command, args []string) error {
	code, resMap, err := justGet(cmd, "resources", "system.lease")
	if err != nil {
		return err
	}

	cfgFile := cmd.Flag("config").Value.String()

	cfg, err := Load(cfgFile)
	if err != nil {
		return err
	}

	if code != http.StatusOK {
		return ErrQuery
	}

	reqbool := cmd.Flag("request").Value.String()
	if l, ok := resMap.Resources["system.lease"]; ok {
		if reqbool != cmd.Flag("request").DefValue {
			printResourceReq(getUserLeases(cfg.Email, l.Resources))
		} else {
			printLeases(l.Resources)
		}
	}

	return nil
}

func showUsers(cmd *cobra.Command, args []string) error {
	code, resMap, err := justGet(cmd, "resources", "system.user")
	if err != nil {
		return err
	}

	if code != http.StatusOK {
		return ErrQuery
	}

	if l, ok := resMap.Resources["system.user"]; ok {
		printUsers(l.Resources)
	}

	return nil
}

func showRegistrations(cmd *cobra.Command, args []string) error {
	code, resMap, err := justGet(cmd, "resources", "system.user")
	if err != nil {
		return err
	}

	if code != http.StatusOK {
		return ErrQuery
	}

	if l, ok := resMap.Resources["system.user"]; ok {
		printUsers(l.Resources)
	}

	return nil
}

func showPublicKey(cmd *cobra.Command, args []string) error {
	cfgFile := cmd.Flag("config").Value.String()

	cfg, e := Load(cfgFile)
	if e != nil {
		return e
	}

	fmt.Println(cfg.Key.Public())

	return nil
}

func state(r zebra.Resource) string {
	return r.GetStatus().State.String()
}

func usedBy(r zebra.Resource) string {
	s := r.GetStatus().UsedBy
	if s == "" {
		s = "--"
	}

	return s
}

func getUserLeases(email string, res []zebra.Resource) []zebra.Resource {
	var ulease []zebra.Resource

	for _, r := range res {
		if r.GetMeta().Owner == email {
			ulease = append(ulease, r)
		}
	}

	return ulease
}

func printResources(resources *zebra.ResourceMap) {
	tw := table.NewWriter()
	tw.AppendHeader(table.Row{"Name", "ID", "Type", "Group(s)", "Status", "Lease State"})

	for t, l := range resources.Resources {
		for _, resource := range l.Resources {
			tw.AppendRow(table.Row{
				resource.GetMeta().Name,
				resource.GetMeta().ID,
				t,
				resource.GetMeta().Labels,
				state(resource),
				resource.GetStatus().LeaseStatus,
			})
		}
	}

	fmt.Println(tw.Render())
}

// print server resource types - servers, esx, vcenters, vms.
func printServers(servers []zebra.Resource) {
	tw := table.NewWriter()
	tw.AppendHeader(table.Row{
		"Name", "ID", "Board IP", "Group(s)",
		"Model", "Serialnumber", "User", "Status", "Lease State",
	})

	for _, s := range servers {
		if server, ok := s.(*compute.Server); ok {
			tw.AppendRow(table.Row{
				server.GetMeta().Name,
				server.GetMeta().ID,
				server.BoardIP,
				server.GetMeta().Labels,
				server.Model,
				server.SerialNumber,
				usedBy(server),
				state(server),
				server.GetStatus().LeaseStatus,
			})
		}
	}

	fmt.Println(tw.Render())
}

func printESX(manyESX []zebra.Resource) {
	data := table.NewWriter()
	data.AppendHeader(table.Row{
		"Name", "ID", "Server ID", "Group(s)", "IP",
		"Credentials", "User", "Status", "Lease State",
	})

	for _, e := range manyESX {
		if esx, ok := e.(*compute.ESX); ok {
			data.AppendRow(table.Row{
				esx.GetMeta().Name,
				esx.GetMeta().ID,
				esx.ServerID,
				esx.GetMeta().Labels,
				esx.IP.String(),
				esx.Credentials.Keys,
				usedBy(esx),
				state(esx),
				esx.GetStatus().LeaseStatus,
			})
		}
	}

	fmt.Println(data.Render())
}

func printVCenters(manyVC []zebra.Resource) {
	data := table.NewWriter()
	data.AppendHeader(table.Row{"Name", "ID", "IP", "Group(s)", "Credentials", "User", "Status", "Lease State"})

	for _, vc := range manyVC {
		if vcenter, ok := vc.(*compute.VCenter); ok {
			data.AppendRow(table.Row{
				vcenter.GetMeta().Name,
				vcenter.GetMeta().ID,
				vcenter.IP.String(),
				vcenter.GetMeta().Labels,
				vcenter.Credentials.Keys,
				usedBy(vcenter),
				state(vcenter),
			})
		}
	}

	fmt.Println(data.Render())
}

func printVM(manyVM []zebra.Resource) {
	data := table.NewWriter()
	data.AppendHeader(table.Row{
		"Name", "ID", "IP", "Group(s)", "Credentials",
		"ESXID", "VCID", "User", "Status", "Lease State",
	})

	for _, vm := range manyVM {
		if machine, ok := vm.(*compute.VM); ok {
			data.AppendRow(table.Row{
				machine.GetMeta().Name,
				machine.GetMeta().ID,
				machine.ManagementIP.String(),
				machine.GetMeta().Labels,
				machine.Credentials.Keys,
				machine.ESXID,
				machine.VCenterID,
				usedBy(machine),
				state(machine),
				machine.GetStatus().LeaseStatus,
			})
		}
	}

	fmt.Println(data.Render())
}

// print dc resource types - datacenter, lab, rack.
func printDatacenters(dcs []zebra.Resource) {
	data := table.NewWriter()
	data.AppendHeader(table.Row{"Name", "ID", "Address", "Group(s)", "User", "Status", "Lease State"})

	for _, d := range dcs {
		if dc, ok := d.(*dc.Datacenter); ok {
			data.AppendRow(table.Row{
				dc.GetMeta().Name,
				dc.GetMeta().ID,
				dc.Address,
				dc.GetMeta().Labels,
				usedBy(dc),
				state(dc),
				dc.GetStatus().LeaseStatus,
			})
		}
	}

	fmt.Println(data.Render())
}

func printLabs(labs []zebra.Resource) {
	data := table.NewWriter()
	data.AppendHeader(table.Row{"Name", "ID", "User", "Group(s)", "Type", "Status", "Lease Status"})

	for _, lb := range labs {
		if lab, ok := lb.(*dc.Lab); ok {
			data.AppendRow(table.Row{
				lab.GetMeta().Name,
				lab.Meta.ID,
				usedBy(lab),
				lab.GetMeta().Labels,
				lab.GetMeta().Type.Name,
				state(lab),
				lab.GetStatus().LeaseStatus,
			})
		}
	}

	fmt.Println(data.Render())
}

func printRacks(racks []zebra.Resource) {
	data := table.NewWriter()
	data.AppendHeader(table.Row{"Name", "ID", "Row", "Group(s)", "User", "Status", "Lease State"})

	for _, r := range racks {
		if rack, ok := r.(*dc.Rack); ok {
			data.AppendRow(table.Row{
				rack.GetMeta().Name,
				rack.GetMeta().ID,
				rack.Row,
				rack.GetMeta().Labels,
				usedBy(rack),
				state(rack),
				rack.GetStatus().LeaseStatus,
			})
		}
	}

	fmt.Println(data.Render())
}

// print network resource types: vlan, switch, IPAddressPool.
func printVlans(vlans []zebra.Resource) {
	data := table.NewWriter()
	data.AppendHeader(table.Row{"VLanPool", "ID", "User", "Group(s)", "Status", "Lease State"})

	for _, v := range vlans {
		vlan, ok := v.(*network.VLANPool)

		if ok {
			data.AppendRow(table.Row{
				vlan.String(),
				vlan.GetMeta().ID,
				usedBy(vlan),
				vlan.GetMeta().Labels,
				state(vlan),
				vlan.GetStatus().LeaseStatus,
			})
		}
	}

	fmt.Println(data.Render())
}

func printSwitches(switches []zebra.Resource) {
	data := table.NewWriter()
	data.AppendHeader(table.Row{
		"Name", "ID", "Management IP", "Group(s)", "Credentials",
		"Serial Number", "Model", "Ports", "User", "Status", "Lease State",
	})

	for _, s := range switches {
		if sw, ok := s.(*network.Switch); ok {
			data.AppendRow(table.Row{
				sw.GetMeta().Name,
				sw.GetMeta().ID,
				sw.ManagementIP.String(),
				sw.GetMeta().Labels,
				sw.Credentials.Keys,
				sw.SerialNumber,
				sw.Model,
				sw.NumPorts,
				usedBy(sw),
				state(sw),
				sw.GetStatus().LeaseStatus,
			})
		}
	}

	fmt.Println(data.Render())
}

func printIPs(ips []zebra.Resource) {
	data := table.NewWriter()
	data.AppendHeader(table.Row{"Subnets", "ID", "User", "Group(s)", "Status", "Lease State"})

	for _, addr := range ips {
		if pool, ok := addr.(*network.IPAddressPool); ok {
			data.AppendRow(table.Row{
				pool.Subnets,
				pool.GetMeta().ID,
				usedBy(pool),
				pool.GetMeta().Labels,
				state(pool),
				pool.GetStatus().LeaseStatus,
			})
		}
	}

	fmt.Println(data.Render())
}

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

func printResourceReq(leases []zebra.Resource) {
	tw := table.NewWriter()
	tw.AppendHeader(table.Row{
		"Name", "Type", "Group", "Count", "Resources Assigned",
	})

	// print the table here.
	for _, s := range leases {
		if l, ok := s.(*lease.Lease); ok {
			for _, req := range l.Request {
				tw.AppendRow(table.Row{
					req.Name,
					req.Type,
					req.Group,
					req.Count,
					req.Resources,
				})
			}
		}
	}

	fmt.Println(tw.Render())
}
