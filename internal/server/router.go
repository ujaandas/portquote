package server

import (
	"database/sql"
	"net/http"
)

type DBHandlerFunc func(db *sql.DB, w http.ResponseWriter, r *http.Request)

type Router struct {
	mux *http.ServeMux
	db  *sql.DB
}

func NewRouter(db *sql.DB) *Router {
	return &Router{
		mux: http.NewServeMux(),
		db:  db,
	}
}

func (rt *Router) Handle(method, pattern string, handler DBHandlerFunc) {
	rt.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handler(rt.db, w, r)
	})
}

func (rt *Router) Static(prefix, dir string) {
	fs := http.FileServer(http.Dir(dir))
	rt.mux.Handle(prefix, http.StripPrefix(prefix, fs))
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rt.mux.ServeHTTP(w, r)
}
