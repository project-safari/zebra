package zebra_test

import (
	"testing"

	"github.com/project-safari/zebra"
	"github.com/stretchr/testify/assert"
)

func TestCredentials(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	credentials := zebra.NewCredentials("")
	assert.Equal(zebra.ErrIDEmpty, credentials.Validate())

	assert.Equal(zebra.ErrUnknownCredential, credentials.Add("blah", "blah"))

	credentials = zebra.NewCredentials("id123")
	assert.NotNil(credentials.Validate())

	assert.Nil(credentials.Add("password", "a"))
	assert.Nil(credentials.Add("ssh-key", "test"))
	assert.NotNil(credentials.Validate())

	credentials.Keys["password"] = "abcdefghijklm"
	assert.NotNil(credentials.Validate())

	credentials.Keys["password"] = "ABCDEFGHIJKLM"
	assert.NotNil(credentials.Validate())

	credentials.Keys["password"] = "ABCDEFghijklm"
	assert.NotNil(credentials.Validate())

	credentials.Keys["password"] = "ABCDEFghijklm1"
	assert.NotNil(credentials.Validate())

	credentials.Keys["password"] = "properPass123$"
	assert.Nil(credentials.Validate())
}
