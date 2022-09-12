package main //nolint:testpackage

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/project-safari/zebra/auth"
	"github.com/project-safari/zebra/model"
	"github.com/stretchr/testify/assert"
)

func TestAuthAdapter(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	a := authAdapter()
	assert.NotNil(a)

	handler := a(nil)
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/api/v1/resources", nil)
	assert.NotNil(req)
	assert.Nil(err)

	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusUnauthorized, rr.Code)

	root := "test_auth_adapter"

	defer func() { os.RemoveAll(root) }()

	resources := NewResourceAPI(model.Factory())
	resources.Store = makeQueryStore(root, assert, makeUser(assert))

	ctx := context.WithValue(context.Background(), ResourcesCtxKey, resources)
	ctx = context.WithValue(ctx, AuthCtxKey, authKey)
	rr = httptest.NewRecorder()
	req, err = http.NewRequestWithContext(ctx, "GET", "/", nil)
	assert.Nil(err)

	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusUnauthorized, rr.Code)
}

func TestRSAKey(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "test_rsa_key"

	defer func() { os.RemoveAll(root) }()

	user := makeUser(assert)
	priKey := user.Key

	user.Key = priKey.Public()

	resources := NewResourceAPI(model.Factory())
	resources.Store = makeQueryStore(root, assert, user)

	ctx := context.WithValue(context.Background(), ResourcesCtxKey, resources)
	ctx = context.WithValue(ctx, AuthCtxKey, authKey)

	tb, err := priKey.Sign([]byte(user.Email))
	assert.Nil(err)
	assert.NotEmpty(tb)

	token := base64.StdEncoding.EncodeToString(tb)
	assert.NotEmpty(token)

	req, err := http.NewRequestWithContext(ctx, "GET", "/", nil)
	assert.Nil(err)

	req.Header.Set("Zebra-Auth-User", user.Email)
	req.Header.Set("Zebra-Auth-Token", token)

	rr := httptest.NewRecorder()

	a := authAdapter()
	handler := a(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
	}))
	handler.ServeHTTP(rr, req)

	assert.Equal(http.StatusOK, rr.Code)
}

func TestJWT(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "test_jwt"

	defer func() { os.RemoveAll(root) }()

	user := makeUser(assert)
	priKey := user.Key

	user.Key = priKey.Public()

	resources := NewResourceAPI(model.Factory())
	resources.Store = makeQueryStore(root, assert, user)

	ctx := context.WithValue(context.Background(), ResourcesCtxKey, resources)
	ctx = context.WithValue(ctx, AuthCtxKey, authKey)

	req, err := http.NewRequestWithContext(ctx, "GET", "/", nil)
	assert.Nil(err)

	claims := auth.NewClaims("zebra", user.Meta.Name, user.Role, user.Email)
	jwt := claims.JWT(authKey)
	req.AddCookie(makeCookie(jwt))

	rr := httptest.NewRecorder()

	a := authAdapter()
	handler := a(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
	}))
	handler.ServeHTTP(rr, req)

	assert.Equal(http.StatusOK, rr.Code)
}

func TestBadRSAKey(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "test_bad_rsa_key"

	defer func() { os.RemoveAll(root) }()

	user := makeUser(assert)
	priKey := user.Key

	user.Key = priKey.Public()

	resources := NewResourceAPI(model.Factory())
	resources.Store = makeQueryStore(root, assert, user)

	ctx := context.WithValue(context.Background(), ResourcesCtxKey, resources)
	ctx = context.WithValue(ctx, AuthCtxKey, authKey)

	tb, err := priKey.Sign([]byte(user.Email))
	assert.Nil(err)
	assert.NotEmpty(tb)

	token := base64.StdEncoding.EncodeToString(tb)
	assert.NotEmpty(token)

	// Bad user
	req, err := http.NewRequestWithContext(ctx, "GET", "/", nil)
	assert.Nil(err)

	req.Header.Set("Zebra-Auth-User", "doesnotexist")
	req.Header.Set("Zebra-Auth-Token", token)

	rr := httptest.NewRecorder()

	a := authAdapter()
	handler := a(nil)
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusUnauthorized, rr.Code)

	// Bad token
	req, err = http.NewRequestWithContext(ctx, "GET", "/", nil)
	assert.Nil(err)

	req.Header.Set("Zebra-Auth-User", user.Email)
	req.Header.Set("Zebra-Auth-Token", "badtoken")

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusUnauthorized, rr.Code)
}

func TestBadJWT(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "test_bad_jwt"

	defer func() { os.RemoveAll(root) }()

	user := makeUser(assert)
	priKey := user.Key

	user.Key = priKey.Public()

	resources := NewResourceAPI(model.Factory())
	resources.Store = makeQueryStore(root, assert, user)

	ctx := context.WithValue(context.Background(), ResourcesCtxKey, resources)

	// No authKey
	req, err := http.NewRequestWithContext(ctx, "GET", "/", nil)
	assert.Nil(err)

	rr := httptest.NewRecorder()

	a := authAdapter()
	handler := a(nil)
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusInternalServerError, rr.Code)

	// Bad jwt
	ctx = context.WithValue(ctx, AuthCtxKey, authKey)
	req, err = http.NewRequestWithContext(ctx, "GET", "/", nil)
	assert.Nil(err)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusUnauthorized, rr.Code)

	// bad token
	req, err = http.NewRequestWithContext(ctx, "GET", "/", nil)
	req.AddCookie(makeCookie("badjwt"))
	assert.Nil(err)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusUnauthorized, rr.Code)
}
