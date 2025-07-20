package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"portquote/internal/store"
)

type DashboardRecord struct {
	Port      store.Port
	Quotation store.Quotation
}

func Dashboard(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "invalid request method", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := store.GetUserBySession(db, cookie.Value)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
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
}
