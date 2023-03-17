package main //nolint:testpackage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/auth"
	"github.com/project-safari/zebra/model"
	"github.com/project-safari/zebra/model/user"
	"github.com/project-safari/zebra/store"
	"github.com/stretchr/testify/assert"
)

// Mock const values for authentication to be used in tests.
const (
	authKey   = "abracadabra"
	jiniWords = "youGetThreewishes!"
)

// Test function to find a user.
func TestFindUser(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "test_find_user"

	defer func() { os.RemoveAll(root) }()

	store := makeQueryStore(root, assert, nil)
	assert.Nil(findUser(store, "email@domain"))

	assert.Nil(os.RemoveAll(root))

	store = makeQueryStore(root, assert, makeUser(assert))

	assert.Nil(findUser(store, ""))
	assert.NotNil(findUser(store, "email@domain"))
}

// Mock function to create a user to be used in tests.
// It returns a pointer to user.User.
func makeUser(assert *assert.Assertions) *user.User {
	ctx := context.Background()
	all, e := auth.NewPriv("", true, true, true, true)
	assert.Nil(e)
	assert.NotNil(all)

	jiniKey, err := auth.Generate()
	assert.Nil(err)
	assert.NotNil(jiniKey)

	jini := user.NewUser("jini", "email@domain", jiniWords, jiniKey,
		&auth.Role{Name: "admin", Privileges: []*auth.Priv{all}})

	jini.Key = jiniKey
	assert.Nil(jini.Validate(ctx))
	assert.Nil(jini.AuthenticatePassword(jiniWords))

	return jini
}

// Mock function to make a query store to be used in tests.
// It returns a zebra.Store.
func makeQueryStore(root string, assert *assert.Assertions, user *user.User) zebra.Store {
	factory := model.Factory()
	assert.NotNil(factory)

	store := store.NewResourceStore(root, factory)
	err := store.Initialize()
	assert.Nil(err)

	if user != nil {
		err = store.Create(user)
		assert.Nil(err)
	}

	return store
}

// Mock function to make a query store to be used in tests.
// It returns a zebra.Store.
func makeLoginRequest(assert *assert.Assertions, user string, password string,
	email string, resources *ResourceAPI,
) *http.Request {
	ctx := context.WithValue(context.Background(), ResourcesCtxKey, resources)
	ctx = context.WithValue(ctx, AuthCtxKey, authKey)
	req, err := http.NewRequestWithContext(ctx, "POST", "/login", nil)
	assert.Nil(err)
	assert.NotNil(req)

	v := fmt.Sprintf("{\"name\":\"%s\",\"password\":\"%s\",\"email\":\"%s\"}", user, password, email)
	req.Body = ioutil.NopCloser(bytes.NewBufferString(v))

	return req
}

// Test function for a bad / incorrect user.
func TestBadUser(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "test_bad_user_login"

	defer func() { os.RemoveAll(root) }()

	resources := NewResourceAPI(model.Factory())
	resources.Store = makeQueryStore(root, assert, makeUser(assert))

	req := makeLoginRequest(assert, "ali", "aliOfAgrabha", "email@domain1", resources)

	h := loginAdapter()
	rr := httptest.NewRecorder()
	handler := h(nil)

	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusUnauthorized, rr.Code)

	req = makeLoginRequest(assert, "\"", "\"", "\"", resources)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusBadRequest, rr.Code)
}

// Test function for a bad / incorrect login.
func TestBadLogin(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "test_bad_login"

	defer func() { os.RemoveAll(root) }()

	resources := NewResourceAPI(model.Factory())
	resources.Store = makeQueryStore(root, assert, makeUser(assert))

	req := makeLoginRequest(assert, "jini", "aliOfAgrabha", "email@domain", resources)
	h := loginAdapter()
	rr := httptest.NewRecorder()
	handler := h(nil)

	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusUnauthorized, rr.Code)
}

// Test function for login.
func TestLogin(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "test_login"

	defer func() { os.RemoveAll(root) }()

	resources := NewResourceAPI(model.Factory())
	resources.Store = makeQueryStore(root, assert, makeUser(assert))

	req := makeLoginRequest(assert, "jini", jiniWords, "email@domain", resources)

	h := loginAdapter()
	rr := httptest.NewRecorder()
	handler := h(nil)

	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusOK, rr.Code)

	jwt := &struct {
		JWT string `json:"jwt"`
	}{}
	assert.Nil(json.Unmarshal(rr.Body.Bytes(), jwt))
	assert.NotNil(jwt)
}

// Test function for testForward.
func TestLoginForward(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	testForward(assert, loginAdapter())
}

// Test function for login w/ context.
func TestLoginContext(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	req, err := http.NewRequest("POST", "/login", nil)
	assert.Nil(err)

	h := loginAdapter()
	rr := httptest.NewRecorder()
	handler := h(nil)

	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusInternalServerError, rr.Code)

	ctx := context.WithValue(context.Background(), ResourcesCtxKey,
		NewResourceAPI(model.Factory()))

	req, err = http.NewRequestWithContext(ctx, "POST", "/login", nil)
	assert.Nil(err)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusInternalServerError, rr.Code)
}
