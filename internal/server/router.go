package server

import (
	"context"
	"net/http"
)

type Router struct {
	mux *http.ServeMux
	ctx context.Context
}

func NewRouter(ctx context.Context) *Router {
	return &Router{mux: http.NewServeMux(), ctx: ctx}
}

func (rt *Router) Handle(method, pattern string, handler http.Handler) {
	route := method + " " + pattern
	rt.mux.Handle(route, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r.WithContext(rt.ctx))
	}))
}

func (rt *Router) Static(prefix, dir string) {
	fs := http.FileServer(http.Dir(dir))
	rt.mux.Handle(prefix, http.StripPrefix(prefix, fs))
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rt.mux.ServeHTTP(w, r)
}
