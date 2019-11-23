package services

import (
	"database/sql"
	"fmt"
)

// structure for functions that access the database
type DB struct {
	*sql.DB
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "pass"
	dbname   = "postgres"
)

func NewDB() (*DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	// open connection
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	// check connection
	if err = db.Ping(); err != nil {
		return nil, err
	}
	fmt.Println("Successfully connected!")
	return &DB{db}, nil
}
