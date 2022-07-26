package main //nolint:testpackage

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/project-safari/zebra/auth"
	"github.com/stretchr/testify/assert"
	"gojini.dev/web"
)

func TestRefresh(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "teststore5"

	t.Cleanup(func() { os.RemoveAll(root) })

	jini := makeUser(assert)
	claims := auth.NewClaims("zebra", jini.Name, jini.Role, "email@domain")
	jwtStr := claims.JWT(authKey)

	req := makeRefreshRequest(assert, claims, authKey)
	req.AddCookie(makeCookie(jwtStr))

	h := refreshAdapter()
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

func makeRefreshRequest(assert *assert.Assertions, claims *auth.Claims, authKey string) *http.Request {
	ctx := context.WithValue(context.Background(), ClaimsCtxKey, claims)
	ctx = context.WithValue(ctx, AuthCtxKey, authKey)

	req, err := http.NewRequestWithContext(ctx, "POST", "/refresh", nil)
	assert.Nil(err)
	assert.NotNil(req)

	return req
}

func testForward(assert *assert.Assertions, h web.Adapter) {
	handler := web.Wrap(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
	}), h)

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/test", nil)
	assert.Nil(err)
	assert.NotNil(req)

	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusOK, rr.Code)
}

func TestRefreshForward(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	testForward(assert, refreshAdapter())
}

func TestBadRefreshRequest(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	h := refreshAdapter()
	handler := h(nil)

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("POST", "/refresh", nil)
	assert.Nil(err)
	assert.NotNil(req)

	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusInternalServerError, rr.Code)

	rr = httptest.NewRecorder()
	claims := auth.NewClaims("zebra", "a", nil, "email@domain")
	ctx := context.WithValue(context.Background(), ClaimsCtxKey, claims)
	req, err = http.NewRequestWithContext(ctx, "POST", "/refresh", nil)
	assert.Nil(err)
	assert.NotNil(req)

	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusInternalServerError, rr.Code)
}
