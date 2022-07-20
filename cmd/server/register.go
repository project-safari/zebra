package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-logr/logr"
	"github.com/julienschmidt/httprouter"
	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/auth"
)

type NewUser struct {
	auth.User
	Email  string `json:"email"`
	Active bool   `json:"active"`
}

func handleRegister(ctx context.Context, store zebra.Store) httprouter.Handle {
	return func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
		log := logr.FromContextOrDiscard(ctx)
		registryData := &struct {
			Name     string            `json:"user"`
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

		if err := json.Unmarshal(body, registryData); err != nil {
			res.WriteHeader(http.StatusBadRequest)

			return
		}
		user := findUser(store, registryData.Name)
		if user != nil {
			log.Error(err, "user already exist", "user", registryData.Name)
			res.WriteHeader(http.StatusForbidden)

			return
		}

		newuser := createNewUser(registryData.Email, registryData.Name, registryData.Key, registryData.Password)

		if store.Initialize() != nil || store.Create(newuser) != nil {
			log.Error(err, "user cant be stored", "user", registryData.Name)
			return
		}

		responseRegister(log, res, newuser)
		log.Info("Registry succeeded", "user", registryData.Name)
	}
}

func responseRegister(log logr.Logger, res http.ResponseWriter, newuser *NewUser) {
	regData := &struct {
		Name   string     `json:"user"`
		Email  string     `json:"email"`
		Active bool       `json:"active"`
		Role   *auth.Role `json:"role"`
	}{
		Name:   newuser.Name,
		Email:  newuser.Email,
		Active: false,
		Role:   newuser.Role,
	}

	bytes, err := json.Marshal(regData)
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

func createNewUser(email string, name string, key *auth.RsaIdentity, passward string) *NewUser {
	newuser := &NewUser{
		Email:  email,
		Active: false,
	}
	newuser.Key = key
	newuser.PasswordHash = auth.HashPassword(passward)
	newuser.Name = name
	none, _ := auth.NewPriv("", false, false, false, false)
	newuser.Role = &auth.Role{
		Name:       "No role",
		Privileges: []*auth.Priv{none},
	}
	newuser.BaseResource = *zebra.NewBaseResource("User", nil)

	return newuser
}

func (nw *NewUser) createRole(role *auth.Role) {
	nw.Role = role
}

func (nw *NewUser) changeActiveStatus(store zebra.Store, status bool) {
	if !status {
		nw.deleteUser(store)
	}
	nw.Active = status
}

func (nw *NewUser) deleteUser(store zebra.Store) error {
	err := store.Delete(nw)
	if err != nil {
		return err
	}
	return nil
}

func (nw *NewUser) changePassword(newpass string) {
	nw.PasswordHash = auth.HashPassword(newpass)
}

func (nw *NewUser) changePubKey() error {
	rsa, err := auth.Generate()
	if err != nil {
		return err
	}
	nw.Key = rsa
	return nil
}
