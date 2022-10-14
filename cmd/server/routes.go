package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// routeHandler returns a http handler that handles all routes under the
// /api/v1 endpoint. It is expected that this handler is the final handler
// and requires the request context to be set with log, store, auth etc.
func routeHandler() http.Handler {
	router := httprouter.New()
	router.GET("/api/v1/types", handleTypes())
	router.GET("/api/v1/labels", handleLabels())
	router.GET("/api/v1/resources", handleQuery())
	router.POST("/api/v1/resources", handlePost())
	// add the resources from the database on to the zebra server.
	// using put to avoid panic errors in the handle/path for post.
	router.PUT("/api/v1/resources", postToZebra())
	router.DELETE("/api/v1/resources", handleDelete())

	return router
}
