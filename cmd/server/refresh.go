package main

import (
	"log"
	"net/http"

	"github.com/project-safari/zebra/auth"
	"gojini.dev/web"
)

func refreshAdapter() web.Adapter {
	return func(nextHandler http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			if req.URL.Path != "/refresh" {
				// This is not a refresh request just forward it
				callNext(nextHandler, res, req)

				return
			}

			ctx := req.Context()
			err := OpenLogFile("./zebralog.log")
			if err != nil {
				log.Fatal(err)
			}
			jwtClaims, ok := ctx.Value(ClaimsCtxKey).(*auth.Claims)

			if !ok {
				log.Println(nil, "claims not in context")
				res.WriteHeader(http.StatusUnauthorized)

				return
			}

			authKey, ok := ctx.Value(AuthCtxKey).(string)
			if !ok {
				log.Println(nil, "authKey not in context")
				res.WriteHeader(http.StatusUnauthorized)

				return
			}

			// Create a new token and cookie
			claims := auth.NewClaims("zebra", jwtClaims.Subject, jwtClaims.Role, jwtClaims.Email)
			respondWithClaims(ctx, res, claims, authKey)

			log.Println("refresh succeeded", "user", jwtClaims.Subject)
		})
	}
}
