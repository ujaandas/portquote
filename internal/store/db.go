package store

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations/*.sql
//go:embed seeds/*.sql
var sqlFS embed.FS

type Store struct {
	*sql.DB
}

func NewDB(ctx context.Context, path string) (*Store, error) {
	if path != ":memory:" {
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			return nil, fmt.Errorf("remove database file %q: %w", path, err)
		}
	}

	db, err := sql.Open("sqlite3", path+"?_foreign_keys=1")
	if err != nil {
		return nil, fmt.Errorf("open sqlite3: %w", err)
	}

	store := &Store{DB: db}

	if err := runSQLFiles(ctx, store, "migrations"); err != nil {
		store.Close()
		return nil, fmt.Errorf("migrate: %w", err)
	}
	if err := runSQLFiles(ctx, store, "seeds"); err != nil {
		store.Close()
		return nil, fmt.Errorf("seed: %w", err)
	}

	return store, nil
}

func runSQLFiles(ctx context.Context, s *Store, subdir string) error {
	entries, err := fs.ReadDir(sqlFS, subdir)
	if err != nil {
		return fmt.Errorf("read %s: %w", subdir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !fs.ValidPath(entry.Name()) {
			continue
		}
		sqlBytes, err := sqlFS.ReadFile(fmt.Sprintf("%s/%s", subdir, entry.Name()))
		if err != nil {
			return fmt.Errorf("read %s/%s: %w", subdir, entry.Name(), err)
		}
		if _, err := s.ExecContext(ctx, string(sqlBytes)); err != nil {
			return fmt.Errorf("exec %s/%s: %w", subdir, entry.Name(), err)
		}
	}
	return nil
}
