package lease //nolint:testpackage

import (
	"context"
	"testing"
	"time"

	"github.com/project-safari/zebra"
	"github.com/stretchr/testify/assert"
)

// Test function for creating a new lease.
func TestNewLease(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	l := getLease()
	assert.NotNil(l)
}

// Test function for activation of a lease.
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

// Test function for deactivation of a lease.
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

// Test function for leases of bad - incorrect - leases.
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

// Test function for valid but expired leases.
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

// Test function for validation of leases and their request(s).
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

// Test functuon that verifies the list of resources in a lease request
// and ensures it has the appropriate length.
func TestRequestList(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	l := getEmptyLease()
	assert.Empty(l.RequestList())
}

// Test functuon that verifies that the identity the owner
// of a resource and its lease (request) is as expected.
func TestOwner(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	l := getEmptyLease()
	assert.Equal("tester@quality.com", l.Owner())
}

// Mock function that gets an empy lease request to be used in tests.
// It returns a pointer to the Lease struct.
func getEmptyLease() *Lease {
	d, err := time.ParseDuration("4h")
	if err != nil {
		return nil
	}

	return NewLease("tester@quality.com", d, make([]*ResourceReq, 0))
}

// Mock function that gets a lease request to be used in tests.
// It returns a pointer to the Lease struct.
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

	return NewLease("tester@quality.com", d, resources)
}

// Mock function that a resource to be used in tests.
// It returns a zebra.Resource.
func getRes() zebra.Resource {
	return zebra.NewBaseResource(
		zebra.Type{Name: "dummy", Description: "dummy"},
		"test_res", "tester", "test_group")
}
