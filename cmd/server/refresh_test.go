package main //nolint:testpackage

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/project-safari/zebra/auth"
	"github.com/stretchr/testify/assert"
)

func TestRefresh(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "teststore5"

	t.Cleanup(func() { os.RemoveAll(root) })

	jini := makeUser(assert)
	claims := auth.NewClaims("zebra", jini.Name, jini.Role)
	jwtStr := claims.JWT(authKey)

	req := makeRequest(assert, "jini", jiniWords)
	req.AddCookie(makeCookie(jwtStr))

	store := makeQueryStore(root, assert, jini)
	h := handleRefresh(setupLogger(nil), store, authKey)
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
