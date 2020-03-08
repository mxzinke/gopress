package gopress

import "net/http"

// A Request represents an HTTP request received by a server or to be sent by a client.
type Request struct {
	http.Request
	Path string
}
