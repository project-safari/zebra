package main //nolint:testpackage

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/auth"
	"github.com/project-safari/zebra/store"
	"github.com/stretchr/testify/assert"
)

var (
	Email    = "test@cisco.com"
	Name     = "testuser"
	Password = "bigword"
)

func makeStore(root string, assert *assert.Assertions, user *NewUser) zebra.Store {
	factory := initTypes()
	assert.NotNil(factory)

	store := store.NewResourceStore(root, factory)
	if store.Initialize() != nil || store.Create(user) != nil {
		return nil
	}

	return store
}

func TestCreateNewUser(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "teststore1"

	t.Cleanup(func() { os.RemoveAll(root) })

	key, _ := auth.Generate()
	user := createNewUser(Email, Name, key.Public(), Password)

	assert.NotNil(user)
	assert.NotNil(user.BaseResource)
	assert.NotNil(user.User)
	assert.NotNil(user.Role)
	assert.NotNil(user.Key)
	assert.NotNil(user.PasswordHash)

	//Not Part of the test just using for debugging
	factory := initTypes()
	assert.NotNil(factory)

	store := store.NewResourceStore(root, factory)
	assert.Nil(store.Initialize())
	assert.Nil(store.Create(user))
}

func TestUpdateUser(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "teststore2"

	t.Cleanup(func() { os.RemoveAll(root) })

	key, _ := auth.Generate()
	user := createNewUser(Email, Name, key.Public(), Password)
	oldpassword := user.PasswordHash
	assert.Nil(user.Role)

	role := &auth.Role{
		Name:       "rolling",
		Privileges: nil,
	}
	user.createRole(role)

	assert.NotNil(user.Role)
	store := makeStore(root, assert, user)

	user.changeActiveStatus(store, false)

	user.changePassword("newpassword")
	assert.NotEqual(oldpassword, user.PasswordHash)

	err := user.changePubKey()
	assert.Nil(err)
}

func TestDeleteUser(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "teststore3"

	t.Cleanup(func() { os.RemoveAll(root) })

	key, _ := auth.Generate()
	user := createNewUser(Email, Name, key.Public(), Password)

	store := makeStore(root, assert, user)
	assert.NotNil(store)
	//user.changeActiveStatus(store, false)

	//assert.Nil(findUser(store, user.Name))

}

func TestRegistry(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "teststore4"

	t.Cleanup(func() { os.RemoveAll(root) })

	req, err := http.NewRequest("POST", "/register", nil)
	assert.Nil(err)
	assert.NotNil(req)

	key, _ := auth.Generate()
	v := fmt.Sprintf("{\"user\":\"%s\",\"password\":\"%s\"\"email\":\"%s\",\"key\":\"%p\"}", Name, Password, Email, key)
	req.Body = ioutil.NopCloser(bytes.NewBufferString(v))

	user := createNewUser(Email, Name, key.Public(), Password)
	store := makeStore(root, assert, user)
	h := handleRegister(setupLogger(nil), store)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h(w, r, nil)
	})

	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusCreated, rr.Code)
}
