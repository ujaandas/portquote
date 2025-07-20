package auth

import (
	"database/sql"

	"github.com/gorilla/securecookie"
)

const cookieName = "session"

type Authenticator struct {
	db     *sql.DB
	cookie *securecookie.SecureCookie
}

func New(db *sql.DB, hashKey, blockKey []byte) *Authenticator {
	cookie := securecookie.New(hashKey, blockKey)
	return &Authenticator{db: db, cookie: cookie}
}
