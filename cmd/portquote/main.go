package main

import (
	"log"
	"net/http"
	"path/filepath"
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

	// fs := http.FileServer(http.Dir("web/static"))
	// mux.Handle("/static/", http.StripPrefix("/static/", fs))

	addr := ":8080"
	log.Printf("starting server on %s\n", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server failed: %v", err)
	}

}
