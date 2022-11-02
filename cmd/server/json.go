package main

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/go-logr/logr"
)

// ErrEmptyBody is an error that occurs if the request body is empty.
var ErrEmptyBody = errors.New("request body is empty")

// Function that reads from a request body, logs it, and unmarshals the data.
//
// It takes in a context.Context, a pointer to http.Request, and an interface.
// It returns an error or nil in the absence thereof.
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

// Function that writes the response header.
//
// It takes in a context.Context, a pointer to http.ResponceWriter, and an interface.
// It returns an error or nil in the absence thereof.
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
