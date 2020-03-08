package gopress

import (
	"net/http"

	"github.com/jinzhu/copier"
)

// A Handler represents a function which returns a usable net/http HandlerFunc.
// This make
type Handler func(http.HandlerFunc) http.HandlerFunc

// The HandlerFunc is defining how the function for Handler creation should look like
type HandlerFunc func(*Request, Response)

// CreateHandler gives the possibility to create middleware or route handlers
// by passing in a HandlerFunc which contains a Request and a Response.
func CreateHandler(handleFunc HandlerFunc) Handler {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			request := new(Request)
			copier.Copy(request, r)

			request.Path = r.URL.Path

			handleFunc(request, Response{w})

			f(w, r)
		}
	}
}

// Chain combines multiple handlers to on handler
func Chain(handlers ...Handler) Handler {
	return func(f http.HandlerFunc) http.HandlerFunc {
		for _, handler := range reverseHandlers(handlers) {
			f = handler(f)
		}
		return f
	}
}

// MakeHandlerExecuteable converts the Handler into a http.HandlerFunc by passing in
// a start http.HandlerFunc which does nothing. Be aware, that you are leaving
// the package with that method and relay on the "http/net" functionalities.
func MakeHandlerExecuteable(handler Handler) http.HandlerFunc {
	dummyHandlerFunc := func(wf http.ResponseWriter, rf *http.Request) {}
	return handler(dummyHandlerFunc)
}

func reverseHandlers(s []Handler) []Handler {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}

	return s
}
