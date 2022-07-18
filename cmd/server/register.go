package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-logr/logr"
	"github.com/julienschmidt/httprouter"
)

func handleRegister(ctx context.Context) httprouter.Handle {
	return func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
		log := logr.FromContextOrDiscard(ctx)
		registryData := &struct {
			Name     string `json:"user"`
			Password string `json:"password"`
			Email    string `json:"email"`
			Key      string `json:"key"`
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

	}
}

func responseRegister(log logr.Logger, res http.ResponseWriter) {
	regData := &struct {
		Name   string `json:"user"`
		Email  string `json:"email"`
		Active bool   `json:"active"`
		Role   string `json:"role"`
	}{}

	bytes, err := json.Marshal(regData)
	if err != nil {
		log.Error(err, "crazy we can't marshal our own data!")
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)

	if _, err := res.Write(bytes); err != nil {
		log.Error(err, "error writing response")
	}
}
