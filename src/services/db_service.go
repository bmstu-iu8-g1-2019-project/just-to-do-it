package services

import (
	"database/sql"
	"fmt"
	"io"
	"os"

	"github.com/DATA-DOG/go-sqlmock"
)

type DB struct {
	*sql.DB
}

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

func NewMockDB() (*DB, error) {
	db, mock, err := sqlmock.New()

	if err != nil {
		return &DB{}, err
	}
	// Closes the database and prevents new queries from starting.
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "email", "login", "fullname", "password", "acc_verified"}).
		AddRow(1, "Just@mail.com", "To", "Do", "It", false).
		AddRow(2, "a", "b", "c", "d", true)

	mock.ExpectQuery("^SELECT (.+) FROM user_table*").
		WithArgs(2).
		WillReturnRows(rows)

	return &DB{db}, nil
}

func OpenConfigFile(filename string) (config string){
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	data := make([]byte, 64)
	for {
		n, err := file.Read(data)
		if err == io.EOF {
			break
		}
		config += string(data[:n])
	}
	return config
}
