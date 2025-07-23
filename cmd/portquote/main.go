package main

import (
	"log"
	"net/http"
	"path/filepath"
	"portquote/internal/handlers"
	"portquote/internal/server"
	"portquote/internal/store"
)

func main() {
	dbPath := filepath.Join("internal", "store", "portquote.db")
	db, err := store.NewDB(dbPath)
	if err != nil {
		log.Fatalf("failed to setup db: %v", err)
	}
	defer db.Close()

	router := server.NewRouter(db)

	router.Static("/static/", "web/static")

	router.Handle(http.MethodGet, "/login", handlers.Login)
	router.Handle(http.MethodGet, "/agent/dashboard", handlers.AgentDashboard)
	router.Handle(http.MethodPost, "/agent/dashboard/edit", handlers.AgentDashboardEdit)
	router.Handle(http.MethodGet, "/crew/dashboard", handlers.CrewDashboard)

	addr := ":8080"
	log.Printf("starting server on %s\n", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
