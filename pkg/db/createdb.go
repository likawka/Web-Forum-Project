package db

import (
	"database/sql"
	"io/ioutil"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(dbPath, schemaPath string) error {
	var err error
	DB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	sqlBytes, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		return err
	}

	_, err = DB.Exec(string(sqlBytes))
	if err != nil {
		return err
	}
	return nil
}