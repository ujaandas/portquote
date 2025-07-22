// handlers/dashboard_edit.go

package handlers

import (
	"database/sql"
	"net/http"
	"portquote/internal/store"
	"portquote/web/templates"
	"strconv"
)

func DashboardEdit(db *sql.DB, w http.ResponseWriter, r *http.Request) {
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

	if r.Header.Get("HX-Request") != "true" || r.Method != http.MethodGet {
		http.Error(w, "not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, _ := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	qt, err := store.GetQuotationByID(db, id, int64(user.ID))
	if err != nil || qt == nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	port, _ := store.GetPortByID(db, qt.PortID)

	rec := DashboardRecord{Port: *port, Quotation: *qt}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	templates.T.ExecuteTemplate(w, "dashboard_edit_fragment.html", rec)
}
