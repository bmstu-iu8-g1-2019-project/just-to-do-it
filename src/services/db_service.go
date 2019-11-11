package services

import (
	"fmt"
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
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

//для тестов (пока ничего не написал(((((
func NewMockDB() (*DB, error) {
	db, mock, err := sqlmock.New()

	if err != nil {
		return &DB{}, err
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "email", "login", "fullname", "password", "acc_verified"}).
		AddRow(1, "Just@mail.com", "To", "Do", "It", false).
		AddRow(2, "a", "b", "c", "d", true)

	mock.ExpectQuery("^SELECT (.+) FROM user_table*").
		WithArgs(2).
		WillReturnRows(rows)

	return &DB{db}, nil
}
