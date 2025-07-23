package server

import (
	"net/http"
	"slices"
)

func roleWrap(next http.Handler, allowed []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := CurrentUser(r.Context())
		if user == nil || !contains(allowed, string(user.Role)) {
			redirectToLogin(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func contains(list []string, s string) bool {
	return slices.Contains(list, s)
}

func redirectToLogin(w http.ResponseWriter, r *http.Request) {
	loginURL := "/login"
	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", loginURL)
	} else {
		http.Redirect(w, r, loginURL, http.StatusSeeOther)
	}
}
