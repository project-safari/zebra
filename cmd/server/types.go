package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/go-logr/logr"
	"github.com/julienschmidt/httprouter"
	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/store"
)

func handleTypes(ctx context.Context) httprouter.Handle {
	allTypes := store.DefaultFactory()

	return func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
		typeReq := &struct {
			Types []string `json:"types"`
		}{Types: []string{}}

		if err := readReq(ctx, req, typeReq); err != nil {
			res.WriteHeader(http.StatusBadRequest)

			return
		}

		typeRes := &struct {
			Types []zebra.Type `json:"types"`
		}{Types: []zebra.Type{}}

		if len(typeReq.Types) == 0 {
			// return all types
			typeRes.Types = allTypes.Types()
		} else {
			for _, t := range typeReq.Types {
				if aType, ok := allTypes.Type(t); ok {
					typeRes.Types = append(typeRes.Types, aType)
				}
			}
		}

		writeJSON(ctx, res, typeRes)
	}
}

func readReq(ctx context.Context, req *http.Request, data interface{}) error {
	log := logr.FromContextOrDiscard(ctx)

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil
	}

	log.Info("request", "body", string(body))

	if len(body) > 0 {
		err = json.Unmarshal(body, data)
	}

	return err
}
