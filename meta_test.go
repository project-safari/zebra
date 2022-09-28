package zebra_test

import (
	"testing"

	"github.com/project-safari/zebra"
	"github.com/stretchr/testify/assert"
)

// Test function that verifies the validation function for the meta operations.
func TestMetaValidate(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	m := zebra.Meta{}
	assert.Equal(zebra.ErrIDEmpty, m.Validate())

	m.ID = "ab"
	assert.Equal(zebra.ErrIDShort, m.Validate())

	m.ID = "abcdefgh"
	assert.Equal(zebra.ErrNameEmpty, m.Validate())

	m.Name = "mock"
	assert.NotNil(m.Validate())

	m.Type = zebra.Type{Name: "dummy", Description: "dummy"}
	assert.Equal(zebra.ErrLabel, m.Validate())

	m = zebra.NewMeta(m.Type, "test_meta", "test_group", "test_owner")
	assert.Nil(m.Validate())

	m = zebra.NewMeta(m.Type, "", "test_group", "test_owner")
	assert.Nil(m.Validate())
}
