package store

import (
	"database/sql"
	"time"
)

type Quotation struct {
	ID         int64
	AgentID    int64
	PortID     int64
	Rate       float64
	ValidUntil string
	UpdatedAt  time.Time
}

func GetQuotationsByAgent(db *sql.DB, agentID int64) ([]Quotation, error) {
	const q = `
        SELECT id, agent_id, port_id, rate, valid_until, updated_at
          FROM quotations
         WHERE agent_id = ?`
	rows, err := db.Query(q, agentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Quotation
	for rows.Next() {
		var qt Quotation
		if err := rows.Scan(
			&qt.ID,
			&qt.AgentID,
			&qt.PortID,
			&qt.Rate,
			&qt.ValidUntil,
			&qt.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, qt)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func UpsertQuotation(db *sql.DB, q Quotation) error {
	const stmt = `
        INSERT INTO quotations (agent_id, port_id, rate, valid_until)
        VALUES (?, ?, ?, ?)
        ON CONFLICT(agent_id, port_id) DO UPDATE
          SET rate        = excluded.rate,
              valid_until = excluded.valid_until,
              updated_at  = CURRENT_TIMESTAMP`
	_, err := db.Exec(stmt,
		q.AgentID,
		q.PortID,
		q.Rate,
		q.ValidUntil,
	)
	return err
}
