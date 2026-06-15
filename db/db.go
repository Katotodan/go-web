package db

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var Database *sql.DB
var DabaseErr error

func ConnectDb() error {
	db, err := sql.Open("mysql", "root:my-secret-pw@tcp(127.0.0.1:3306)/book_store?parseTime=true")

	if err != nil {
		Database = nil
		return err
	}

	err = db.Ping()
	if err != nil {
		Database = nil
		return err
	}
	Database = db
	return nil
}
