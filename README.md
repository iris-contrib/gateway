# Gateway

[![build status](https://img.shields.io/travis/com/iris-contrib/gateway/master.svg?style=for-the-badge&logo=travis)](https://travis-ci.com/github/iris-contrib/gateway) [![report card](https://img.shields.io/badge/report%20card-a%2B-ff3333.svg?style=for-the-badge)](https://goreportcard.com/report/github.com/iris-contrib/gateway) [![godocs](https://img.shields.io/badge/go-%20docs-488AC7.svg?style=for-the-badge)](https://pkg.go.dev/github.com/iris-contrib/gateway)

Gateway is a simple [iris.Runner](https://github.com/kataras/iris/blob/8e049d77c92fb4550c836e8c64c3052c3d05af15/iris.go#L634). It runs Iris Web Applications through AWS Lambda & API Gateway aka **Serverless**. This includes the [Netlify functions (free and paid)](https://docs.netlify.com/functions/overview/) too. Thanks to [apex/gateway](https://github.com/apex/gateway).

## Installation

The only requirement is the [Go Programming Language](https://golang.org/dl).

```sh
$ go get github.com/iris-contrib/gateway
```

**Until this [PR](https://github.com/apex/gateway/pull/33) is merged, you have to use a `replace statement` inside your go.mod file**:

```text
module my_iris_function

go 1.15

require (
	github.com/iris-contrib/gateway v0.0.0-20200823143335-771cd2392f72
	github.com/kataras/iris/v12 v12.1.9-0.20200823145529-ef5685bf7eeb
)

replace github.com/apex/gateway/v2 v2.0.0-20200703123654-59bba3473042 => github.com/kataras/gateway/v2 v2.0.0-20200823133619-5f644b75fcd5
```

## Getting Started

Simply as:

```go
app := iris.New()
// [...]

runner := gateway.Listen(gateway.Options{})
app.Run(runner)
```

### Netlify

*1.* Create an account on [netlify.com](https://app.netlify.com/signup)

*2.* Link a new website with a repository (GitHub or GitLab, public or private)

*3.* Add a `main.go` in the root of that repository:

```go
// Read and Write JSON only.
package main

func main() {
    app := iris.New()
    app.OnErrorCode(iris.StatusNotFound, notFound)

    app.Get("/", index)
    app.Get("/ping", status)

    // IMPORTANT:
    app.Run(gateway.Listen(gateway.Options{
        URLPathParameter: "path",
    }))
}

func notFound(ctx iris.Context){
    code := ctx.GetStatusCode()
    msg := iris.StatusText(code)
    if err := ctx.GetErr(); err!=nil{
        msg = err.Error(),
    }

    ctx.JSON(iris.Map{
        "Message": msg,
        "Code": code,
    })
}

func index(ctx iris.Context) {
    var req map[string]interface{}
    ctx.ReadJSON(req)
    ctx.JSON(req)
}

func status(ctx iris.Context) {
    ctx.JSON(iris.Map{"Message": "OK"})
}
```

*4.* Create or open the `netlify.toml` file, edit its contents so they look like the following:

```tml
[build]
  publish = "public"
  command = "make build"
  functions = "./functions"
  

[build.environment]
  GO_VERSION = "1.14.7"
  GIMME_GO_VERSION = "1.14.7"
  GO_IMPORT_PATH = "github.com/your_username/your_repo"

[[redirects]]
   from = "/api/*"
   to = '/.netlify/functions/my_iris_function/:splat'
   status = 200
```

**Makefile**

```sh
build:
	go build -o ./functions/my_iris_function
	chmod +x ./functions/my_iris_function
```

*5.* Use `git push` to deploy to Netlify.

The serverless Iris application of will be reachable through _your_site.com/api_, e.g. `https://example.com/api?path=ping`. Have fun!

## License

This software is licensed under the [MIT License](LICENSE).
