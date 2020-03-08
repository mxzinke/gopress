package main

import (
	"log"
	"time"

	"github.com/mxzinke/gopress"
)

// User as a test object
type User struct {
	Username   string
	Repository string
}

func main() {
	// Passed Middleware to the router will be always executed (at every request)
	router, err := gopress.NewRouter(gopress.RouterSettings{
		TemplatesPath: "./public/templates",
		UseSSL:        false, // Or make use of ENV with something like os.GetEnv("USE_SSL")
		SSLCertFile:   "./path/to/cert/file",
		SSLKeyFile:    "./path/to/key/file",
	}, StaticLoggingMiddleware)

	if err != nil {
		log.Panic(err)
	}

	router.AddRoute("/", gopress.RouteMethods{
		Get: GetRouteHandler(true),
		// Get: gopress.Chain(AuthMiddleware, GetRouteHandler(true),
		// Post: gopress.Chain(AuthMiddleware, AdminPermissionMiddleware, AddItemRouteHandler()),
		// ...
	})

	router.AddRoute("/login", gopress.RouteMethods{
		// Post: usermanagement.RouteHandlerFunc(),
		// ...
	})

	// a Route which uses Template in the GET
	router.AddRoute("/user", gopress.RouteMethods{
		// NOT IMPLEMENTED YET
		// Get: gopress.Template(getUserInfo),
		// ...
	})

	// Adding a FileServer serving files under the static routes from path ./public/static
	// NOT IMPLEMENTED YET
	// router.AddFileServer("/static", "./public/static")

	log.Println("Listening...")
	err = router.Start(8080)
	if err != nil {
		log.Panic(err)
	}
}

// This is gopress.TemplateDataHandler function
func getUserInfo(req *gopress.Request) interface{} {
	return User{
		Username:   "mxzinke",
		Repository: "gopress",
	}
}

// Variant 1:
// Create a function, which then returns then a gopress.Handler.
// This is very useful when you want pass something into the function
// to reuse the Middleware/RouteHandler in different situations.

// GetRouteHandler ... Returns a gopress.Handler function,
// depending if request should be logged to database
func GetRouteHandler(saveRequestToDatabase bool) gopress.Handler {
	return gopress.CreateHandler(func(req *gopress.Request, res gopress.Response, next gopress.NextHandler) {
		if saveRequestToDatabase {
			// ... Do something to save the request to database
			// for example to save the path with req.Path
			err := saveRequest(req.Path)
			if err != nil {
				res.SendStatus(gopress.StatusTeapot)
			}
		}

		res.SendJSON(User{
			Username:   "mxzinke",
			Repository: "gopress",
		})
	})
}

// Variant 2:
// Storing the Handlers in a variable
// This kind of creating Handlers make sense,
// when you anyway just want to have a static handler

// StaticLoggingMiddleware ... logs every request with log package
var StaticLoggingMiddleware = gopress.CreateHandler(func(req *gopress.Request, res gopress.Response, next gopress.NextHandler) {
	start := time.Now()

	defer func() {
		log.Println(req.RemoteAddr, req.Method, req.Path, time.Since(start))
	}()

	next()
})

func saveRequest(path string) error {
	log.Printf("Save the request on %s", path)
	c := time.After(20 * time.Millisecond)
	<-c

	return nil
}

// ...
