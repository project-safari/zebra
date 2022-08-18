/*
file contains functions to read / write json.
*/

package main

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/go-logr/logr"
)

// ErrEmptyBody returns an error if the request body is empty.
var ErrEmptyBody = errors.New("request body is empty")

// Read JSON.
func readJSON(ctx context.Context, req *http.Request, data interface{}) error {
	log := logr.FromContextOrDiscard(ctx)

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}

	log.Info("request", "body", string(body))

	if len(body) > 0 {
		err = json.Unmarshal(body, data)
	} else {
		err = ErrEmptyBody
	}

	return err
}

// Write json.
func writeJSON(ctx context.Context, res http.ResponseWriter, data interface{}) {
	log := logr.FromContextOrDiscard(ctx)

	bytes, err := json.Marshal(data)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)

	if _, err := res.Write(bytes); err != nil {
		log.Error(err, "error writing response")
	}
}
