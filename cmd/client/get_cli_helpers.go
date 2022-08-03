package main

import (
	"path"

	"github.com/project-safari/zebra/auth"
	"github.com/project-safari/zebra/compute"
	"github.com/project-safari/zebra/dc"
	"github.com/project-safari/zebra/network"
)

func GetNetPaths(config *Config, p, typ string, data interface{}) (int, error) {
	client, err := NewClient(config)
	if err != nil {
		return 0, err
	}

	switch typ {
	case "users", "registrations":
		return client.Get(path.Join("resources", config.User),
			new(auth.User), map[string]*auth.User{})

	case "vlans":
		return client.Get(path.Join("resources", config.User),
			new(network.VLANPool), map[string]*network.VLANPool{})

	case "switches":
		return client.Get(path.Join("resources", config.User),
			new(network.Switch), map[string]*network.Switch{})

	case "ips":
		return client.Get(path.Join("resources", config.User),
			new(network.IPAddressPool), map[string]*network.IPAddressPool{})
	}

	return 0, nil
}

func GetDCPaths(config *Config, p, typ string, data interface{}) (int, error) {
	client, err := NewClient(config)
	if err != nil {
		return 0, err
	}

	switch typ {
	case "datacenters":
		return client.Get(path.Join("resources", config.User),
			new(dc.Datacenter), map[string]*dc.Datacenter{})

	case "labs":
		return client.Get(path.Join("resources", config.User),
			new(dc.Lab), map[string]*dc.Lab{})

	case "racks":
		return client.Get(path.Join("resources", config.User),
			new(dc.Rack), map[string]*dc.Rack{})
	}

	return 0, nil
}

func GetComputePaths(config *Config, p, typ string, data interface{}) (int, error) {
	client, err := NewClient(config)
	if err != nil {
		return 0, err
	}

	switch typ {
	case "servers":
		return client.Get(path.Join("resources", config.User),
			new(compute.Server), map[string]*compute.Server{})

	case "esxes":
		return client.Get(path.Join("resources", config.User),
			new(compute.ESX), map[string]*compute.ESX{})

	case "vcenter":
		return client.Get(path.Join("resources", config.User),
			new(compute.VCenter), map[string]*compute.VCenter{})

	case "vms":
		return client.Get(path.Join("resources", config.User),
			new(compute.VM), map[string]*compute.VM{})
	}

	return 0, nil
}
