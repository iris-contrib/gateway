package gateway

import (
	"net/http"

	"github.com/iris-contrib/gateway/gateway"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/router"
)

// Options holds the gateway options. All fields are optional.
type Options struct {
	// When not empty then a URL parameter of that key will be used to route the requests,
	// changes the default Iris Router behavior based on the Request's URI's Path, e.g. "path".
	//
	// This is extremely useful when the deployment allows only one request path to be ran
	// under a particular lambda function.
	//
	// Defaults to empty.
	URLPathParameter string
}

// New returns a pair of iris Runner and Configurator
// to convert the http application to a lambda function
// using the Apex Gateway. That allows Iris-powered web application
// to be deployed and ran on host services like Netlify and Amazon AWS.
//
// Usage:
// app := iris.New()
// [...routes]
// runner, configurator := gateway.New(gateway.Options{URLPathParameter: "path"})
// app.Run(runner, configurator)
//
// Get the original API Gateway Request object through:
// req, ok := gateway.GetRequest(ctx.Request().Context())
func New(opts Options) (iris.Runner, iris.Configurator) {
	runner := func(app *iris.Application) error {
		g := gateway.NewGateway(app)
		lambda.StartHandler(g)
		return nil
	}

	configurator := func(app *iris.Application) {
		if opts.URLPathParameter != "" {
			wrapper := urlToPath(opts.URLPathParameter)
			app.WrapRouter(wrapper)
		}
	}

	return runner, configurator
}

func urlToPath(key string) router.WrapperFunc {
	return func(w http.ResponseWriter, r *http.Request, router http.HandlerFunc) {
		req, _ := gateway.GetRequest(r.Context())
		path := req.QueryStringParameters["path"]
		if len(path) > 0 {
			if path[0] != '/' {
				path = "/" + path
			}
			r.URL.Path = path
			r.URL.RawPath = path
			r.RequestURI = path
		}

		router(w, r)
	}
}
