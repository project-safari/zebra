package lease //nolint:testpackage

import (
	"testing"
	"time"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/auth"
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

	l := getEmptyLease()
	assert.NotNil(l)

	assert.Equal(zebra.Inactive, l.Status.State)

	assert.Nil(l.Activate())
	assert.Equal(zebra.Active, l.Status.State)
	assert.False(l.activationTime.IsZero())
}

func TestDeactivate(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	l := getEmptyLease()
	assert.NotNil(l)

	assert.Equal(zebra.Inactive, l.Status.State)

	assert.Nil(l.Activate())
	assert.Equal(zebra.Active, l.Status.State)
	assert.False(l.activationTime.IsZero())

	l.Deactivate()
	assert.Equal(zebra.Inactive, l.Status.State)
}

func TestBadResources(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	l := getEmptyLease()
	assert.NotNil(l)

	res := getRes()

	l.Assign(res)
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

func getEmptyLease() *Lease {
	d, err := time.ParseDuration("4h")
	if err != nil {
		return nil
	}

	return NewLease(getUser(), d, make(map[*zebra.Group]*ResourceReq, 0))
}

func getLease() *Lease {
	d, err := time.ParseDuration("4h")
	if err != nil {
		return nil
	}

	resources := map[*zebra.Group]*ResourceReq{
		zebra.NewGroup("san-jose-building-14"): {
			Type:  "Server",
			Group: "san-jose-building-14",
			Name:  "linux blah blah",
			Count: 2,
		},
		zebra.NewGroup("san-jose-building-18"): {
			Type:  "VM",
			Group: "san-jose-building-18",
			Name:  "virtual",
			Count: 1,
		},
	}

	return NewLease(getUser(), d, resources)
}

func getRes() zebra.Resource {
	res := &network.VLANPool{
		BaseResource: *zebra.NewBaseResource("VLANPool", nil),
		RangeStart:   0,
		RangeEnd:     10,
	}

	res.Status.Lease = zebra.Free

	return res
}

func getUser() auth.User {
	return auth.User{
		NamedResource: zebra.NamedResource{
			BaseResource: *zebra.NewBaseResource("User", nil),
			Name:         "shravya",
		},
		Email:        "shravya@cisco.com",
		Key:          nil,
		PasswordHash: "",
		Role:         nil,
	}
}
