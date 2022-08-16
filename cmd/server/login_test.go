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
	"github.com/project-safari/zebra/cmd/herd/pkg"
	"github.com/project-safari/zebra/store"
	"github.com/stretchr/testify/assert"
)

const (
	authKey   = "abracadabra"
	jiniWords = "youGetThreewishes!"
)

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

func makeUser(assert *assert.Assertions) *auth.User {
	ctx := context.Background()
	all, e := auth.NewPriv("", true, true, true, true)
	assert.Nil(e)
	assert.NotNil(all)

	jini := new(auth.User)
	jini.Name = "jini"
	jini.Email = "email@domain"
	jini.ID = "007"
	jini.Type = auth.UserType()
	jini.PasswordHash = auth.HashPassword(jiniWords)
	jini.Role = &auth.Role{Name: "admin", Privileges: []*auth.Priv{all}}
	jini.Labels = pkg.CreateLabels()
	jini.Labels = pkg.GroupLabels(jini.Labels, "sampleGroup")

	jiniKey, err := auth.Generate()
	assert.Nil(err)
	assert.NotNil(jiniKey)

	jini.Key = jiniKey
	assert.Nil(jini.Validate(ctx, jini.Type.Name))
	assert.Nil(jini.AuthenticatePassword(jiniWords))

	return jini
}

func makeQueryStore(root string, assert *assert.Assertions, user *auth.User) zebra.Store {
	factory := store.DefaultFactory()
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

func TestBadUser(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "test_bad_user_login"

	defer func() { os.RemoveAll(root) }()

	resources := NewResourceAPI(store.DefaultFactory())
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

func TestBadLogin(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "test_bad_login"

	defer func() { os.RemoveAll(root) }()

	resources := NewResourceAPI(store.DefaultFactory())
	resources.Store = makeQueryStore(root, assert, makeUser(assert))

	req := makeLoginRequest(assert, "jini", "aliOfAgrabha", "email@domain", resources)
	h := loginAdapter()
	rr := httptest.NewRecorder()
	handler := h(nil)

	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusUnauthorized, rr.Code)
}

func TestLogin(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "test_login"

	defer func() { os.RemoveAll(root) }()

	resources := NewResourceAPI(store.DefaultFactory())
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

func TestLoginForward(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	testForward(assert, loginAdapter())
}

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
		NewResourceAPI(store.DefaultFactory()))

	req, err = http.NewRequestWithContext(ctx, "POST", "/login", nil)
	assert.Nil(err)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusInternalServerError, rr.Code)
}
