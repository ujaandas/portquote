package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"portquote/internal/repository"
	"portquote/internal/store"
	"portquote/web/templates"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func Login(db *store.Store, w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		templates.T.ExecuteTemplate(w, "login.html", nil)
		return

	case http.MethodPost:
		username, password := r.FormValue("username"), r.FormValue("password")
		user, err := repository.GetUserByUsername(db, username)
		if err != nil || user == nil || passwordInvalid(password, user.PasswordHash) {
			w.WriteHeader(http.StatusUnauthorized)
			templates.T.ExecuteTemplate(w, "login.html", map[string]string{
				"Error": "Invalid username or password",
			})
			return
		}

		token, err := generateToken(32)
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
		if err := repository.UpdateUserSession(db, int64(user.ID), token); err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_token",
			Value:    token,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
		})

		if r.Header.Get("HX-Request") == "true" {
			w.Header().Set("HX-Redirect", fmt.Sprintf("/%s/dashboard", user.Role))
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/%s/dashboard", user.Role), http.StatusSeeOther)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
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
