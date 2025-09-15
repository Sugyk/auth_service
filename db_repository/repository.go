package db_repository

import "database/sql"

type DBRepo struct {
	db *sql.DB
}

func NewDBRepo(db *sql.DB) *DBRepo {
	return &DBRepo{
		db: db,
	}
}
