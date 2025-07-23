package repository

import (
	"context"
	"database/sql"
	"portquote/internal/store"
)

type Port struct {
	ID      int64
	Name    string
	Country string
	City    string
}

type PortRepo struct {
	db *store.Store
}

func NewPortRepo(db *store.Store) *PortRepo {
	return &PortRepo{db: db}
}

func (r *PortRepo) GetAll(ctx context.Context) ([]Port, error) {
	const q = `
        SELECT id, name, country, city
          FROM ports`
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Port
	for rows.Next() {
		var p Port
		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Country,
			&p.City,
		); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *PortRepo) GetByID(ctx context.Context, id int64) (*Port, error) {
	const q = `
        SELECT id, name, country, city
          FROM ports
         WHERE id = ?`
	row := r.db.QueryRowContext(ctx, q, id)

	p := &Port{}
	if err := row.Scan(
		&p.ID,
		&p.Name,
		&p.Country,
		&p.City,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return p, nil
}
