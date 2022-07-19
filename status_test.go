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
	assert.Equal("none", s.Fault)
	assert.Equal("free", s.Lease)
	assert.Equal("", s.UsedBy)
	assert.Equal("active", s.State)
}

func TestValidateStatus(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	s := zebra.Status{
		Fault:       "random1",
		Lease:       "random2",
		UsedBy:      "random3",
		State:       "random4",
		CreatedTime: time.Now().AddDate(2, 1, 3),
	}
	ctx := context.Background()

	assert.NotNil(s.Validate(ctx))
	s.Fault = "critical"

	assert.NotNil(s.Validate(ctx))
	s.Lease = "free"

	assert.NotNil(s.Validate(ctx))
	s.State = "inactive"

	assert.NotNil(s.Validate(ctx))
	s.CreatedTime = time.Now()

	assert.Nil(s.Validate(ctx))
}
