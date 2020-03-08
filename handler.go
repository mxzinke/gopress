package gopress

import "net/http"

// A Handler represents a function which returns a usable net/http HandlerFunc.
// This make
type Handler func(http.HandlerFunc) http.HandlerFunc

// The HandlerFunc is defining how the function for Handler creation should look like
type HandlerFunc func(*Request, Response)

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
