package auth

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"net/http"
	"portquote/internal/store"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func Login(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "invalid req method", http.StatusMethodNotAllowed)
		return
	}

	username, password := r.FormValue("username"), r.FormValue("password")
	user, err := store.GetUserByUsername(db, username)

	if err != nil || user == nil || passwordInvalid(password, user.PasswordHash) {
		http.Error(w, "invalid login", http.StatusUnauthorized)
		return
	}

	token, err := generateToken(32)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	if err := store.UpdateUserSession(db, int64(user.ID), token); err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)

}

func passwordInvalid(pw, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw)) != nil
}

func generateToken(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
