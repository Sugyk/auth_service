package db_repository

import (
	"database/sql"
	"fmt"
)

func (d *DBRepo) GetUserPasswordHash(login string) (string, error) {
	query := `
		SELECT password_hash FROM Users
		WHERE login = $1
	`

	var password_hash string

	err := d.db.QueryRow(query, login).Scan(&password_hash)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", ErrUserNotFound
		}
		return "", fmt.Errorf("db error: %w", err)
	}
	return password_hash, nil
}
