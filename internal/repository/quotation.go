package repository

import (
	"database/sql"
	"portquote/internal/store"
	"time"
)

type Quotation struct {
	ID         int64
	AgentID    int64
	PortID     int64
	Rate       float64
	ValidUntil time.Time
	UpdatedAt  time.Time
}

func GetQuotationsByAgent(db *store.Store, agentID int64) ([]Quotation, error) {
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

func GetQuotationByID(db *store.Store, id, agentID int64) (*Quotation, error) {
	const q = `
        SELECT id, agent_id, port_id, rate, valid_until, updated_at
          FROM quotations
         WHERE id = ? AND agent_id = ?`
	row := db.QueryRow(q, id, agentID)

	var qt Quotation
	if err := row.Scan(
		&qt.ID,
		&qt.AgentID,
		&qt.PortID,
		&qt.Rate,
		&qt.ValidUntil,
		&qt.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &qt, nil
}

func UpdateQuotation(db *store.Store, id, agentID int64, rate float64, validUntil string) error {
	const stmt = `
    UPDATE quotations
       SET rate = ?, valid_until = ?, updated_at = CURRENT_TIMESTAMP
     WHERE id = ? AND agent_id = ?`
	_, err := db.Exec(stmt, rate, validUntil, id, agentID)
	return err
}

func DeleteQuotation(db *store.Store, id, agentID int64) error {
	const stmt = `
        DELETE FROM quotations
         WHERE id = ? AND agent_id = ?`
	_, err := db.Exec(stmt, id, agentID)
	return err
}

func GetQuotationsByPort(db *store.Store, portID int64) ([]Quotation, error) {
	const q = `
	SELECT id, agent_id, port_id, rate, valid_until, updated_at
		FROM quotations
	 WHERE port_id = ?
	 ORDER BY rate ASC`
	rows, err := db.Query(q, portID)
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
