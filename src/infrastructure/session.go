package infrastructure

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func GetSession() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./database/database.sqlite")

	if err != nil {
		panic(err)
	}

	return db, nil
}
