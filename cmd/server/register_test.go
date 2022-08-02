package main //nolint:testpackage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/project-safari/zebra/auth"
	"github.com/project-safari/zebra/cmd/herd/pkg"
	"github.com/project-safari/zebra/store"
	"github.com/stretchr/testify/assert"
)

func TestCreateNewUser(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "regteststore1"

	t.Cleanup(func() { os.RemoveAll(root) })

	key, _ := auth.Generate()
	user := createNewUser("testuser", "test@cisco.com", "bigword", key.Public())

	user.Labels = pkg.CreateLabels()
	user.Labels = pkg.GroupLabels(user.Labels, "groupSample")

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
	user.Labels = pkg.CreateLabels()
	user.Labels = pkg.GroupLabels(user.Labels, "groupSample")

	oldpassword := user.PasswordHash

	store := makeQueryStore(root, assert, user)
	assert.NotNil(store)
	err := deleteUser(user, store)
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

	user.Labels = pkg.CreateLabels()
	user.Labels = pkg.GroupLabels(user.Labels, "groupLabel")

	store := makeQueryStore(root, assert, user)
	assert.NotNil(store)
	err := deleteUser(user, store)
	assert.Nil(err)
	assert.Empty(store.Query().Resources)
}

type RData struct {
	Name     string            `json:"name"`
	Password string            `json:"password"`
	Email    string            `json:"email"`
	Key      *auth.RsaIdentity `json:"key"`
}

func newRData(name string, password string, email string, needKey bool) *RData {
	return &RData{
		Name:     name,
		Password: password,
		Email:    email,
		Key: func() *auth.RsaIdentity {
			if !needKey {
				return nil
			}

			key, _ := auth.Generate()
			pubKey := key.Public()

			return pubKey
		}(),
	}
}

func (r *RData) Body() io.ReadCloser {
	v, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}

	return ioutil.NopCloser(bytes.NewBuffer(v))
}

func TestRegistry(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "regteststore4"

	t.Cleanup(func() { os.RemoveAll(root) })

	req, err := http.NewRequest("POST", "/register", nil)
	assert.Nil(err)
	assert.NotNil(req)

	// Invalid Context
	h := registerAdapter()
	rr := httptest.NewRecorder()
	handler := h(nil)

	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusInternalServerError, rr.Code)

	data := newRData("testuser2", "secrect", "myemail@domain", true)

	user := createNewUser("testuser", "testuser@domain", "secrect", data.Key)
	user.Labels = pkg.CreateLabels()
	user.Labels = pkg.GroupLabels(user.Labels, "groupSample")

	resources := NewResourceAPI(store.DefaultFactory())
	resources.Store = makeQueryStore(root, assert, user)

	ctx := context.WithValue(context.Background(), ResourcesCtxKey, resources)

	req, err = http.NewRequestWithContext(ctx, "POST", "/register", nil)
	req.Body = data.Body()

	assert.Nil(err)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusCreated, rr.Code)

	testForward(assert, h)
}

func TestNoKeyUser(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "regteststore5"

	t.Cleanup(func() { os.RemoveAll(root) })

	key, _ := auth.Generate()
	user := createNewUser("testuser", "test@cisco.com", "bigword", key.Public())
	resources := NewResourceAPI(store.DefaultFactory())
	resources.Store = makeQueryStore(root, assert, user)

	ctx := context.WithValue(context.Background(), ResourcesCtxKey, resources)

	req, err := http.NewRequestWithContext(ctx, "POST", "/register", nil)
	assert.Nil(err)
	assert.NotNil(req)

	data := newRData("testuser2", "secrect", "myemail@domain", false)
	req.Body = data.Body()

	h := registerAdapter()
	rr := httptest.NewRecorder()
	handler := h(nil)

	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusInternalServerError, rr.Code)
}

func TestSameUser(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "regteststore6"

	t.Cleanup(func() { os.RemoveAll(root) })

	data := newRData("testuser", "secrect", "test@cisco123.com", true)
	user := createNewUser("testuser", "test@cisco123.com", "bigword", data.Key)
	user.Labels = pkg.CreateLabels()
	user.Labels = pkg.GroupLabels(user.Labels, "sampleGroup")

	resources := NewResourceAPI(store.DefaultFactory())
	resources.Store = makeQueryStore(root, assert, user)

	ctx := context.WithValue(context.Background(), ResourcesCtxKey, resources)

	req, err := http.NewRequestWithContext(ctx, "POST", "/register", nil)
	assert.Nil(err)
	assert.NotNil(req)

	req.Body = data.Body()

	h := registerAdapter()
	rr := httptest.NewRecorder()
	handler := h(nil)

	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusForbidden, rr.Code)
}

func TestBadRequest(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "regteststore7"

	t.Cleanup(func() { os.RemoveAll(root) })

	key, _ := auth.Generate()
	registryData := fmt.Sprintf("{\"user\":\"%s\"\"password\":\"%s\"}", "Bad", "Request")

	v, err := json.Marshal(registryData)
	assert.Nil(err)
	assert.NotEmpty(v)

	user := createNewUser("testuser", "test@cisco.com", "bigword", key.Public())
	resources := NewResourceAPI(store.DefaultFactory())
	resources.Store = makeQueryStore(root, assert, user)

	ctx := context.WithValue(context.Background(), ResourcesCtxKey, resources)

	req, err := http.NewRequestWithContext(ctx, "POST", "/register", nil)
	assert.Nil(err)
	assert.NotNil(req)

	req.Body = ioutil.NopCloser(bytes.NewBuffer(v))

	h := registerAdapter()
	rr := httptest.NewRecorder()
	handler := h(nil)

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

	user.Labels = pkg.CreateLabels()
	user.Labels = pkg.GroupLabels(user.Labels, "sampleGroup")

	user2 := createNewUser("testuser", "test@cisco.com", "bigword", nil)

	user2.Labels = pkg.CreateLabels()
	user2.Labels = pkg.GroupLabels(user2.Labels, "groupUser2")

	store := makeQueryStore(root, assert, user)
	assert.NotNil(store)
	err := deleteUser(user2, store)
	assert.NotNil(err)
}
