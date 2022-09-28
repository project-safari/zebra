package main

import (
	"context"
	"encoding/base64"
	"net/http"

	"github.com/go-logr/logr"
	"github.com/project-safari/zebra/auth"
	"gojini.dev/web"
)

func authAdapter() web.Adapter {
	return func(nextHandler http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			if nextReq := rsaKey(res, req); nextReq != nil {
				callNext(nextHandler, res, req)
			} else if nextReq := jwtClaims(res, req); nextReq != nil {
				callNext(nextHandler, res, req)
			} else {
				// No auth token so return unautorized status
				res.WriteHeader(http.StatusUnauthorized)
			}
		})
	}
}

func creds(r *http.Request) (string, string) {
	user := r.Header.Get("Zebra-Auth-User")
	token := r.Header.Get("Zebra-Auth-Token")
	bt, _ := base64.StdEncoding.DecodeString(token)

	return user, string(bt)
}

func rsaKey(res http.ResponseWriter, req *http.Request) *http.Request {
	ctx := req.Context()
	log := logr.FromContextOrDiscard(ctx)
	api, ok := ctx.Value(ResourcesCtxKey).(*ResourceAPI)

	if !ok {
		log.Error(nil, "resources not in context")

		return nil
	}

	userEmail, userToken := creds(req)
	if userEmail == "" || userToken == "" {
		// No valid auth token
		return nil
	}

	// Make sure the user still exists
	user := findUser(api.Store, userEmail)
	if user == nil {
		log.Error(nil, "user not found", "user", userEmail)
		res.WriteHeader(http.StatusUnauthorized)

		return nil
	}

	// Verify that token is valid
	if e := user.Authenticate(userToken); e != nil {
		log.Error(e, "user token invalid")
		res.WriteHeader(http.StatusUnauthorized)

		return nil
	}

	// Set the claims into request
	claims := auth.NewClaims("zebra", user.Meta.Name, user.Role, user.Email)
	ctx = context.WithValue(ctx, ClaimsCtxKey, claims)

	return req.Clone(ctx)
}

// jwtClaims extracts the jwt claims from the request cookie, validates it
// and sets it into the request context for all subsequest http handlers
// to use. It returns an error if the jwt is not present or invalid.
func jwtClaims(res http.ResponseWriter, req *http.Request) *http.Request {
	ctx := req.Context()
	log := logr.FromContextOrDiscard(ctx)
	api, ok := ctx.Value(ResourcesCtxKey).(*ResourceAPI)

	if !ok {
		log.Error(nil, "resources not in context")

		return nil
	}

	authKey, ok := ctx.Value(AuthCtxKey).(string)
	if !ok {
		log.Error(nil, "authKey not in context")
		res.WriteHeader(http.StatusInternalServerError)

		return nil
	}

	jwtCookie, err := req.Cookie("jwt")
	if err != nil {
		log.Error(err, "jwt cookie missing")

		return nil
	}

	// Parse the claims
	jwtClaims, err := auth.FromJWT(jwtCookie.Value, authKey)
	if err != nil {
		log.Error(err, "bad jwt token")
		res.WriteHeader(http.StatusUnauthorized)

		return nil
	}

	// Make sure the jwt is still valid
	if err := jwtClaims.Valid(); err != nil {
		log.Error(err, "invalid jwt token")
		res.WriteHeader(http.StatusUnauthorized)

		return nil
	}

	// Make sure the user still exists
	user := findUser(api.Store, jwtClaims.Email)
	if user == nil {
		log.Error(err, "user not found", "user", jwtClaims.Subject)
		res.WriteHeader(http.StatusUnauthorized)

		return nil
	}

	// Set the claims into request
	ctx = context.WithValue(ctx, ClaimsCtxKey, jwtClaims)

	return req.Clone(ctx)
}
