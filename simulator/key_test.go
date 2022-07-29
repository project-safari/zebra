package simulator_test

import (
	"testing"

	"github.com/project-safari/zebra/auth"
	"github.com/stretchr/testify/assert"
)

func TestKey(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	key, err := auth.Generate()
	assert.Nil(err)
	assert.NotNil(key)

	assert.Nil(key.Save("user.key"))
}
