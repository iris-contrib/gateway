package gateway

import (
	"bytes"
	"testing"
)

func Test_JSON_isTextMime(t *testing.T) {
	assertEqual(t, isTextMime("application/json"), true)
	assertEqual(t, isTextMime("application/json; charset=utf-8"), true)
	assertEqual(t, isTextMime("Application/JSON"), true)
}

func Test_XML_isTextMime(t *testing.T) {
	assertEqual(t, isTextMime("application/xml"), true)
	assertEqual(t, isTextMime("application/xml; charset=utf-8"), true)
	assertEqual(t, isTextMime("ApPlicaTion/xMl"), true)
}

func TestResponseWriter_Header(t *testing.T) {
	w := NewResponse()
	w.Header().Set("Foo", "bar")
	w.Header().Set("Bar", "baz")

	var buf bytes.Buffer
	w.header.Write(&buf)

	assertEqual(t, "Bar: baz\r\nFoo: bar\r\n", buf.String())
}

func TestResponseWriter_multiHeader(t *testing.T) {
	w := NewResponse()
	w.Header().Set("Foo", "bar")
	w.Header().Set("Bar", "baz")
	w.Header().Add("X-APEX", "apex1")
	w.Header().Add("X-APEX", "apex2")

	var buf bytes.Buffer
	w.header.Write(&buf)

	assertEqual(t, "Bar: baz\r\nFoo: bar\r\nX-Apex: apex1\r\nX-Apex: apex2\r\n", buf.String())
}

func TestResponseWriter_Write_text(t *testing.T) {
	types := []string{
		"text/x-custom",
		"text/plain",
		"text/plain; charset=utf-8",
		"application/json",
		"application/json; charset=utf-8",
		"application/xml",
		"image/svg+xml",
	}

	for _, kind := range types {
		t.Run(kind, func(t *testing.T) {
			w := NewResponse()
			w.Header().Set("Content-Type", kind)
			w.Write([]byte("hello world\n"))

			e := w.End()
			assertEqual(t, 200, e.StatusCode)
			assertEqual(t, "hello world\n", e.Body)
			assertEqual(t, kind, e.Headers["Content-Type"])
			assertFalse(t, e.IsBase64Encoded)
			assertTrue(t, <-w.CloseNotify())
		})
	}
}

func TestResponseWriter_Write_binary(t *testing.T) {
	w := NewResponse()
	w.Header().Set("Content-Type", "image/png")
	w.Write([]byte("data"))

	e := w.End()
	assertEqual(t, 200, e.StatusCode)
	assertEqual(t, "ZGF0YQ==", e.Body)
	assertEqual(t, "image/png", e.Headers["Content-Type"])
	assertTrue(t, e.IsBase64Encoded)
}

func TestResponseWriter_Write_gzip(t *testing.T) {
	w := NewResponse()
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Encoding", "gzip")
	w.Write([]byte("data"))

	e := w.End()
	assertEqual(t, 200, e.StatusCode)
	assertEqual(t, "ZGF0YQ==", e.Body)
	assertEqual(t, "text/plain", e.Headers["Content-Type"])
	assertTrue(t, e.IsBase64Encoded)
}

func TestResponseWriter_WriteHeader(t *testing.T) {
	w := NewResponse()
	w.WriteHeader(404)
	w.Write([]byte("Not Found\n"))

	e := w.End()
	assertEqual(t, 404, e.StatusCode)
	assertEqual(t, "Not Found\n", e.Body)
	assertEqual(t, "text/plain; charset=utf8", e.Headers["Content-Type"])
}
