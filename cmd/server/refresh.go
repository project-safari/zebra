package main

import (
	"context"
	"net/http"

	"github.com/go-logr/logr"
	"github.com/julienschmidt/httprouter"
	"github.com/project-safari/zebra/auth"
	"github.com/project-safari/zebra/query"
)

func handleRefresh(ctx context.Context, store *query.QueryStore, authKey string) httprouter.Handle {
	return func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
		log := logr.FromContextOrDiscard(ctx)

		jwtCookie, err := req.Cookie("jwt")
		if err != nil {
			log.Error(err, "jwt cookie missing")
			res.WriteHeader(http.StatusBadRequest)

			return
		}

		// Parse the claims
		jwtClaims, err := auth.FromJWT(jwtCookie.Value, authKey)
		if err != nil {
			log.Error(err, "bad jwt token")
			res.WriteHeader(http.StatusUnauthorized)

			return
		}

		// Make sure the jwt is still valid
		if err := jwtClaims.Valid(); err != nil {
			log.Error(err, "invalid jwt token")
			res.WriteHeader(http.StatusUnauthorized)

			return
		}

		// Make sure the user still exists
		user := findUser(store, jwtClaims.Subject)
		if user == nil {
			log.Error(err, "user not found", "user", jwtClaims.Subject)
			res.WriteHeader(http.StatusUnauthorized)

			return
		}

		// Create a new token and cookie
		claims := auth.NewClaims("zebra", user.Name, user.Role)
		respondWithClaims(log, res, claims, authKey)

		log.Info("refresh succeeded", "user", jwtClaims.Subject)
	}
}
