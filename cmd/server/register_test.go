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
	"github.com/project-safari/zebra/model"
	"github.com/project-safari/zebra/model/user"
	"github.com/stretchr/testify/assert"
)

// Test function for creating a new user.
func TestCreateNewUser(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "test_create_new_user"

	defer func() { os.RemoveAll(root) }()

	key, _ := auth.Generate()
	testUser := user.NewUser("testuser", "test@cisco.com", "bigword", key.Public(), DefaultRole())

	assert.NotNil(testUser)
	assert.NotNil(testUser.Meta.Name)
	assert.NotNil(testUser.Role)
	assert.NotNil(testUser.Key)
	assert.NotNil(testUser.PasswordHash)
}

func TestUpdateUser(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "test_update_user"

	defer func() { os.RemoveAll(root) }()

	key, _ := auth.Generate()
	testUser := user.NewUser("testuser", "test@cisco.com", "bigword", key.Public(), DefaultRole())

	oldpassword := testUser.PasswordHash

	store := makeQueryStore(root, assert, testUser)
	assert.NotNil(store)
	err := deleteUser(testUser, store)
	assert.Nil(err)

	changePassword(testUser, "newpassword")
	assert.NotEqual(oldpassword, testUser.PasswordHash)
}

func TestDeleteUser(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "test_delete_user"

	defer func() { os.RemoveAll(root) }()

	key, _ := auth.Generate()
	testUser := user.NewUser("testuser", "test@cisco.com", "bigword", key.Public(), DefaultRole())

	store := makeQueryStore(root, assert, testUser)
	assert.NotNil(store)
	err := deleteUser(testUser, store)
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

func TestRegister(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "test_register"

	defer func() { os.RemoveAll(root) }()

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
	testUser := user.NewUser("testuser", "testuser@domain", "secrect", data.Key, DefaultRole())

	resources := NewResourceAPI(model.Factory())
	resources.Store = makeQueryStore(root, assert, testUser)

	ctx := context.WithValue(context.Background(), ResourcesCtxKey, resources)

	req, err = http.NewRequestWithContext(ctx, "POST", "/register", nil)
	req.Body = data.Body()

	assert.Nil(err)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusCreated, rr.Code)

	// Re-register should fail
	req, err = http.NewRequestWithContext(ctx, "POST", "/register", nil)
	req.Body = data.Body()

	assert.Nil(err)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusForbidden, rr.Code)

	testForward(assert, h)
}

func TestNoKeyUser(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "test_no_key_user"

	defer func() { os.RemoveAll(root) }()

	key, _ := auth.Generate()
	testUser := user.NewUser("testuser", "test@cisco.com", "bigword", key.Public(), DefaultRole())

	resources := NewResourceAPI(model.Factory())
	resources.Store = makeQueryStore(root, assert, testUser)

	ctx := context.WithValue(context.Background(), ResourcesCtxKey, resources)

	req, err := http.NewRequestWithContext(ctx, "POST", "/register", nil)
	assert.Nil(err)
	assert.NotNil(req)

	data := newRData("testeruser", "secrect", "tester@cisco.com", false)
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

	root := "test_same_user"

	defer func() { os.RemoveAll(root) }()

	data := newRData("testuser", "secrect", "test@cisco.com", true)
	testUser := user.NewUser("testuser", "test@cisco.com", "bigword", data.Key, DefaultRole())

	resources := NewResourceAPI(model.Factory())
	resources.Store = makeQueryStore(root, assert, testUser)

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

	root := "test_bad_request"

	defer func() { os.RemoveAll(root) }()

	key, _ := auth.Generate()
	registryData := fmt.Sprintf("{\"user\":\"%s\"\"password\":\"%s\"}", "Bad", "Request")

	v, err := json.Marshal(registryData)
	assert.Nil(err)
	assert.NotEmpty(v)

	testUser := user.NewUser("testuser", "test@cisco.com", "bigword", key.Public(), DefaultRole())
	resources := NewResourceAPI(model.Factory())
	resources.Store = makeQueryStore(root, assert, testUser)

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

	root := "test_delete_user"

	defer func() { os.RemoveAll(root) }()

	key, _ := auth.Generate()
	testUser := user.NewUser("testuser", "test@cisco.com", "bigword", key.Public(), DefaultRole())
	testUser2 := user.NewUser("testuser", "test@cisco.com", "bigword", key.Public(), DefaultRole())

	store := makeQueryStore(root, assert, testUser)
	assert.NotNil(store)
	err := deleteUser(testUser2, store)
	assert.NotNil(err)
}
