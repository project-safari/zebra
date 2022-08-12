package main

import (
	"net/http"

	"github.com/go-logr/logr"
	"github.com/project-safari/zebra/auth"
	"gojini.dev/web"
)

// function to refresh the web adapter.
//
// if everything succeeds, it creates a new token and cookie for the refreshed session.
//
// returns web.Adapter.
func refreshAdapter() web.Adapter {
	return func(nextHandler http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			if req.URL.Path != "/refresh" {
				// This is not a refresh request just forward it
				callNext(nextHandler, res, req)

				return
			}

			ctx := req.Context()
			log := logr.FromContextOrDiscard(ctx)
			jwtClaims, ok := ctx.Value(ClaimsCtxKey).(*auth.Claims)
			if !ok {
				log.Error(nil, "claims not in context")
				res.WriteHeader(http.StatusUnauthorized)

				return
			}

			authKey, ok := ctx.Value(AuthCtxKey).(string)
			if !ok {
				log.Error(nil, "authKey not in context")
				res.WriteHeader(http.StatusUnauthorized)

				return
			}

			// Create a new token and cookie
			claims := auth.NewClaims("zebra", jwtClaims.Subject, jwtClaims.Role, jwtClaims.Email)
			respondWithClaims(ctx, res, claims, authKey)

			log.Info("refresh succeeded", "user", jwtClaims.Subject)
		})
	}
}
