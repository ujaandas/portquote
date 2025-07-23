package repository

import (
	"context"
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

type QuotationRepo struct {
	db *store.Store
}

func NewQuotationRepo(db *store.Store) *QuotationRepo {
	return &QuotationRepo{db: db}
}

func (r *QuotationRepo) GetByAgent(ctx context.Context, agentID int64) ([]Quotation, error) {
	const q = `
        SELECT id, agent_id, port_id, rate, valid_until, updated_at
          FROM quotations
         WHERE agent_id = ?`
	rows, err := r.db.QueryContext(ctx, q, agentID)
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

func (r *QuotationRepo) GetById(ctx context.Context, id, agentID int64) (*Quotation, error) {
	const q = `
        SELECT id, agent_id, port_id, rate, valid_until, updated_at
          FROM quotations
         WHERE id = ? AND agent_id = ?`
	row := r.db.QueryRowContext(ctx, q, id, agentID)

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

func (r *QuotationRepo) UpdateById(ctx context.Context, id, agentID int64, rate float64, validUntil string) error {
	const stmt = `
    UPDATE quotations
       SET rate = ?, valid_until = ?, updated_at = CURRENT_TIMESTAMP
     WHERE id = ? AND agent_id = ?`
	_, err := r.db.ExecContext(ctx, stmt, rate, validUntil, id, agentID)
	return err
}

func (r *QuotationRepo) DeleteById(ctx context.Context, id, agentID int64) error {
	const stmt = `
        DELETE FROM quotations
         WHERE id = ? AND agent_id = ?`
	_, err := r.db.ExecContext(ctx, stmt, id, agentID)
	return err
}

func (r *QuotationRepo) GetByPort(ctx context.Context, portID int64) ([]Quotation, error) {
	const q = `
	SELECT id, agent_id, port_id, rate, valid_until, updated_at
		FROM quotations
	 WHERE port_id = ?
	 ORDER BY rate ASC`
	rows, err := r.db.QueryContext(ctx, q, portID)
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
