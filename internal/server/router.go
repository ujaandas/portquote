package server

import (
	"context"
	"net/http"
)

type Router struct {
	mux         *http.ServeMux
	ctx         context.Context
	middlewares []Middleware
}

func NewRouter(ctx context.Context) *Router {
	return &Router{mux: http.NewServeMux(), ctx: ctx}
}

func (rt *Router) Handle(method, pattern string, handler http.Handler, allowedRoles ...string) {
	h := handler

	if len(allowedRoles) > 0 {
		h = roleWrap(h, allowedRoles)
	}

	for i := len(rt.middlewares) - 1; i >= 0; i-- {
		h = rt.middlewares[i](h)
	}

	route := method + " " + pattern

	rt.mux.Handle(route, http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r.WithContext(rt.ctx))
		}))
}

func (rt *Router) Use(mw Middleware) {
	rt.middlewares = append(rt.middlewares, mw)
}

func (rt *Router) Static(prefix, dir string) {
	fs := http.FileServer(http.Dir(dir))
	rt.mux.Handle(prefix, http.StripPrefix(prefix, fs))
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rt.mux.ServeHTTP(w, r)
}
