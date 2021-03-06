package gopress

import (
	"errors"
	"fmt"
	"net/http"
)

// A Router holds private information of routes and settings of the webserver.
type Router struct {
	routes     []*Route
	settings   RouterSettings
	middleware Handler
	mux        *http.ServeMux
}

// The RouterSettings defining paths and behaviors of the new router object.
type RouterSettings struct {
	// TemplatesPath pointing to a relative or absolute path where the .template files are existing
	TemplatesPath string
	UseSSL        bool
	SSLCertFile   string
	SSLKeyFile    string
}

// NewRouter creates a new Router object by passing the settings
func NewRouter(settings RouterSettings, middlewares ...Handler) (*Router, error) {
	if settings.UseSSL && (settings.SSLCertFile == "" || settings.SSLKeyFile == "") {
		return nil, errors.New(fmt.Sprint("When using SSL/TLS, you have to pass a cert and key file path!"))
	}

	router := new(Router)
	router.settings = settings
	router.middleware = Chain(middlewares...)
	router.mux = http.NewServeMux()

	return router, nil
}

// Start will start up a webserver, depending on the given router settings and given port.
func (router *Router) Start(port int) error {
	if router.settings.UseSSL {
		return http.ListenAndServeTLS(
			fmt.Sprintf(":%d", port),
			router.settings.SSLCertFile,
			router.settings.SSLKeyFile,
			makeHTTPHandler(router.mux, router.middleware),
		)
	}

	return http.ListenAndServe(fmt.Sprintf(":%d", port), makeHTTPHandler(router.mux, router.middleware))
}

func makeHTTPHandler(handler http.Handler, middleware Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlerFunc := middleware(func(wf http.ResponseWriter, rf *http.Request) {
			handler.ServeHTTP(wf, rf)
		})

		handlerFunc(w, r)
	})
}
