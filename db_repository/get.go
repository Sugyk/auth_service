package db_repository

func (d *DBRepo) GetUser(login string) (string, bool) {
	password, ok := d.db[login]
	return password, ok
}
