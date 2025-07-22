package main

import (
	"log"
	"net/http"
	"path/filepath"
	"portquote/internal/handlers"
	"portquote/internal/store"
)

func main() {
	dbPath := filepath.Join("internal", "store", "portquote.db")
	db, err := store.NewDB(dbPath)
	if err != nil {
		log.Fatalf("failed to setup db: %v", err)
	}
	defer db.Close()

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		handlers.Login(db, w, r)
	})

	mux.HandleFunc("/agent/dashboard", func(w http.ResponseWriter, r *http.Request) {
		handlers.AgentDashboard(db, w, r)
	})

	mux.HandleFunc("/agent/dashboard/edit", func(w http.ResponseWriter, r *http.Request) {
		handlers.AgentDashboardEdit(db, w, r)
	})

	addr := ":8080"
	log.Printf("starting server on %s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
