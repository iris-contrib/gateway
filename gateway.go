package gateway

import (
	"net/http"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/router"

	"github.com/apex/gateway/v2"
	"github.com/aws/aws-lambda-go/lambda"
)

// Options holds the Listen options. All fields are optional.
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

// Listen is given to iris.Application.Run to run the application
// using the Apex Gateway. That allows Iris-powered web application
// to be deployed and ran on host services like Netlify and Amazon AWS.
func Listen(opts Options) iris.Runner {
	return func(app *iris.Application) error {
		if opts.URLPathParameter != "" {
			wrapper := urlToPath(opts.URLPathParameter)
			app.WrapRouter(wrapper)
			app.RefreshRouter()
		}

		g := gateway.NewGateway(app)
		lambda.StartHandler(g)
		return nil
	}
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
