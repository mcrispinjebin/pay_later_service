package repository

import (
	"database/sql"
	"fmt"
)

var (
	DB  *sql.DB
	err error
)

func InitDB(server string) error {

	DB, err = sql.Open("mysql", server)

	if err != nil {
		return fmt.Errorf("%s %s", "mySql connection failed", err.Error())
	}

	return nil
}
