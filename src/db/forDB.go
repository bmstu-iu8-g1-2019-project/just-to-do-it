package db

import (
	"database/sql"
	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/tree/dev-s/src/auth"
)

func NewDB(dbSourceName string) (*DB, error) {
	db, err := sql.Open("postgres",dbSourceName)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &DB{db}, nil
}
