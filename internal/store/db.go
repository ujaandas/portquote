package store

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
)

var schemaFS embed.FS

func NewDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path+"?_foreign_keys=1")
	if err != nil {
		return nil, fmt.Errorf("open sqlite3: %w", err)
	}

	if err := migrate(db); err != nil {
		db.Close()
		return nil, fmt.Errorf("migrate: %w", err)
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
