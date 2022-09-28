package auth_test

import (
	"testing"

	"github.com/project-safari/zebra/auth"
	"github.com/stretchr/testify/assert"
)

func TestClaims(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	crud, err := auth.NewPriv("", true, true, true, true)
	assert.NotNil(crud)
	assert.Nil(err)

	role := &auth.Role{"boy", []*auth.Priv{crud}}
	claims := auth.NewClaims("zebra", "adam", role, "email@domain")

	assert.True(claims.Create("anything"))
	assert.True(claims.Read("anything"))
	assert.True(claims.Write("anything"))
	assert.True(claims.Update("anything"))
	assert.True(claims.Delete("anything"))

	jwt := claims.JWT("abracadabra")
	assert.NotEmpty(jwt)

	jwtClaims, err := auth.FromJWT(jwt, "abracadabra")
	assert.NotNil(jwtClaims)
	assert.Nil(err)

	jwtClaims, err = auth.FromJWT(jwt, "blahblahblah")
	assert.Nil(jwtClaims)
	assert.NotNil(err)
}
