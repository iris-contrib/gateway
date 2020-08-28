package gateway

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello World from Go")
}

func TestGateway_Invoke(t *testing.T) {

	e := []byte(`{"version": "2.0", "rawPath": "/pets/luna", "requestContext": {"http": {"method": "POST"}}}`)

	gw := NewGateway(http.HandlerFunc(hello))

	payload, err := gw.Invoke(context.Background(), e)
	assertNoError(t, err)
	assertJSONEqual(t, `{"body":"Hello World from Go\n", "cookies": null, "headers":{"Content-Type":"text/plain; charset=utf8"}, "multiValueHeaders":{}, "statusCode":200}`, string(payload))
}

// Let's make this package as light as possible, lambda is already heavy enough.

func assertNoError(t *testing.T, err error) {
	if err != nil {
		fail(t, "unexpected error: %v", err)
	}
}

func assertJSONEqual(t *testing.T, expected, got string) {
	var exp, gt interface{}

	if err := json.Unmarshal([]byte(expected), &exp); err != nil {
		fail(t, "expected: %s: invalid json: %v", expected, err)
	}

	if err := json.Unmarshal([]byte(got), &gt); err != nil {
		fail(t, "got: %s: invalid json: %v", got, err)
	}

	assertEqual(t, exp, gt)
}

func assertTrue(t *testing.T, v bool) {
	if !v {
		fail(t, "expected to be true but got false")
	}
}

func assertFalse(t *testing.T, v bool) {
	if v {
		fail(t, "expected to be false but got true")
	}
}

func assertEqual(t *testing.T, expected, got interface{}) {
	if !isEqual(expected, got) {
		fail(t, "expected:\n%#+v \nbut got:\n%#+v", expected, got)
	}
}

func fail(t *testing.T, format string, args ...interface{}) {
	t.Errorf(format, args...)
	t.FailNow()
}

func isEqual(expected, got interface{}) bool {
	if expected == nil || got == nil {
		return expected == got
	}

	exp, ok := expected.([]byte)
	if !ok {
		return reflect.DeepEqual(expected, got)
	}

	act, ok := got.([]byte)
	if !ok {
		return false
	}
	if exp == nil || act == nil {
		return exp == nil && act == nil
	}
	return bytes.Equal(exp, act)
}

//
