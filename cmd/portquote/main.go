package main

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"portquote/internal/handlers"
	"portquote/internal/repository"
	"portquote/internal/server"
	"portquote/internal/store"
)

func main() {
	baseCtx := context.Background()
	l := server.NewLogger()

	dbPath := filepath.Join("internal", "store", "portquote.db")
	db, err := store.NewDB(baseCtx, dbPath)
	if err != nil {
		l.Error("DB setup failed", err)
	}
	defer db.Close()
	l.Info("DB setup complete")

	userRepo := repository.NewUserRepo(db)
	portRepo := repository.NewPortRepo(db)
	quoteRepo := repository.NewQuotationRepo(db)
	l.Info("Repo setup complete")

	router := server.NewRouter(baseCtx)
	l.Info("Router setup complete")

	router.Use(server.SessionMiddleware(userRepo))
	router.Use(server.LoggingMiddleware)
	l.Info("Middleware setup complete")

	router.Static("/static/", "web/static")
	l.Info("Static routes registered")

	router.Handle(http.MethodGet, "/login", handlers.LoginGET())
	router.Handle(http.MethodPost, "/login", handlers.LoginPOST(userRepo))
	router.Handle(http.MethodGet, "/agent/dashboard", handlers.AgentDashboard(userRepo, portRepo, quoteRepo), "agent", "admin")
	router.Handle(http.MethodPost, "/agent/dashboard/edit", handlers.AgentDashboardEdit(userRepo, portRepo, quoteRepo), "agent", "admin")
	router.Handle(http.MethodGet, "/crew/dashboard", handlers.CrewDashboard(userRepo, portRepo, quoteRepo), "crew", "admin")
	l.Info("Other routes registered")

	addr := ":8080"
	l.Info(fmt.Sprintf("Server started on http://localhost%s", addr))
	if err := http.ListenAndServe(addr, router); err != nil {
		l.Error("Server failed catastrophically", err)
	}
}
