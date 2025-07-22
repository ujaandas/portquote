package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"portquote/internal/store"
	"portquote/web/templates"
	"strconv"
)

type DashboardRecord struct {
	Port      store.Port
	Quotation store.Quotation
}

func AgentDashboard(db *sql.DB, w http.ResponseWriter, r *http.Request) {
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
	if err != nil || user == nil || user.Role != "agent" {
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

func AgentDashboardEdit(db *sql.DB, w http.ResponseWriter, r *http.Request) {
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
	if err != nil || user == nil || user.Role != "agent" {
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

type CrewDashboardRecord struct {
	Agent     store.User
	Quotation store.Quotation
}

func CrewDashboard(db *sql.DB, w http.ResponseWriter, r *http.Request) {
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
	if err != nil || user == nil || user.Role != "crew" {
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
			ports, _ := store.GetAllPorts(db)
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			templates.T.ExecuteTemplate(w, "crew_dashboard.html", struct {
				Ports []store.Port
			}{
				Ports: ports,
			})
			return
		}

		pid, _ := strconv.ParseInt(r.URL.Query().Get("port_id"), 10, 64)
		quotes, _ := store.GetQuotationsByPort(db, pid)
		var rows []CrewDashboardRecord
		for _, qt := range quotes {
			agent, _ := store.GetUserByID(db, qt.AgentID)
			if agent == nil {
				continue
			}
			rows = append(rows, CrewDashboardRecord{
				Agent:     *agent,
				Quotation: qt,
			})
		}

		fmt.Print(rows)

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		templates.T.ExecuteTemplate(w, "crew_dashboard_fragment.html", rows)

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}
