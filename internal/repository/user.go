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

type UserRepo interface {
	GetByID(ctx context.Context, id int64) (*User, error)
}

func GetUserByID(db *store.Store, id int64) (*User, error) {
	const q = `
	SELECT id, username, password_hash, role, session
		FROM users
	 WHERE id = ?`
	row := db.QueryRow(q, id)

	u := &User{}
	if err := row.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Role, &u.Session); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return u, nil
}

func GetUserByUsername(db *store.Store, username string) (*User, error) {
	const query = `
    SELECT id, username, password_hash, role, session
      FROM users
     WHERE username = ?
  `

	u := &User{}
	row := db.QueryRow(query, username)
	if err := row.Scan(&u.ID, &u.Username, &u.PasswordHash, &u.Role, &u.Session); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return u, nil
}

func GetUserBySession(db *store.Store, session string) (*User, error) {
	const q = `
	SELECT id, username, password_hash, role, session
		FROM users
	 WHERE session = ?`
	row := db.QueryRow(q, session)

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

func UpdateUserSession(db *store.Store, userID int64, token string) error {
	_, err := db.Exec(`
			UPDATE users
			SET session = ?
			WHERE id = ?`, token, userID)
	return err
}
