package main

import (
	"context"
	"net/http"
	"time"

	"github.com/go-logr/logr"
	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/auth"
	"github.com/project-safari/zebra/model/noifications"
	"github.com/project-safari/zebra/model/user"
	"gojini.dev/web"
)

func loginAdapter() web.Adapter {
	return func(nextHandler http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			if req.URL.Path != "/login" {
				// This is not a login request just forward it
				callNext(nextHandler, res, req)

				return
			}

			ctx := req.Context()
			log := logr.FromContextOrDiscard(ctx)
			api, ok := ctx.Value(ResourcesCtxKey).(*ResourceAPI)
			if !ok {
				res.WriteHeader(http.StatusInternalServerError)

				return
			}

			authKey, ok := ctx.Value(AuthCtxKey).(string)
			if !ok {
				res.WriteHeader(http.StatusInternalServerError)

				return
			}

			userData := &struct {
				Password string `json:"password"`
				Email    string `json:"email"`
			}{}

			if err := readJSON(ctx, req, userData); err != nil {
				res.WriteHeader(http.StatusBadRequest)

				return
			}

			user := findUser(api.Store, userData.Email)
			if user == nil {
				log.Error(nil, "user not found", "user", userData.Email)
				res.WriteHeader(http.StatusUnauthorized)

				return
			}

			if err := user.AuthenticatePassword(userData.Password); err != nil {
				log.Error(err, "user auth failed", "user", user.Email)
				res.WriteHeader(http.StatusUnauthorized)

				return
			}

			claims := auth.NewClaims("zebra", user.Meta.Name, user.Role, user.Email)
			respondWithClaims(ctx, res, claims, authKey)

			log.Info("login succeeded", "user", user.Email)

			this := noifications.NotifyLogin()
			noifications.SendAccountNotification("New Login into Your Account", this, user.Email)
		})
	}
}

func makeCookie(jwt string) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = "jwt"
	cookie.Value = jwt
	cookie.Expires = time.Now().Add(auth.TokenDuration)

	return cookie
}

func findUser(store zebra.Store, email string) *user.User {
	resMap := store.QueryType([]string{"system.user"})

	users := resMap.Resources["system.user"]
	if users == nil {
		return nil
	}

	for _, u := range users.Resources {
		user, ok := u.(*user.User)
		if ok && user.Email == email {
			return user
		}
	}

	return nil
}

func respondWithClaims(ctx context.Context, res http.ResponseWriter,
	claims *auth.Claims, authKey string,
) {
	resData := &struct {
		JWT   string     `json:"jwt"`
		Role  *auth.Role `json:"role"`
		Email string     `json:"email"`
		Name  string     `json:"name"`
	}{JWT: claims.JWT(authKey), Role: claims.Role, Email: claims.Email, Name: claims.Id}

	http.SetCookie(res, makeCookie(resData.JWT))
	writeJSON(ctx, res, resData)
}
