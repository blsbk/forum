package repository

import (
	"database/sql"
	"io/ioutil"
)

func SetUpDB(dbType, dbName string) (*sql.DB, error) {
	db, err := sql.Open(dbType, dbName)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	if err = migrate(db); err != nil {
		return nil, err
	}
	return db, err
}

func migrate(db *sql.DB) error {
	fileByte, err := ioutil.ReadFile("migrations.sql")
	if err != nil {
		return err
	}
	_, err = db.Exec(string(fileByte))
	if err != nil {
		return err
	}
	return nil
}
