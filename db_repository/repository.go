package db_repository

type DBRepo struct {
	db map[string]string
}

func NewDBRepo() *DBRepo {
	return &DBRepo{
		db: make(map[string]string),
	}
}
