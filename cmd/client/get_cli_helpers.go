package main

import (
	"path"

	"github.com/project-safari/zebra"
)

func GetType(typ string) *zebra.ResourceMap {
	factory := new(zebra.ResourceFactory)

	named := new(zebra.Type)
	named.Name = typ

	// (*factory).Add(*named)

	(*factory).New(typ)

	res := new(zebra.Resource)

	mapped := zebra.NewResourceMap(*factory)

	mapped.Add(*res, typ)

	return mapped
}

func GetPath(config *Config, p string, resMap interface{}) (int, error) {
	theRes := new(zebra.Resource)

	client, err := NewClient(config)
	if err != nil {
		return 0, err
	}

	return client.Get(path.Join("resources", config.User, p),
		theRes, resMap)
}
