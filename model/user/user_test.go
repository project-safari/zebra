package user_test

import (
	"context"
	"testing"

	"github.com/project-safari/zebra/auth"
	"github.com/project-safari/zebra/model/user"
	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) { //nolint:funlen
	t.Parallel()

	assert := assert.New(t)
	all, e := auth.NewPriv("", true, true, true, true)
	assert.Nil(e)
	assert.NotNil(all)

	admin := &auth.Role{Name: "admin", Privileges: []*auth.Priv{all}}

	writeAll, e := auth.NewPriv("", true, false, true, true)
	assert.Nil(e)
	assert.NotNil(writeAll)

	readAll, e := auth.NewPriv("", false, true, false, false)
	assert.Nil(e)
	assert.NotNil(all)

	rwOne, e := auth.NewPriv("eden", true, true, true, true)
	assert.Nil(e)
	assert.NotNil(rwOne)

	roleUser := &auth.Role{Name: "user", Privileges: []*auth.Priv{readAll, rwOne}}

	godKey, err := auth.Generate()
	assert.Nil(err)
	assert.NotNil(godKey)

	ctx := context.Background()

	uType := user.Type()
	god, ok := user.Empty().(*user.User)
	assert.True(ok)

	god.Meta.Name = "almighty"
	god.Meta.ID = "00000000000001"
	god.Meta.Type = uType
	assert.NotNil(god.Validate(ctx))

	god.Key = godKey
	assert.NotNil(god.Validate(ctx))

	god.Role = admin
	assert.NotNil(god.Validate(ctx))

	god.PasswordHash = user.HashPassword("youhaveachoice")
	god.Email = "god@heaven.com"
	assert.NotNil(god.Validate(ctx)) // it will have an error because it does not have group label, so not nil.

	adamKey, err := auth.Generate()
	assert.Nil(err)
	assert.NotNil(adamKey)

	adam := new(user.User)
	adam.Meta.Name = "adam"
	adam.Meta.ID = "00000000000003"
	adam.Meta.Type = uType
	adam.Key = adamKey
	adam.Role = roleUser
	adam.PasswordHash = user.HashPassword("iloveeve")

	eveKey, err := auth.Generate()
	assert.Nil(err)
	assert.NotNil(eveKey)

	eve := new(user.User)
	eve.Meta.Name = "eve"
	eve.Meta.ID = "00000000000004"
	eve.Meta.Type = uType
	eve.Key = eveKey
	eve.Role = roleUser
	eve.PasswordHash = user.HashPassword("iloveadam")

	token, err := godKey.Sign([]byte(god.Email))
	assert.Nil(err)
	assert.NotEmpty(token)

	assert.Nil(god.Authenticate(string(token)))
	assert.Nil(god.AuthenticatePassword("youhaveachoice"))

	assert.True(god.Create("universe"))
	assert.True(god.Read("universe"))
	assert.True(god.Delete("universe"))
	assert.True(god.Update("universe"))
	assert.True(god.Write("universe"))

	assert.False(adam.Write("something"))
	assert.False(adam.Delete("something"))
	assert.False(eve.Create("universe"))
	assert.False(eve.Delete("universe"))
	assert.True(adam.Create("eden"))
	assert.True(eve.Write("eden"))
	assert.True(eve.Update("eden"))
	assert.False(eve.Update("universe"))

	assert.True(eve.Write("eden"))
	assert.True(eve.Update("eden"))
	assert.False(eve.Update("universe"))

	newUser := user.NewUser("eve", "eve@email.com", "eve123", eveKey, roleUser)
	assert.NotNil(newUser)
	assert.Nil(newUser.AuthenticatePassword("eve123"))
	assert.NotNil(newUser.AuthenticatePassword("adam123"))
}
