package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/auth"
	"github.com/project-safari/zebra/compute"
	"github.com/project-safari/zebra/lease"
	"github.com/project-safari/zebra/store"
	"github.com/rodaine/table"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var ErrShow = errors.New("error with show command")

func NewShow() *cobra.Command {
	showCmd := &cobra.Command{
		Use:          "show",
		Short:        "show resources (filter by type)",
		RunE:         showResources,
		SilenceUsage: true,
		Args:         cobra.MaximumNArgs(1),
	}

	return showCmd
}

func showResources(cmd *cobra.Command, args []string) error {
	cfg, req, err := makeShowReq(cmd, args)
	if err != nil {
		return err
	}

	client, err := NewClient(cfg)
	if err != nil {
		return err
	}

	res := zebra.NewResourceMap(store.DefaultFactory())

	resCode, err := client.Get("api/v1/resources", req, res)
	if resCode != http.StatusOK {
		return ErrShow
	}

	if err != nil {
		return err
	}

	printResMap(res)

	return err
}

func makeShowReq(cmd *cobra.Command, args []string) (*Config, interface{}, error) {
	cfgFile := cmd.Flag("config").Value.String()

	cfg, err := Load(cfgFile)
	if err != nil {
		return nil, nil, err
	}

	types := args

	// If all, make an empty request to show all resources.
	if len(types) == 0 || strings.ToLower(types[0]) == "all" {
		types = []string{}
	}

	// Right now, manually change caps
	for i, t := range types {
		switch strings.ToLower(t) {
		case "ipaddresspool":
			types[i] = "IPAddressPool"
		case "vlanpool":
			types[i] = "VLANPool"
			types[i] = "ESX"
		case "vcenter":
			types[i] = "VCenter"
		case "vm":
			types[i] = "VM"
		default:
			types[i] = cases.Title(language.AmericanEnglish).String(t)
		}
	}

	qr := struct {
		IDs        []string      `json:"ids,omitempty"`
		Types      []string      `json:"types,omitempty"`
		Labels     []zebra.Query `json:"labels,omitempty"`
		Properties []zebra.Query `json:"properties,omitempty"`
	}{
		Types: types,
	}

	return cfg, qr, nil
}

func printResMap(resMap *zebra.ResourceMap) {
	first := true
	for t, l := range resMap.Resources {
		if first {
			fmt.Println()
		}

		printResList(t, l)
		fmt.Println()
	}
}

func printResList(t string, l *zebra.ResourceList) {
	switch t {
	case "User":
		fmt.Printf("-- USERS --\n\n")

		table.DefaultHeaderFormatter = func(format string, vals ...interface{}) string {
			return strings.ToUpper(fmt.Sprintf(format, vals...))
		}
		tbl := table.New("Name", "Email", "Role")

		for _, res := range l.Resources {
			u, _ := res.(*auth.User)
			tbl.AddRow(u.Name, u.Email, u.Role.Name)
		}

		tbl.Print()
	case "Server":
		fmt.Printf("-- SERVERS --\n\n")

		table.DefaultHeaderFormatter = func(format string, vals ...interface{}) string {
			return strings.ToUpper(fmt.Sprintf(format, vals...))
		}
		tbl := table.New("Name", "Serial Number", "Board IP", "Model", "Group")

		for _, res := range l.Resources {
			s, _ := res.(*compute.Server)
			tbl.AddRow(s.Name, s.SerialNumber, s.BoardIP, s.Model, s.GetLabels()["system.group"])
		}

		tbl.Print()
	case "Lease":
		fmt.Printf("-- LEASES --\n\n")

		table.DefaultHeaderFormatter = func(format string, vals ...interface{}) string {
			return strings.ToUpper(fmt.Sprintf(format, vals...))
		}
		tbl := table.New("Owner", "Type", "Count", "Duration", "Status")

		for _, res := range l.Resources {
			lease, _ := res.(*lease.Lease)
			tbl.AddRow(lease.Status.UsedBy, lease.Request[0].Type, lease.Request[0].Count, lease.Duration, lease.Status.State)
		}

		tbl.Print()
	}
}
