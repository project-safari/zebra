package zebra_test

import (
	"testing"

	"github.com/project-safari/zebra"
	"github.com/stretchr/testify/assert"
)

func TestSecret(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	s := zebra.Secret{}
	message := []byte(`{"password":"strongPassword"}`)

	err := s.UnmarshalText(message)
	assert.Nil(err)

	b, err := s.MarshalText()
	assert.NotNil(b)
	assert.Nil(err)
}
