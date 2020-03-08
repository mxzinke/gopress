package gopress

import "net/http"

// A Route holds private information for a route.
type Route struct {
	pattern string
	methods RouteMethods
}

// RouteMethods defines HTTP Handler functions for each REST method.
// You don't have to define for each method a function, just pick what you need.
type RouteMethods struct {
	Head    Handler
	Get     Handler
	Post    Handler
	Put     Handler
	Patch   Handler
	Delete  Handler
	Trace   Handler
	Options Handler
}

// AddRoute creates a new route for a pattern which leads
// to the correct gopress.Handler by the request method.
func (router *Router) AddRoute(pattern string, methods RouteMethods) {
	route := new(Route)
	route.pattern = pattern
	route.methods = methods

	router.routes = append(router.routes, route)

	router.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodHead && methods.Head != nil {
			MakeHandlerExecuteable(methods.Head)(w, r)
		} else if r.Method == http.MethodGet && methods.Get != nil {
			MakeHandlerExecuteable(methods.Get)(w, r)
		} else if r.Method == http.MethodPost && methods.Post != nil {
			MakeHandlerExecuteable(methods.Post)(w, r)
		} else if r.Method == http.MethodPut && methods.Put != nil {
			MakeHandlerExecuteable(methods.Put)(w, r)
		} else if r.Method == http.MethodPatch && methods.Patch != nil {
			MakeHandlerExecuteable(methods.Patch)(w, r)
		} else if r.Method == http.MethodDelete && methods.Delete != nil {
			MakeHandlerExecuteable(methods.Delete)(w, r)
		} else if r.Method == http.MethodTrace && methods.Trace != nil {
			MakeHandlerExecuteable(methods.Trace)(w, r)
		} else if r.Method == http.MethodOptions && methods.Options != nil {
			MakeHandlerExecuteable(methods.Options)(w, r)
		} else {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	})
}
