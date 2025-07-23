package handlers

import (
	"fmt"
	"net/http"
	"portquote/internal/repository"
	"portquote/web/templates"
	"strconv"
)

type DashboardRecord struct {
	Port      repository.Port
	Quotation repository.Quotation
}

func AgentDashboard(users *repository.UserRepo, ports *repository.PortRepo, quotations *repository.QuotationRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		cookie, err := r.Cookie("session_token")
		if err != nil {
			if r.Header.Get("HX-Request") == "true" {
				w.Header().Set("HX-Redirect", "/login")
				return
			}
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		user, err := users.GetBySession(ctx, cookie.Value)
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

			quotes, _ := quotations.GetByAgent(ctx, int64(user.ID))
			var records []DashboardRecord
			for _, qt := range quotes {
				port, _ := ports.GetByID(ctx, qt.PortID)
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
			if err := quotations.UpdateById(ctx, id, int64(user.ID), rate, valid); err != nil {
				http.Error(w, "failed to update", http.StatusInternalServerError)
				return
			}

			qt, _ := quotations.GetById(ctx, id, int64(user.ID))
			port, _ := ports.GetByID(ctx, qt.PortID)
			rec := DashboardRecord{Port: *port, Quotation: *qt}

			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			templates.T.ExecuteTemplate(w, "dashboard_row_fragment.html", rec)

		case http.MethodDelete:
			id, _ := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
			if err := quotations.DeleteById(ctx, id, int64(user.ID)); err != nil {
				http.Error(w, "failed to delete", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusNoContent)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
	}
}

func AgentDashboardEdit(users *repository.UserRepo, ports *repository.PortRepo, quotations *repository.QuotationRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		cookie, err := r.Cookie("session_token")
		if err != nil {
			if r.Header.Get("HX-Request") == "true" {
				w.Header().Set("HX-Redirect", "/login")
				return
			}
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		user, err := users.GetBySession(ctx, cookie.Value)
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
		qt, err := quotations.GetById(ctx, id, int64(user.ID))
		if err != nil || qt == nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		port, _ := ports.GetByID(ctx, qt.PortID)

		rec := DashboardRecord{Port: *port, Quotation: *qt}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		templates.T.ExecuteTemplate(w, "dashboard_edit_fragment.html", rec)
	}
}

type CrewDashboardRecord struct {
	Agent     repository.User
	Quotation repository.Quotation
}

func CrewDashboard(users *repository.UserRepo, ports *repository.PortRepo, quotations *repository.QuotationRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		cookie, err := r.Cookie("session_token")
		if err != nil {
			if r.Header.Get("HX-Request") == "true" {
				w.Header().Set("HX-Redirect", "/login")
				return
			}
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		user, err := users.GetBySession(ctx, cookie.Value)
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
				ports, _ := ports.GetAll(ctx)
				w.Header().Set("Content-Type", "text/html; charset=utf-8")
				templates.T.ExecuteTemplate(w, "crew_dashboard.html", struct {
					Ports []repository.Port
				}{
					Ports: ports,
				})
				return
			}

			pid, _ := strconv.ParseInt(r.URL.Query().Get("port_id"), 10, 64)
			quotes, _ := quotations.GetByPort(ctx, pid)
			var rows []CrewDashboardRecord
			for _, qt := range quotes {
				agent, _ := users.GetByID(ctx, qt.AgentID)
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
}
