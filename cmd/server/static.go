package main

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"gojini.dev/web"
)

//go:embed build
var embedFiles embed.FS

func staticAdaptor() web.Adapter {
	return func(nextHandler http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			server := http.FileServer(getFilesSystem())
			if req.URL.Path == "/" || req.URL.Path == "/ZebraIcon.svg" || strings.HasPrefix(req.URL.Path, "/static") {
				server.ServeHTTP(res, req)

				return
			}
			callNext(nextHandler, res, req)
		})
	}
}

func getFilesSystem() http.FileSystem {
	// Get the build subdirectory as the
	// root directory so that it can be passed
	// to the http.FileServer
	fsys, err := fs.Sub(embedFiles, "build")
	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}
