package handlers

import (
	"database/sql"
	"net/http"
	"portquote/internal/store"
	"portquote/web/templates"
	"strconv"
)

type DashboardRecord struct {
	Port      store.Port
	Quotation store.Quotation
}

func Dashboard(db *sql.DB, w http.ResponseWriter, r *http.Request) {
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

	switch r.Method {
	case http.MethodGet:
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
	case http.MethodPut:
		id, _ := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
		rate, _ := strconv.ParseFloat(r.FormValue("rate"), 64)
		valid := r.FormValue("valid_until")
		if err := store.UpdateQuotation(db, id, int64(user.ID), rate, valid); err != nil {
			http.Error(w, "failed to update", http.StatusInternalServerError)
			return
		}
		qt, _ := store.GetQuotationByID(db, id, int64(user.ID))
		port, _ := store.GetPortByID(db, qt.PortID)
		rec := DashboardRecord{Port: *port, Quotation: *qt}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		templates.T.ExecuteTemplate(w, "dashboard_row_fragment.html", rec)

	case http.MethodDelete:
		id, _ := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
		if err := store.DeleteQuotation(db, id, int64(user.ID)); err != nil {
			http.Error(w, "failed to delete", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

}
