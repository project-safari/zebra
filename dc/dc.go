package dc

import (
	"context"
	"fmt"

	"github.com/rchamarthy/zebra"
)

type Datacenter struct {
	zebra.NamedResource
	Address string `json:"address"`
}

func (dc *Datacenter) Validate(ctx context.Context) error {
	if dc.Address == "" {
		return fmt.Errorf("datacenter must have an address")
	}

	return dc.NamedResource.Validate(ctx)
}

type Lab struct {
	zebra.NamedResource
}

type Rack struct {
	zebra.NamedResource
	Row string `json:"row"`
}

func (r *Rack) Validate(ctx context.Context) error {
	if r.Row == "" {
		return fmt.Errorf("row cannot be empty")
	}

	return r.NamedResource.Validate(ctx)
}
