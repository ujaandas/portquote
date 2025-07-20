package store

import (
	"database/sql"
	"errors"
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

func GetUserByUsername(db *sql.DB, username string) (*User, error) {
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

func UpdateUserSession(db *sql.DB, userID int64, token string) error {
	_, err := db.Exec(`
			UPDATE users
			SET session = ?
			WHERE id = ?`, token, userID)
	return err
}
