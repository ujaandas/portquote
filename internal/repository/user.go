package repository

import (
	"context"
	"database/sql"
	"errors"
	"portquote/internal/store"
)

type UserRole string

const (
	AdminR UserRole = "admin"
	AgentR UserRole = "agent"
	UserR  UserRole = "user"
)

type User struct {
	ID           int
	Username     string
	PasswordHash string
	Role         UserRole
	Session      string
}
type UserRepo struct {
	db *store.Store
}

func NewUserRepo(db *store.Store) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) GetByID(ctx context.Context, id int64) (*User, error) {
	const q = `
	SELECT id, username, password_hash, role, session
		FROM users
	 WHERE id = ?`
	row := r.db.QueryRowContext(ctx, q, id)

	u := &User{}
	if err := row.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Role, &u.Session); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return u, nil
}

func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*User, error) {
	const query = `
    SELECT id, username, password_hash, role, session
      FROM users
     WHERE username = ?
  `

	u := &User{}
	row := r.db.QueryRow(query, username)
	if err := row.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Role, &u.Session); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return u, nil
}

func (r *UserRepo) GetBySession(ctx context.Context, session string) (*User, error) {
	const q = `
	SELECT id, username, password_hash, role, session
		FROM users
	 WHERE session = ?`
	row := r.db.QueryRow(q, session)

	u := &User{}
	if err := row.Scan(
		&u.ID,
		&u.Username,
		&u.PasswordHash,
		&u.Role,
		&u.Session,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return u, nil
}

func (r *UserRepo) UpdateSession(ctx context.Context, userID int64, token string) error {
	_, err := r.db.Exec(`
			UPDATE users
			SET session = ?
			WHERE id = ?`, token, userID)
	return err
}
