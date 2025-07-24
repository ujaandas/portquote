package handlers

import (
	"fmt"
	"net/http"
	"portquote/internal/repository"
	"portquote/internal/server"
	"portquote/web/templates"
	"strconv"
)

type DashboardRecord struct {
	Port      repository.Port
	Quotation repository.Quotation
}

func AgentDashboardGET(ports *repository.PortRepo, quotations *repository.QuotationRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user := server.CurrentUser(ctx)

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

	}
}

func AgentDashboardPUT(ports *repository.PortRepo, quotations *repository.QuotationRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user := server.CurrentUser(ctx)

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
	}
}

func AgentDashboardDELETE(ports *repository.PortRepo, quotations *repository.QuotationRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user := server.CurrentUser(ctx)

		id, _ := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
		if err := quotations.DeleteById(ctx, id, int64(user.ID)); err != nil {
			http.Error(w, "failed to delete", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func AgentDashboardEditGET(users *repository.UserRepo, ports *repository.PortRepo, quotations *repository.QuotationRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		user := server.CurrentUser(ctx)

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

func CrewDashboardGET(users *repository.UserRepo, ports *repository.PortRepo, quotations *repository.QuotationRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

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
	}
}
