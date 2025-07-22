package handlers

import (
	"database/sql"
	"net/http"
	"portquote/internal/store"
	"portquote/web/templates"
)

type DashboardRecord struct {
	Port      store.Port
	Quotation store.Quotation
}

func Dashboard(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

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
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		templates.T.ExecuteTemplate(w, "dashboard_fragment.html", []DashboardRecord{})
		return
	}

	quotes, _ := store.GetQuotationsByAgent(db, int64(user.ID))
	var records []DashboardRecord
	for _, qt := range quotes {
		port, _ := store.GetPortByID(db, qt.PortID)
		if port == nil {
			continue
		}
		records = append(records, DashboardRecord{Port: *port, Quotation: qt})
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	templates.T.ExecuteTemplate(w, "dashboard_fragment.html", records)

}
