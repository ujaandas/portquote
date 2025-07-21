// internal/handlers/dashboard.go
package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"portquote/internal/store"
	"portquote/web/templates"
)

type DashboardRecord struct {
	Port      store.Port
	Quotation store.Quotation
}

func Dashboard(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		cookie, err := r.Cookie("session_token")
		if err != nil {
			if r.Header.Get("HX-Request") == "true" {
				w.Header().Set("HX-Redirect", "/login")
				return
			}
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		user, err := store.GetUserBySession(db, cookie.Value)
		if err != nil || user == nil {
			if r.Header.Get("HX-Request") == "true" {
				w.Header().Set("HX-Redirect", "/login")
				return
			}
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		if r.Header.Get("HX-Request") != "true" {
			templates.T.ExecuteTemplate(w, "dashboard.html", nil)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if user.Role != "agent" {
			json.NewEncoder(w).Encode([]DashboardRecord{})
			return
		}

		quotes, err := store.GetQuotationsByAgent(db, int64(user.ID))
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}

		var resp []DashboardRecord
		for _, qt := range quotes {
			port, err := store.GetPortByID(db, qt.PortID)
			if err != nil {
				http.Error(w, "server error", http.StatusInternalServerError)
				return
			}
			if port == nil {
				continue
			}
			resp = append(resp, DashboardRecord{
				Port:      *port,
				Quotation: qt,
			})
		}

		json.NewEncoder(w).Encode(resp)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
