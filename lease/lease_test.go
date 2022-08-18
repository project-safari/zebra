package lease //nolint:testpackage

import (
	"context"
	"testing"
	"time"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/auth"
	"github.com/project-safari/zebra/network"
	"github.com/stretchr/testify/assert"
)

// Test for getting a lease.
func TestNewLease(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	l := getLease()
	assert.NotNil(l)
}

// Test for activation of a lease.
// This function stats with an empty lease and then tests successfulness of activating it.
func TestActivate(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	l := getEmptyLease()
	assert.NotNil(l)

	assert.Equal(zebra.Inactive, l.Status.State)

	assert.Nil(l.Activate())
	assert.Equal(zebra.Active, l.Status.State)
	assert.False(l.ActivationTime.IsZero())
}

// Test for deactivation of a lease.
// This function stats with an empty lease, activates it and then tests successfulness of deactivating it.
func TestDeactivate(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	l := getEmptyLease()
	assert.NotNil(l)

	assert.Equal(zebra.Inactive, l.Status.State)

	assert.Nil(l.Activate())
	assert.Equal(zebra.Active, l.Status.State)
	assert.False(l.ActivationTime.IsZero())

	l.Deactivate()
	assert.Equal(zebra.Inactive, l.Status.State)
}

// Test lease with incorrect resources.
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

// Tests for valid and expired lease.
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

// Tests for validation.
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

// Tests for request list.
func TestRequestList(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	l := getEmptyLease()
	assert.Empty(l.RequestList())
}

// Tests the owner of a lease.
func TestOwner(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	l := getEmptyLease()
	assert.Equal("shravya@cisco.com", l.Owner())
}

// Mock function to make a new empty lease to use in tests.
func getEmptyLease() *Lease {
	d, err := time.ParseDuration("4h")
	if err != nil {
		return nil
	}

	return NewLease(getUser().Email, d, make([]*ResourceReq, 0))
}

// Mock function to make a new lease to use in tests.
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

	return NewLease(getUser().Email, d, resources)
}

// Mock function to create a new resource to use in tests.
func getRes() zebra.Resource {
	res := &network.VLANPool{
		BaseResource: *zebra.NewBaseResource("VLANPool", nil),
		RangeStart:   0,
		RangeEnd:     10,
	}

	return res
}

// Mock function to create a new user to use in tests.
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
