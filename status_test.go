package zebra_test

import (
	"testing"

	"github.com/project-safari/zebra"
	"github.com/stretchr/testify/assert"
)

func TestDefaultStatus(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	s := zebra.DefaultStatus()
	assert.NotNil(s)
	assert.Equal(zebra.None, s.Fault)
	assert.Equal(zebra.Free, s.LeaseStatus)
	assert.Equal("", s.UsedBy)
	assert.Equal(zebra.Inactive, s.State)
}

func TestValidateStatus(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	s := zebra.Status{
		Fault:       8,
		LeaseStatus: 8,
		UsedBy:      "",
		State:       8,
	}
	assert.NotNil(s.Validate())
	s.Fault = zebra.Critical

	assert.NotNil(s.Validate())
	s.LeaseStatus = zebra.Free

	assert.NotNil(s.Validate())
	s.State = zebra.Inactive

	assert.Nil(s.Validate())
}

func TestString(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	var f zebra.Fault = 8

	var l zebra.LeaseStatus = 8

	var s zebra.State = 8

	assert.Equal("unknown", f.String())
	assert.Equal("unknown", l.String())
	assert.Equal("unknown", s.String())
}

func TestUnmarshal(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	f := zebra.None
	assert.Equal("none", f.String())
	b, err := f.MarshalText()
	assert.Nil(err)
	assert.NotEmpty(b)
	assert.Nil(f.UnmarshalText(b))
	assert.Equal(zebra.None, f)

	l := zebra.Free
	assert.Equal("free", l.String())
	b, err = l.MarshalText()
	assert.Nil(err)
	assert.NotEmpty(b)
	assert.Nil(l.UnmarshalText(b))
	assert.Equal(zebra.Free, l)

	s := zebra.Active
	assert.Equal("active", s.String())
	b, err = s.MarshalText()
	assert.Nil(err)
	assert.NotEmpty(b)
	assert.Nil(s.UnmarshalText(b))
	assert.Equal(zebra.Active, s)

	assert.Equal(zebra.ErrLeaseStatus, l.UnmarshalText([]byte("zzz")))
	assert.Equal(zebra.ErrFault, f.UnmarshalText([]byte("zzz")))
	assert.Equal(zebra.ErrState, s.UnmarshalText([]byte("zzz")))
}

func TestStatusChecker(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	d, _ := dummyType()
	res := zebra.NewBaseResource(d, "dummy", "dummy", "dummy")

	// lease status 1 means not leasable.
	res.Status.LeaseStatus = 2

	assert.NotNil(res.Status.LeaseStatus.CanLease())

	// lease status 1 means not free.
	res.Status.LeaseStatus = 1

	assert.NotNil(res.Status.LeaseStatus.IsFree())

	// lease status 0 means free and available.
	res.Status.LeaseStatus = 0

	assert.Nil(res.Status.LeaseStatus.CanLease())
	assert.Nil(res.Status.LeaseStatus.IsFree())
}
