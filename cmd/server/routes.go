package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/project-safari/zebra/script"
)

// routeHandler returns a http handler that handles all routes under the
// /api/v1 endpoint. It is expected that this handler is the final handler
// and requires the request context to be set with log, store, auth etc.
func routeHandler() http.Handler {
	router := httprouter.New()
	router.GET("/api/v1/types", handleTypes())
	router.GET("/api/v1/labels", handleLabels())
	router.GET("/api/v1/resources", handleQuery())
	router.POST("/api/v1/resources", script.HandlePost())
	router.DELETE("/api/v1/resources", handleDelete())

	return router
}
