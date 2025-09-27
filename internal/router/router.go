package router

import (
	"path"

	"github.com/cilginc/webserver-go/internal/http"
)

type Handler interface {
	ServeHTTP(ctx *http.ConnContext, req *http.Request) error
}

type HandlerFunc func(ctx *http.ConnContext, req *http.Request) error

func (f HandlerFunc) ServeHTTP(ctx *http.ConnContext, req *http.Request) error {
	return f(ctx, req)
}

type routeKey struct{ method, pattern string }

type Router struct {
	routes map[routeKey]Handler
}

func New() *Router { return &Router{routes: make(map[routeKey]Handler)} }

func (r *Router) Handle(method, pattern string, h Handler) {
	r.routes[routeKey{method, pattern}] = h
}

func (r *Router) Match(method, reqPath string) Handler {
	if h, ok := r.routes[routeKey{method, reqPath}]; ok {
		return h
	}
	for k, h := range r.routes {
		if k.method != method {
			continue
		}
		if k.pattern == "/*" {
			return h
		}
		if k.pattern == "/" {
			return h
		}
		if len(k.pattern) > 1 && k.pattern[len(k.pattern)-1] == '*' {
			base := path.Dir(k.pattern)
			if len(reqPath) >= len(base) && reqPath[:len(base)] == base {
				return h
			}
		}
	}
	return nil
}
