package main

import (
	"context"
	"log"
	"net/http"
	"path/filepath"
	"portquote/internal/handlers"
	"portquote/internal/repository"
	"portquote/internal/server"
	"portquote/internal/store"
)

func main() {
	baseCtx := context.Background()

	dbPath := filepath.Join("internal", "store", "portquote.db")
	db, err := store.NewDB(baseCtx, dbPath)
	if err != nil {
		log.Fatalf("failed to setup db: %v", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepo(db)
	portRepo := repository.NewPortRepo(db)
	quoteRepo := repository.NewQuotationRepo(db)

	router := server.NewRouter(baseCtx)

	router.Static("/static/", "web/static")

	router.Handle(http.MethodGet, "/login", handlers.LoginHandlerGET())
	router.Handle(http.MethodPost, "/login", handlers.LoginHandlerPOST(userRepo))
	router.Handle(http.MethodGet, "/agent/dashboard", handlers.AgentDashboard(userRepo, portRepo, quoteRepo))
	router.Handle(http.MethodPost, "/agent/dashboard/edit", handlers.AgentDashboardEdit(userRepo, portRepo, quoteRepo))
	router.Handle(http.MethodGet, "/crew/dashboard", handlers.CrewDashboard(userRepo, portRepo, quoteRepo))

	addr := ":8080"
	log.Printf("starting server on %s\n", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
