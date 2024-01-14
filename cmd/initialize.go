package main

import (
	"github.com/jmoiron/sqlx"
)

func newDB(dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", dsn)

	if err != nil {
		return nil, err
	}
	// force a connection and test that it worked
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
