# GoPress

GoPress is a super easy to use router with a stolen idea from the nice Express.js framework used by a wide majority of Node.js users.

The package makes internally usage of the net/http package and provides a easy usage of routes, middlewares and other features like file server and template serving.

## Installation

```bash
$ go get -u github.com/mxzinke/gopress
```

## Build with GoPress

Building a easy server with some routes:

```golang
package main

import (
    "os"

    "github.com/mxzinke/gopress"
)

func main() {
    // Passed Middleware to the router will be always executed to 
    router := gopress.NewRouter(gopress.RouterSettings{
        TemplatesPath: "./public/templates",
        SSL: true,
        SSLCertFile: "./path/to/cert/file",
        SSLKeyFile: "./path/to/key/file",
    }, gopress.Chain(LoggingMiddleware(os.GetEnv("LOGGING_PATH")), OtherMiddleware))

    router.AddRoute("/", gopress.RouteMethods{
        Get: gopress.Chain(AuthMiddleware, GetRouteHandler()),
        Post: gopress.Chain(AuthMiddleware, AdminPermissionMiddleware, AddItemRouteHandler()),
        // ...
    })
    
    router.AddRoute("/login", gopress.RouteMethods{
        Post: usermanagement.RouteHandlerFunc(),
        // ...
    })

    // a Route which uses Template in the GET 
    router.AddRoute("/user", gopress.RouteMethods{
        Get: gopress.Template(getUserInfo),
        // ...
    })

    // Adding a FileServer serving files under the static routes from path ./public/static
    router.AddFileServer("/static", "./public/static")

    port := os.GetEnv("PORT")
    router.Start(8443)
}

// This is gopress.TemplateDataHandler function
func getUserInfo(req *gopress.Request) interface{} {
    return models.User{
        Username: "mxzinke",
        Repository: "gopress"
    }
}

// ...

```

Also this example is using a file server which serves all files which are available under the path `./public/static/`.

Additionally templates can be saved under the path `/public/templates/`. Templates are very useful, when you pass a data *TemplateDataHandler* function. If non available pass just `nil` as argument. The `gopress.Template(dataHandler)` is just a predefined Route Handler which solves serving templates on a easy way.

To use the functionality of a Middleware you can chain multiple Middleware/RouteHandler functions. The first passed Handler will also be executed first at a normal request. In following code I will show you how to build a RouteHandler or Middleware.

## Handler for Routes and Middleware

```golang
package yourpackage

import (
    "log"

    "github.com/mxzinke/gopress"
)

// Variant 1:
// Create a function, which then returns then a gopress.Handler.
// This is very useful when you want pass something into the function
// to reuse the Middleware/RouteHandler in different situations.

// GetRouteHandler ... Returns a gopress.Handler function,
// depending if request should be logged to database
func GetRouteHandler(saveRequestToDatabase boolean) gopress.Handler {
    return gopress.CreateHandler(func(req *gopress.Request, res gopress.Response) {
        if saveRequestToDatabase {
            // ... Do something to save the request to database
            // for example to save the path with req.Path 

            err := saveRequest(req.Path)
            if err != nil {
                res.sendStatus(gopress.StatusServerError)
            }
        }

        res.sendJSON(models.User{
            Username: "mxzinke",
            Repository: "gopress"
        })

        next()
    })
}

// Variant 2:
// Storing the Handlers in a variable
// This kind of creating Handlers make sense,
// when you anyway just want to have a static handler

// StaticLoggingMiddleware ... logs every request with log package
var StaticLoggingMiddleware = gopress.CreateHandler(func(req *gopress.Request, res gopress.Response) {
	start := time.Now()
    defer func() {
        log.Println(req.RemoteAddr, req.Method, req.Path, time.Since(start))
    }()
})

// ...
```

In core functionality **Middleware and RequestHandler are the same**. Thats the reason, why they just called `gopress.Handler`. As described you have two methods of creating such handler.
