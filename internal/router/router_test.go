package router

import (
	"testing"

	"github.com/cilginc/webserver-go/internal/http"
)

type dummyHandler struct{}

func (d *dummyHandler) ServeHTTP(ctx *http.ConnContext, req *http.Request) error { return nil }

func TestRouterMatch(t *testing.T) {
	r := New()
	h := &dummyHandler{}
	r.Handle("GET", "/foo", h)

	if got := r.Match("GET", "/foo"); got == nil {
		t.Error("expected handler match for /foo")
	}
	if got := r.Match("GET", "/bar"); got != nil {
		t.Error("expected no handler for /bar")
	}
}
