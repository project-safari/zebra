package main //nolint:testpackage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/auth"
	"github.com/project-safari/zebra/query"
	"github.com/stretchr/testify/assert"
)

const (
	authKey   = "abracadabra"
	jiniWords = "youGetThreewishes!"
)

func TestFindUser(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	store := makeQueryStore(assert, makeUser(assert))

	assert.Nil(findUser(store, "ali"))
	assert.NotNil(findUser(store, "jini"))
}

func makeUser(assert *assert.Assertions) *auth.User {
	ctx := context.Background()
	all, e := auth.NewPriv("", true, true, true, true)
	assert.Nil(e)
	assert.NotNil(all)

	jini := new(auth.User)
	jini.Name = "jini"
	jini.ID = "007"
	jini.Type = "User"
	jini.PasswordHash = auth.HashPassword(jiniWords)
	jini.Role = &auth.Role{Name: "admin", Privileges: []*auth.Priv{all}}

	jiniKey, err := auth.Generate()
	assert.Nil(err)
	assert.NotNil(jiniKey)

	jini.Key = jiniKey
	assert.Nil(jini.Validate(ctx))
	assert.Nil(jini.AuthenticatePassword(jiniWords))

	return jini
}

func makeQueryStore(assert *assert.Assertions, user *auth.User) *query.QueryStore {
	factory := initTypes()
	assert.NotNil(factory)

	users := zebra.NewResourceMap(factory)
	users.Add(user, "User")

	return query.NewQueryStore(users)
}

func makeRequest(assert *assert.Assertions, user string, password string) *http.Request {
	req, err := http.NewRequest("POST", "/login", nil)
	assert.Nil(err)
	assert.NotNil(req)

	v := fmt.Sprintf("{\"user\":\"%s\",\"password\":\"%s\"}", user, password)
	req.Body = ioutil.NopCloser(bytes.NewBufferString(v))

	return req
}

func TestBadUser(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	req := makeRequest(assert, "ali", "aliOfAgrabha")
	store := makeQueryStore(assert, makeUser(assert))

	h := handleLogin(context.Background(), store, authKey)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h(w, r, nil)
	})

	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusUnauthorized, rr.Code)

	req = makeRequest(assert, "\"", "\"")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusBadRequest, rr.Code)
}

func TestBadLogin(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	req := makeRequest(assert, "jini", "aliOfAgrabha")
	store := makeQueryStore(assert, makeUser(assert))
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

	req := makeRequest(assert, "jini", jiniWords)
	store := makeQueryStore(assert, makeUser(assert))
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
