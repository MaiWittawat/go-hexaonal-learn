package db

import (
	"database/sql"
	"log"
)

func InitializeMysqlDB(dbName string) (*sql.DB, error) {
	log.Println("Connect to mysqldb...")
	db, err := sql.Open(dbName, "root:1234@tcp(127.0.0.1:3306)/mydb")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func DisconnectMysqlDB(db *sql.DB) error {
	if err := db.Close(); err != nil {
		return err
	}
	return nil
}
