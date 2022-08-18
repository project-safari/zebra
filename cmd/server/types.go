package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/store"
)

func handleTypes() httprouter.Handle {
	allTypes := store.DefaultFactory()

	return func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
		ctx := req.Context()

		typeReq := &struct {
			Types []string `json:"types"`
		}{Types: []string{}}

		if err := readJSON(ctx, req, typeReq); err != nil {
			res.WriteHeader(http.StatusBadRequest)

			return
		}

		typeRes := &struct {
			Types []zebra.Type `json:"types"`
		}{Types: []zebra.Type{}}

		if len(typeReq.Types) == 0 {
			// Return all types.
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
