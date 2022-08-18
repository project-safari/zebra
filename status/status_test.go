package status_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/project-safari/zebra/status"
	"github.com/stretchr/testify/assert"
)

func TestDefaultStatus(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	s := status.DefaultStatus()
	assert.NotNil(s)
	assert.Equal(status.None, s.Fault())
	assert.Equal(status.Free, s.Lease())
	assert.Equal("", s.UsedBy())
	assert.Equal(status.Active, s.State())
	assert.True(s.CreatedTime().Before(time.Now()))
}

func TestValidateStatus(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	s := status.NewStatus(8, 8, 8, "", time.Now().AddDate(2, 1, 3))
	ctx := context.Background()

	assert.NotNil(s.Validate(ctx))
	s = status.NewStatus(status.Critical, 8, 8, "", time.Now().AddDate(2, 1, 3))

	assert.NotNil(s.Validate(ctx))
	s = status.NewStatus(status.Critical, status.Free, 8, "", time.Now().AddDate(2, 1, 3))

	assert.NotNil(s.Validate(ctx))
	s = status.NewStatus(status.Critical, status.Free, status.Inactive, "", time.Now().AddDate(2, 1, 3))

	assert.NotNil(s.Validate(ctx))
	s = status.NewStatus(status.None, status.Free, status.Inactive, "", time.Now())

	assert.Nil(s.Validate(ctx))
}

func TestString(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	var f status.FaultState = 8

	var l status.LeaseState = 8

	var s status.ActivityState = 8

	assert.Equal("unknown", f.String())
	assert.Equal("unknown", l.String())
	assert.Equal("unknown", s.String())

	f = status.None
	assert.Equal("none", f.String())

	l = status.Free
	assert.Equal("free", l.String())

	s = status.Active
	assert.Equal("active", s.String())
}

func TestMarshal(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	status := status.DefaultStatus()
	bytes, err := json.Marshal(status)
	assert.Nil(err)
	assert.NotNil(bytes)
}

func TestUnmarshal(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	body := `8`

	var f *status.FaultState

	var l *status.LeaseState

	var s *status.ActivityState

	status := new(status.Status)

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

	s := status.DefaultStatus()
	assert.Equal(status.Active, s.State())

	s.Deactivate()
	assert.Equal(status.Inactive, s.State())

	s.Activate()
	assert.Equal(status.Active, s.State())
}

func TestSetLeaseState(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	s := status.DefaultStatus()
	assert.Equal(status.Free, s.Lease())

	assert.Nil(s.SetLeased())
	assert.Equal(status.Leased, s.Lease())

	assert.NotNil(s.SetLeased())

	assert.Nil(s.SetFree())
	assert.NotNil(s.SetFree())
}

func TestSetOwner(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	s := status.DefaultStatus()
	assert.Equal("", s.UsedBy())

	s.SetUser("Shravya")
	assert.Equal("Shravya", s.UsedBy())
}
