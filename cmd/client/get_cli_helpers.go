package main

import (
	"strings"

	"github.com/project-safari/zebra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func GetType(args []string) interface{} { //nolint:cyclop
	typ := args

	if len(typ) == 0 || strings.ToLower(typ[0]) == "all" {
		typ = []string{}
	}

	for i, t := range typ {
		switch strings.ToLower(t) {
		case "ipaddresspool", "ip", "ips", "ippool":
			typ[i] = "IPAddressPool"
		case "server", "servers":
			typ[i] = "Server"
		case "switch", "switches":
			typ[i] = "Switch"
		case "vlanpool", "vlan", "vlans":
			typ[i] = "VLANPool"
		case "esx", "esxes":
			typ[i] = "ESX"
		case "vcenter", "vcenters":
			typ[i] = "VCenter"
		case "vm", "vms", "VirtualMachine":
			typ[i] = "VM"
		case "dc", "datacenter", "datacenters":
			typ[i] = "Datacenter"
		case "lab", "labs":
			typ[i] = "Lab"
		case "rack", "racks":
			typ[i] = "Rack"
		default:
			typ[i] = cases.Title(language.AmericanEnglish).String(t)
		}
	}

	qr := struct {
		IDs        []string      `json:"ids,omitempty"`
		Types      []string      `json:"types,omitempty"`
		Labels     []zebra.Query `json:"labels,omitempty"`
		Properties []zebra.Query `json:"properties,omitempty"`
	}{
		Types: typ,
	}

	return qr
}

func GetPath(config *Config, resMap interface{}) (int, error) {
	p := "api/v1/resources"

	theRes := new(interface{})

	client, err := NewClient(config)
	if err != nil {
		return 0, err
	}

	// clinet with path and resources
	return client.Get(p, theRes, resMap)
}
