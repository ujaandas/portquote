package auth

import (
	"fmt"
	"net/http"
	"portquote/internal/store"

	"golang.org/x/crypto/bcrypt"
)

func (auth *Authenticator) LoginPage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `
    <form method="POST" action="/login">
      Username: <input name="username"><br>
      Password: <input type="password" name="password"><br>
      <button>Login</button>
    </form>
  `)
}

func (auth *Authenticator) Login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	username, password := r.Form.Get("username"), r.Form.Get("password")
	user, err := store.GetUserByUsername(auth.db, username)
	if err != nil || user == nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	session := map[string]int{"user_id": user.ID}
	encoded, err := auth.cookie.Encode(cookieName, session)
	if err != nil {
		http.Error(w, "could not create session", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    encoded,
		Path:     "/",
		HttpOnly: true,
	})
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
