package main

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

var ErrEmptyBody = errors.New("request body is empty")

func readJSON(ctx context.Context, req *http.Request, data interface{}) error {
	err := OpenLogFile("./zebralog.log")
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}

	log.Println("request", "body", string(body))

	if len(body) > 0 {
		err = json.Unmarshal(body, data)
	} else {
		err = ErrEmptyBody
	}

	return err
}

func writeJSON(ctx context.Context, res http.ResponseWriter, data interface{}) {
	err := OpenLogFile("./zebralog.log")
	if err != nil {
		log.Fatal(err)
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)

	if _, err := res.Write(bytes); err != nil {
		log.Println(err, "error writing response")
	}
}
