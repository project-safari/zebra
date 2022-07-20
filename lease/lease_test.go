package lease_test

import (
	"testing"
	"time"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/auth"
	"github.com/project-safari/zebra/lease"
	"github.com/project-safari/zebra/network"
	"github.com/stretchr/testify/assert"
)

func TestNewLease(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	l := getLease()
	assert.NotNil(l)
}

func TestActivate(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	l := getLease()
	assert.NotNil(l)

	assert.Equal(zebra.Inactive, l.Status)

	assert.Nil(l.Activate())
	assert.Equal(zebra.Active, l.Status)
	assert.False(l.ActivationTime.IsZero())
}

func TestDeactivate(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	l := getLease()
	assert.NotNil(l)

	assert.Equal(zebra.Inactive, l.Status)

	assert.Nil(l.Activate())
	assert.Equal(zebra.Active, l.Status)
	assert.False(l.ActivationTime.IsZero())

	l.Deactivate()
	assert.Equal(zebra.Inactive, l.Status)
}

func TestBadResources(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	l := getLease()
	assert.NotNil(l)

	res := getActiveRes()

	l.Resources.Add(res, res.GetType())

	assert.NotNil(l.Activate())

	res.GetStatus().Lease = zebra.Free

	assert.Nil(l.Activate())
}

func getLease() *lease.Lease {
	d, err := time.ParseDuration("4h")
	if err != nil {
		return nil
	}

	resources := []lease.ResourceReq{
		{
			Type:  "Server",
			Name:  "linux blah blah",
			Count: 2,
		},
		{
			Type:  "VM",
			Name:  "virtual",
			Count: 1,
		},
	}

	req := lease.Request{
		Duration:  d,
		Group:     "san-jose-building-15",
		Resources: resources,
	}

	return lease.NewLease(getUser(), req)
}

func getActiveRes() zebra.Resource {
	res := &network.VLANPool{
		BaseResource: *zebra.NewBaseResource("VLANPool", nil),
		RangeStart:   0,
		RangeEnd:     10,
	}

	status := zebra.DefaultStatus()
	status.Lease = zebra.Leased

	res.Status = status

	return res
}

func getUser() auth.User {
	return auth.User{
		NamedResource: zebra.NamedResource{
			BaseResource: *zebra.NewBaseResource("User", nil),
			Name:         "shravya",
		},
		Key:          nil,
		PasswordHash: "",
		Role:         nil,
	}
}
