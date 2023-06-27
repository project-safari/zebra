package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-logr/logr"
	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/auth"
	"github.com/project-safari/zebra/model/user"
	"gojini.dev/web"
)

func registerAdapter() web.Adapter {
	return func(nextHandler http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			if req.URL.Path != "/register" {
				// This is not a register request just forward it
				callNext(nextHandler, res, req)

				return
			}

			registerHandler(res, req)
		})
	}
}

func registerHandler(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	log := logr.FromContextOrDiscard(ctx)
	api, ok := ctx.Value(ResourcesCtxKey).(*ResourceAPI)

	if !ok {
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	regReq := &struct {
		Name     string            `json:"name"`
		Password string            `json:"password"`
		Email    string            `json:"email"`
		Key      *auth.RsaIdentity `json:"key"`
	}{}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)

		return
	}

	log.Info("request", "body", string(body))

	if err := json.Unmarshal(body, regReq); err != nil {
		log.Error(err, "bad body")
		res.WriteHeader(http.StatusBadRequest)

		return
	}

	store := api.Store

	if findUser(store, regReq.Email) != nil {
		log.Error(err, "user already exist", "user", regReq.Name)
		res.WriteHeader(http.StatusForbidden)

		return
	}

	newuser := user.NewUser(regReq.Name, regReq.Email,
		regReq.Password, regReq.Key, DefaultRole())

	if err := store.Create(newuser); err != nil {
		log.Error(err, "user cant be stored", "user", regReq.Name)
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	responseRegister(log, res, newuser)
	log.Info("Registry succeeded", "user", regReq.Name)
}

/*
HOW TO USE THIS:
	// this, noteType := noifications.NotifyRegister()
	// noifications.SendAccountNotification("New Account Created", this, regReq.Email, noteType)
*/

// nolint

func responseRegister(log logr.Logger, res http.ResponseWriter, newuser *user.User) {
	bytes, err := json.Marshal(newuser)
	if err != nil {
		log.Error(err, "crazy we can't marshal our own data!")
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)

	if _, err := res.Write(bytes); err != nil {
		log.Error(err, "error writing response")
	}
}

func DefaultRole() *auth.Role {
	read, _ := auth.NewPriv("", false, true, false, false)
	role := &auth.Role{
		Name:       "user",
		Privileges: []*auth.Priv{read},
	}

	return role
}

func deleteUser(u *user.User, store zebra.Store) error {
	if err := store.Delete(u); err != nil {
		return err
	}

	return nil
}

func changePassword(u *user.User, newpass string) {
	u.PasswordHash = user.HashPassword(newpass)
}
