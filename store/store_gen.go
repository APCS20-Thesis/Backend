package store

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
)

func (q *Store) Transaction(txFunc func(*sql.Tx) error) (err error) {
	tx, err := q.db.Begin()
	if err != nil {
		return
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			_ = tx.Rollback() // err is non-nil; don't change it
		} else {
			err = tx.Commit() // err is nil; if Commit returns error update err
		}
	}()
	err = txFunc(tx)
	return err
}

type Store struct {
	*Queries
	db *sqlx.DB
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{db: db, Queries: New(db)}
}

func (q *Store) Ping() error {
	return q.db.Ping()
}
