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
	assert.Nil(jini.Validate(ctx))
	assert.Nil(jini.AuthenticatePassword(jiniWords))

	return jini
}

func makeQueryStore(root string, assert *assert.Assertions, user *auth.User) zebra.Store {
	factory := store.DefaultFactory()
	assert.NotNil(factory)

	store := store.NewResourceStore(root, factory)
	if store.Initialize() != nil || store.Create(user) != nil {
		return nil
	}

	return store
}

func makeRequest(assert *assert.Assertions, user string, password string, email string) *http.Request {
	req, err := http.NewRequest("POST", "/login", nil)
	assert.Nil(err)
	assert.NotNil(req)

	v := fmt.Sprintf("{\"name\":\"%s\",\"password\":\"%s\",\"email\":\"%s\"}", user, password, email)
	req.Body = ioutil.NopCloser(bytes.NewBufferString(v))

	return req
}

func TestFindUser(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "teststore1"

	t.Cleanup(func() { os.RemoveAll(root) })

	store := makeQueryStore(root, assert, makeUser(assert))

	assert.Nil(findUser(store, ""))
	assert.NotNil(findUser(store, "email@domain"))
}

func TestBadUser(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "teststore2"

	t.Cleanup(func() { os.RemoveAll(root) })

	req := makeRequest(assert, "ali", "aliOfAgrabha", "email@domain1")
	store := makeQueryStore(root, assert, makeUser(assert))

	h := handleLogin(context.Background(), store, authKey)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h(w, r, nil)
	})

	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusUnauthorized, rr.Code)

	req = makeRequest(assert, "\"", "\"", "\"")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusBadRequest, rr.Code)
}

func TestBadLogin(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "teststore3"

	t.Cleanup(func() { os.RemoveAll(root) })

	req := makeRequest(assert, "jini", "aliOfAgrabha", "email@domain2")
	store := makeQueryStore(root, assert, makeUser(assert))
	h := handleLogin(context.Background(), store, authKey)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h(w, r, nil)
	})

	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusUnauthorized, rr.Code)
}

func TestLogin(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "teststore4"

	t.Cleanup(func() { os.RemoveAll(root) })

	req := makeRequest(assert, "jini", jiniWords, "email@domain")
	store := makeQueryStore(root, assert, makeUser(assert))
	h := handleLogin(setupLogger(nil), store, authKey)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h(w, r, nil)
	})

	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusOK, rr.Code)

	jwt := &struct {
		JWT string `json:"jwt"`
	}{}
	assert.Nil(json.Unmarshal(rr.Body.Bytes(), jwt))
	assert.NotNil(jwt)
}
