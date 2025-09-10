package db_repository

import (
	"fmt"
)

func (d *DBRepo) CreateUser(login string, password string) error {
	if _, ok := d.db[login]; ok {
		return fmt.Errorf("user is already in db")
	}
	d.db[login] = password
	return nil
}
