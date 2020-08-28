package gateway

import (
	"context"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestDecodeRequest_path(t *testing.T) {
	e := events.APIGatewayV2HTTPRequest{
		RawPath: "/pets/luna",
	}

	r, err := NewRequest(context.Background(), e)
	assertNoError(t, err)

	assertEqual(t, "GET", r.Method)
	assertEqual(t, `/pets/luna`, r.URL.Path)
	assertEqual(t, `/pets/luna`, r.URL.String())
}

func TestDecodeRequest_method(t *testing.T) {
	e := events.APIGatewayV2HTTPRequest{
		RawPath: "/pets/luna",
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Method: "DELETE",
				Path:   "/pets/luna",
			},
		},
	}

	r, err := NewRequest(context.Background(), e)
	assertNoError(t, err)

	assertEqual(t, "DELETE", r.Method)
}

func TestDecodeRequest_queryString(t *testing.T) {
	e := events.APIGatewayV2HTTPRequest{
		RawPath:        "/pets",
		RawQueryString: "fields=name%2Cspecies&order=desc",
		QueryStringParameters: map[string]string{
			"order":  "desc",
			"fields": "name,species",
		},
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Method: "GET",
				Path:   "/pets",
			},
		},
	}

	r, err := NewRequest(context.Background(), e)
	assertNoError(t, err)

	assertEqual(t, `/pets?fields=name%2Cspecies&order=desc`, r.URL.String())
	assertEqual(t, `desc`, r.URL.Query().Get("order"))
}

func TestDecodeRequest_multiValueQueryString(t *testing.T) {
	e := events.APIGatewayV2HTTPRequest{
		RawPath:        "/pets",
		RawQueryString: "fields=name%2Cspecies&multi_arr%5B%5D=arr1&multi_arr%5B%5D=arr2&multi_fields=name&multi_fields=species&order=desc",
		QueryStringParameters: map[string]string{
			"multi_fields": strings.Join([]string{"name", "species"}, ","),
			"multi_arr[]":  strings.Join([]string{"arr1", "arr2"}, ","),
			"order":        "desc",
			"fields":       "name,species",
		},
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Method: "GET",
				Path:   "/pets",
			},
		},
	}

	r, err := NewRequest(context.Background(), e)
	assertNoError(t, err)

	assertEqual(t, `/pets?fields=name%2Cspecies&multi_arr%5B%5D=arr1&multi_arr%5B%5D=arr2&multi_fields=name&multi_fields=species&order=desc`, r.URL.String())
	assertEqual(t, []string{"name", "species"}, r.URL.Query()["multi_fields"])
	assertEqual(t, []string{"arr1", "arr2"}, r.URL.Query()["multi_arr[]"])
}

func TestDecodeRequest_remoteAddr(t *testing.T) {
	e := events.APIGatewayV2HTTPRequest{
		RawPath: "/pets",
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Method:   "GET",
				Path:     "/pets",
				SourceIP: "1.2.3.4",
			},
		},
	}

	r, err := NewRequest(context.Background(), e)
	assertNoError(t, err)

	assertEqual(t, `1.2.3.4`, r.RemoteAddr)
}

func TestDecodeRequest_header(t *testing.T) {
	e := events.APIGatewayV2HTTPRequest{
		RawPath: "/pets",
		Body:    `{ "name": "Tobi" }`,
		Headers: map[string]string{
			"Content-Type": "application/json",
			"X-Foo":        "bar",
			"Host":         "example.com",
		},
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			RequestID: "1234",
			Stage:     "prod",
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Path:   "/pets",
				Method: "POST",
			},
		},
	}

	r, err := NewRequest(context.Background(), e)
	assertNoError(t, err)

	assertEqual(t, `example.com`, r.Host)
	assertEqual(t, `prod`, r.Header.Get("X-Stage"))
	assertEqual(t, `1234`, r.Header.Get("X-Request-Id"))
	assertEqual(t, `18`, r.Header.Get("Content-Length"))
	assertEqual(t, `application/json`, r.Header.Get("Content-Type"))
	assertEqual(t, `bar`, r.Header.Get("X-Foo"))
}

func TestDecodeRequest_multiHeader(t *testing.T) {
	e := events.APIGatewayV2HTTPRequest{
		RawPath: "/pets",
		Body:    `{ "name": "Tobi" }`,
		Headers: map[string]string{
			"X-APEX":       strings.Join([]string{"apex1", "apex2"}, ","),
			"X-APEX-2":     strings.Join([]string{"apex-1", "apex-2"}, ","),
			"Content-Type": "application/json",
			"X-Foo":        "bar",
			"Host":         "example.com",
		},
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			RequestID: "1234",
			Stage:     "prod",
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Path:   "/pets",
				Method: "POST",
			},
		},
	}

	r, err := NewRequest(context.Background(), e)
	assertNoError(t, err)

	assertEqual(t, `example.com`, r.Host)
	assertEqual(t, `prod`, r.Header.Get("X-Stage"))
	assertEqual(t, `1234`, r.Header.Get("X-Request-Id"))
	assertEqual(t, `18`, r.Header.Get("Content-Length"))
	assertEqual(t, `application/json`, r.Header.Get("Content-Type"))
	assertEqual(t, `bar`, r.Header.Get("X-Foo"))
	assertEqual(t, []string{"apex1", "apex2"}, r.Header["X-Apex"])
	assertEqual(t, []string{"apex-1", "apex-2"}, r.Header["X-Apex-2"])
}

func TestDecodeRequest_body(t *testing.T) {
	e := events.APIGatewayV2HTTPRequest{
		RawPath: "/pets",
		Body:    `{ "name": "Tobi" }`,
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Method: "POST",
				Path:   "/pets",
			},
		},
	}

	r, err := NewRequest(context.Background(), e)
	assertNoError(t, err)

	b, err := ioutil.ReadAll(r.Body)
	assertNoError(t, err)

	assertEqual(t, `{ "name": "Tobi" }`, string(b))
}

func TestDecodeRequest_bodyBinary(t *testing.T) {
	e := events.APIGatewayV2HTTPRequest{
		RawPath:         "/pets",
		Body:            `aGVsbG8gd29ybGQK`,
		IsBase64Encoded: true,
		RequestContext: events.APIGatewayV2HTTPRequestContext{
			HTTP: events.APIGatewayV2HTTPRequestContextHTTPDescription{
				Method: "POST",
				Path:   "/pets",
			},
		},
	}

	r, err := NewRequest(context.Background(), e)
	assertNoError(t, err)

	b, err := ioutil.ReadAll(r.Body)
	assertNoError(t, err)

	assertEqual(t, "hello world\n", string(b))
}

func TestDecodeRequest_context(t *testing.T) {
	var key = struct{}{}
	e := events.APIGatewayV2HTTPRequest{}
	ctx := context.WithValue(context.Background(), key, "value")
	r, err := NewRequest(ctx, e)
	assertNoError(t, err)
	v := r.Context().Value(key)
	assertEqual(t, "value", v)
}
