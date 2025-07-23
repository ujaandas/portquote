package repository

import (
	"database/sql"
	"portquote/internal/store"
)

type Port struct {
	ID      int64
	Name    string
	Country string
	City    string
}

func GetAllPorts(db *store.Store) ([]Port, error) {
	const q = `
        SELECT id, name, country, city
          FROM ports`
	rows, err := db.Query(q)
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

func GetPortByID(db *store.Store, id int64) (*Port, error) {
	const q = `
        SELECT id, name, country, city
          FROM ports
         WHERE id = ?`
	row := db.QueryRow(q, id)

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
