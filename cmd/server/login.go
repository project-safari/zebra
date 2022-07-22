package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/go-logr/logr"
	"github.com/julienschmidt/httprouter"
	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/auth"
)

func handleLogin(ctx context.Context, store zebra.Store, authKey string) httprouter.Handle {
	return func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
		log := logr.FromContextOrDiscard(ctx)
		userData := &struct {
			Name     string `json:"name"`
			Password string `json:"password"`
			Email    string `json:"email"`
		}{}

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			res.WriteHeader(http.StatusBadRequest)

			return
		}

		log.Info("request", "body", string(body))

		if err := json.Unmarshal(body, userData); err != nil {
			res.WriteHeader(http.StatusBadRequest)

			return
		}

		user := findUser(store, userData.Email)
		if user == nil {
			log.Error(err, "user not found", "user", userData.Name)
			res.WriteHeader(http.StatusUnauthorized)

			return
		}

		if err := user.AuthenticatePassword(userData.Password); err != nil {
			log.Error(err, "user auth failed", "user", userData.Name)
			res.WriteHeader(http.StatusUnauthorized)

			return
		}

		claims := auth.NewClaims("zebra", user.Name, user.Role, user.Email)
		respondWithClaims(log, res, claims, authKey)

		log.Info("login succeeded", "user", userData.Name)
	}
}

func makeCookie(jwt string) *http.Cookie {
	cookie := new(http.Cookie)
	cookie.Name = "jwt"
	cookie.Value = jwt
	cookie.Expires = time.Now().Add(auth.TokenDuration)

	return cookie
}

func findUser(store zebra.Store, email string) *auth.User {
	resMap := store.QueryType([]string{"User"})
	users := resMap.Resources["User"]

	for _, u := range users.Resources {
		user, ok := u.(*auth.User)
		if ok && user.Email == email {
			return user
		}
	}

	return nil
}

func respondWithClaims(log logr.Logger, res http.ResponseWriter,
	claims *auth.Claims, authKey string,
) {
	resData := &struct {
		JWT string `json:"jwt"`
	}{JWT: claims.JWT(authKey)}

	bytes, err := json.Marshal(resData)
	if err != nil {
		log.Error(err, "crazy we can't marshal our own data!")
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	http.SetCookie(res, makeCookie(resData.JWT))
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)

	if _, err := res.Write(bytes); err != nil {
		log.Error(err, "error writing response")
	}
}
