package main //nolint:testpackage

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/project-safari/zebra/auth"
	"github.com/stretchr/testify/assert"
)

func TestCreateNewUser(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "regteststore1"

	t.Cleanup(func() { os.RemoveAll(root) })

	key, _ := auth.Generate()
	user := createNewUser("testuser", "test@cisco.com", "bigword", key.Public())

	assert.NotNil(user)
	assert.NotNil(user.BaseResource)
	assert.NotNil(user.Role)
	assert.NotNil(user.Key)
	assert.NotNil(user.PasswordHash)
}

func TestUpdateUser(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "regteststore2"

	t.Cleanup(func() { os.RemoveAll(root) })

	key, _ := auth.Generate()
	user := createNewUser("testuser", "test@cisco.com", "bigword", key.Public())
	oldpassword := user.PasswordHash

	store := makeQueryStore(root, assert, user)
	assert.NotNil(store)
	err := changeActiveStatus(user, store, false)
	assert.Nil(err)

	changePassword(user, "newpassword")
	assert.NotEqual(oldpassword, user.PasswordHash)
}

func TestDeleteUser(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "regteststore3"

	t.Cleanup(func() { os.RemoveAll(root) })

	key, _ := auth.Generate()
	user := createNewUser("testuser", "test@cisco.com", "bigword", key.Public())

	store := makeQueryStore(root, assert, user)
	assert.NotNil(store)
	err := changeActiveStatus(user, store, false)
	assert.Nil(err)
	assert.Empty(store.Query().Resources)
}

func TestRegistry(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "regteststore4"

	t.Cleanup(func() { os.RemoveAll(root) })

	req, err := http.NewRequest("POST", "/register", nil)
	assert.Nil(err)
	assert.NotNil(req)

	key, _ := auth.Generate()
	pubKey := key.Public()
	registryData := &struct {
		Name     string            `json:"name"`
		Password string            `json:"password"`
		Email    string            `json:"email"`
		Key      *auth.RsaIdentity `json:"key"`
	}{Name: "testuser2",
		Password: "secrect",
		Email:    "myemail@domain",
		Key:      pubKey}

	v, err := json.Marshal(registryData)
	assert.Nil(err)
	assert.NotEmpty(v)
	req.Body = ioutil.NopCloser(bytes.NewBuffer(v))

	user := createNewUser("testuser", "test@cisco.com", "bigword", pubKey)
	store := makeQueryStore(root, assert, user)
	h := handleRegister(setupLogger(nil), store)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h(w, r, nil)
	})

	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusCreated, rr.Code)
}

func TestNoKeyUser(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "regteststore5"

	t.Cleanup(func() { os.RemoveAll(root) })

	req, err := http.NewRequest("POST", "/register", nil)
	assert.Nil(err)
	assert.NotNil(req)

	key, _ := auth.Generate()
	registryData := &struct {
		Name     string            `json:"name"`
		Password string            `json:"password"`
		Email    string            `json:"email"`
		Key      *auth.RsaIdentity `json:"key"`
	}{Name: "testuser2",
		Password: "secrect",
		Email:    "myemail@domain",
	}

	v, err := json.Marshal(registryData)
	assert.Nil(err)
	assert.NotEmpty(v)
	req.Body = ioutil.NopCloser(bytes.NewBuffer(v))

	user := createNewUser("testuser", "test@cisco.com", "bigword", key.Public())
	store := makeQueryStore(root, assert, user)
	h := handleRegister(setupLogger(nil), store)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h(w, r, nil)
	})

	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusInternalServerError, rr.Code)
}

func TestSameUser(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "regteststore6"

	t.Cleanup(func() { os.RemoveAll(root) })

	req, err := http.NewRequest("POST", "/register", nil)
	assert.Nil(err)
	assert.NotNil(req)

	key, _ := auth.Generate()
	pubKey := key.Public()
	registryData := &struct {
		Name     string            `json:"name"`
		Password string            `json:"password"`
		Email    string            `json:"email"`
		Key      *auth.RsaIdentity `json:"key"`
	}{Name: "testuser",
		Key:      pubKey,
		Email:    "test@cisco123.com",
		Password: "secrect",
	}

	v, err := json.Marshal(registryData)
	assert.Nil(err)
	assert.NotEmpty(v)
	req.Body = ioutil.NopCloser(bytes.NewBuffer(v))

	user := createNewUser("testuser", "test@cisco123.com", "bigword", key.Public())
	store := makeQueryStore(root, assert, user)
	h := handleRegister(setupLogger(nil), store)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h(w, r, nil)
	})

	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusForbidden, rr.Code)
}

func TestBadRequest(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "regteststore7"

	t.Cleanup(func() { os.RemoveAll(root) })

	req, err := http.NewRequest("POST", "/register", nil)
	assert.Nil(err)
	assert.NotNil(req)

	key, _ := auth.Generate()
	registryData := fmt.Sprintf("{\"user\":\"%s\"\"password\":\"%s\"}", "Bad", "Request")

	v, err := json.Marshal(registryData)
	assert.Nil(err)
	assert.NotEmpty(v)
	req.Body = ioutil.NopCloser(bytes.NewBuffer(v))

	user := createNewUser("testuser", "test@cisco.com", "bigword", key.Public())
	store := makeQueryStore(root, assert, user)
	h := handleRegister(setupLogger(nil), store)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h(w, r, nil)
	})

	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusBadRequest, rr.Code)
}

func TestDelete1User(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "regteststore8"

	t.Cleanup(func() { os.RemoveAll(root) })

	key, _ := auth.Generate()
	user := createNewUser("testuser", "test@cisco.com", "bigword", key.Public())
	user2 := createNewUser("testuser", "test@cisco.com", "bigword", nil)

	store := makeQueryStore(root, assert, user)
	assert.NotNil(store)
	err := changeActiveStatus(user2, store, false)
	assert.NotNil(err)
}
