package db_repository

import "database/sql"

type DBRepo struct {
	db map[string]string
}

func NewDBRepo(db *sql.DB) *DBRepo {
	return &DBRepo{
		db: make(map[string]string),
	}
}
