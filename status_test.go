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
	assert.Equal(zebra.None, s.Fault)
	assert.Equal(zebra.Free, s.Lease)
	assert.Equal("", s.UsedBy)
	assert.Equal(zebra.Active, s.State)
}

func TestValidateStatus(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	s := zebra.Status{
		Fault:       8,
		Lease:       8,
		UsedBy:      "",
		State:       8,
		CreatedTime: time.Now().AddDate(2, 1, 3),
	}
	ctx := context.Background()

	assert.NotNil(s.Validate(ctx))
	s.Fault = zebra.Critical

	assert.NotNil(s.Validate(ctx))
	s.Lease = zebra.Free

	assert.NotNil(s.Validate(ctx))
	s.State = zebra.Inactive

	assert.NotNil(s.Validate(ctx))
	s.CreatedTime = time.Now()

	assert.Nil(s.Validate(ctx))
}

func TestString(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	var f zebra.Fault = 8

	var l zebra.Lease = 8

	var s zebra.State = 8

	assert.Equal("unknown", f.String())
	assert.Equal("unknown", l.String())
	assert.Equal("unknown", s.String())
}

func TestUnmarshal(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	body := `8`

	var f *zebra.Fault

	var l *zebra.Lease

	var s *zebra.State

	assert.NotNil(f.UnmarshalText([]byte(body)))
	assert.NotNil(l.UnmarshalText([]byte(body)))
	assert.NotNil(s.UnmarshalText([]byte(body)))
}
