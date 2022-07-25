package zebra_test

import (
	"context"
	"testing"
	"time"

	"github.com/project-safari/zebra"
	"github.com/stretchr/testify/assert"
)

func TestDefaultStatus(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	s := zebra.DefaultStatus()
	assert.NotNil(s)
	assert.Equal(zebra.None, s.Fault())
	assert.Equal(zebra.Free, s.Lease())
	assert.Equal("", s.UsedBy())
	assert.Equal(zebra.Active, s.State())
	assert.True(s.CreatedTime().Before(time.Now()))
}

func TestValidateStatus(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	s := zebra.NewStatus(8, 8, 8, "", time.Now().AddDate(2, 1, 3))
	ctx := context.Background()

	assert.NotNil(s.Validate(ctx))
	s = zebra.NewStatus(zebra.Critical, 8, 8, "", time.Now().AddDate(2, 1, 3))

	assert.NotNil(s.Validate(ctx))
	s = zebra.NewStatus(zebra.Critical, zebra.Free, 8, "", time.Now().AddDate(2, 1, 3))

	assert.NotNil(s.Validate(ctx))
	s = zebra.NewStatus(zebra.Critical, zebra.Free, zebra.Inactive, "", time.Now().AddDate(2, 1, 3))

	assert.NotNil(s.Validate(ctx))
	s = zebra.NewStatus(zebra.None, zebra.Free, zebra.Inactive, "", time.Now())

	assert.Nil(s.Validate(ctx))
}

func TestString(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	var f zebra.FaultState = 8

	var l zebra.LeaseState = 8

	var s zebra.ActivityState = 8

	assert.Equal("unknown", f.String())
	assert.Equal("unknown", l.String())
	assert.Equal("unknown", s.String())

	f = zebra.None
	assert.Equal("none", f.String())

	l = zebra.Free
	assert.Equal("free", l.String())

	s = zebra.Active
	assert.Equal("active", s.String())
}

func TestMarshal(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	f := zebra.None
	bytes, err := f.MarshalText()
	assert.Nil(err)
	assert.Equal([]byte(`none`), bytes)

	l := zebra.Leased
	bytes, err = l.MarshalText()
	assert.Nil(err)
	assert.Equal([]byte(`leased`), bytes)

	s := zebra.Inactive
	bytes, err = s.MarshalText()
	assert.Nil(err)
	assert.Equal([]byte(`inactive`), bytes)
}

func TestUnmarshal(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	body := `8`

	var f *zebra.FaultState

	var l *zebra.LeaseState

	var s *zebra.ActivityState

	status := new(zebra.Status)

	assert.NotNil(f.UnmarshalText([]byte(body)))
	assert.NotNil(l.UnmarshalText([]byte(body)))
	assert.NotNil(s.UnmarshalText([]byte(body)))

	// Test bad unmarshal, should fail
	body = `{"fault":"none","lease":"blahblah"}`
	assert.NotNil(status.UnmarshalJSON([]byte(body)))

	body = `{"fault":"none","lease":"free","state":"active","usedBy":"","createdTime":"2018-09-22T19:42:31+07:00"}`
	assert.Nil(status.UnmarshalJSON([]byte(body)))
}

func TestActivateDeactivate(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	s := zebra.DefaultStatus()
	assert.Equal(zebra.Active, s.State())

	s.Deactivate()
	assert.Equal(zebra.Inactive, s.State())

	s.Activate()
	assert.Equal(zebra.Active, s.State())
}

func TestSetLeaseState(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	s := zebra.DefaultStatus()
	assert.Equal(zebra.Free, s.Lease())

	assert.Nil(s.SetLeased())
	assert.Equal(zebra.Leased, s.Lease())

	assert.NotNil(s.SetLeased())
}

func TestSetOwner(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	s := zebra.DefaultStatus()
	assert.Equal("", s.UsedBy())

	s.SetOwner("Shravya")
	assert.Equal("Shravya", s.UsedBy())
}
