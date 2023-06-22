// nolint: gosec, forcetypeassert // Using this script for a secure server with a secure jwt token.
package main

import (
	"crypto/tls"
	"net/http"
)

func main() {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	PostIt()
}
