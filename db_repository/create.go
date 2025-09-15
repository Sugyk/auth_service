package db_repository

import (
	"database/sql"
	"fmt"
)

func (d *DBRepo) CreateUser(login string, password string) error {
	query := `
        INSERT INTO users (login, password_hash) 
        VALUES ($1, $2)
        ON CONFLICT (login) DO NOTHING
        RETURNING id
    `
	var id int

	err := d.db.QueryRow(query, login, password).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("login '%s' is already exists", login)
		}
		return fmt.Errorf("db error: %w", err)
	}
	return nil
}
