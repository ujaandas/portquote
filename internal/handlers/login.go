package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"portquote/internal/repository"
	"portquote/web/templates"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func LoginPOST(users *repository.UserRepo,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		username, password := r.FormValue("username"), r.FormValue("password")
		user, err := users.GetByUsername(ctx, username)
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
		if err := users.UpdateSession(ctx, int64(user.ID), token); err != nil {
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
	}
}

func LoginGET() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		templates.T.ExecuteTemplate(w, "login.html", nil)
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
