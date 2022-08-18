package lease //nolint:testpackage

import (
	"context"
	"testing"
	"time"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/network"
	"github.com/project-safari/zebra/status"
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

	l := getEmptyLease()
	assert.NotNil(l)

	assert.Equal(status.Inactive, l.Status.State())

	assert.Nil(l.Activate())
	assert.Equal(status.Active, l.Status.State())
	assert.False(l.ActivationTime.IsZero())
}

func TestDeactivate(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	l := getEmptyLease()
	assert.NotNil(l)

	assert.Equal(status.Inactive, l.Status.State())

	assert.Nil(l.Activate())
	assert.Equal(status.Active, l.Status.State())
	assert.False(l.ActivationTime.IsZero())

	l.Deactivate()
	assert.Equal(status.Inactive, l.Status.State())
}

func TestBadResources(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	l := getEmptyLease()
	l.Request = []*ResourceReq{
		{
			Type:  "VLANPool",
			Group: "sj-building-20",
			Name:  "blah blah give a name",
			Count: 1,
		},
	}
	assert.NotNil(l)

	res := getRes()

	assert.Nil(l.Request[0].Assign(res))
	assert.Nil(l.Activate())
}

func TestValidExpired(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	l := getEmptyLease()
	assert.NotNil(l)

	assert.Nil(l.Activate())
	assert.True(l.IsValid())
	assert.False(l.IsExpired())

	l.Deactivate()
	assert.False(l.IsValid())
	assert.True(l.IsExpired())
}

func TestValidate(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	l := getEmptyLease()
	assert.Nil(l.Validate(context.Background()))

	dur, err := time.ParseDuration("5h")
	assert.Nil(err)

	l.ActivationTime = time.Now().Add(dur)
	assert.NotNil(l.Validate(context.Background()))

	l.Request = nil
	assert.NotNil(l.Validate(context.Background()))

	l.Duration = dur
	assert.NotNil(l.Validate(context.Background()))
}

func TestRequestList(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	l := getEmptyLease()
	assert.Empty(l.RequestList())
}

func TestOwner(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	l := getEmptyLease()
	assert.Equal("test@zebra.project-safari.io", l.Owner())
}

func getEmptyLease() *Lease {
	d, err := time.ParseDuration("4h")
	if err != nil {
		return nil
	}

	return NewLease(getEmail(), d, make([]*ResourceReq, 0))
}

func getLease() *Lease {
	d, err := time.ParseDuration("4h")
	if err != nil {
		return nil
	}

	resources := []*ResourceReq{
		{
			Type:  "Server",
			Group: "san-jose-building-14",
			Name:  "linux blah blah",
			Count: 2,
		},
		{
			Type:  "VM",
			Group: "san-jose-building-18",
			Name:  "virtual",
			Count: 1,
		},
	}

	return NewLease(getEmail(), d, resources)
}

func getRes() zebra.Resource {
	res := &network.VLANPool{
		BaseResource: *zebra.NewBaseResource("VLANPool", nil),
		RangeStart:   0,
		RangeEnd:     10,
	}

	return res
}

func getEmail() string {
	return "test@zebra.project-safari.io"
}
