package dc

import (
	"context"
	"errors"

	"github.com/rchamarthy/zebra"
)

var ErrAddressEmpty = errors.New("address is empty")

var ErrRowEmpty = errors.New("row is empty")

type Datacenter struct {
	zebra.NamedResource
	Address string `json:"address"`
}

func (dc *Datacenter) Validate(ctx context.Context) error {
	if dc.Address == "" {
		return ErrAddressEmpty
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
		return ErrRowEmpty
	}

	return r.NamedResource.Validate(ctx)
}
