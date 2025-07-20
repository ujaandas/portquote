package store

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema.sql
var schemaFS embed.FS

//go:embed seed.sql
var seedFS embed.FS

func NewDB(path string) (*sql.DB, error) {
	if err := os.Remove(path); err != nil {
		log.Fatalf("remove database file '%s': %v\n", path, err)
	}

	db, err := sql.Open("sqlite3", path+"?_foreign_keys=1")
	if err != nil {
		return nil, fmt.Errorf("open sqlite3: %w", err)
	}

	if err := migrate(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}

	if err := seed(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("seed: %w", err)
	}

	return db, nil
}

func migrate(db *sql.DB) error {
	data, err := fs.ReadFile(schemaFS, "schema.sql")

	if err != nil {
		return fmt.Errorf("read schema.sql: %w", err)
	}

	if _, err := db.Exec(string(data)); err != nil {
		return fmt.Errorf("exec schema: %w", err)
	}

	return nil
}

func seed(db *sql.DB) error {
	data, err := fs.ReadFile(seedFS, "seed.sql")
	if err != nil {
		return fmt.Errorf("read seed.sql: %w", err)
	}

	if _, err := db.Exec(string(data)); err != nil {
		return fmt.Errorf("exec seed.sql: %w", err)
	}

	return nil
}
