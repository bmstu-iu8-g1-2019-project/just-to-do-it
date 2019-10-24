package services

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/models"
)

// structure for functions that access the database
type DB struct {
	*sql.DB
}

func NewDB(dbSourceName string) (*DB, error) {
	db, err := sql.Open("postgres", dbSourceName)
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

func OpenConfigFile(filename string) (config string) {
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

func NewMockGetDB(users []models.User) (*DB, error) {
	db, mock, err := sqlmock.New()

	if err != nil {
		return &DB{}, err
	}
	rows := sqlmock.NewRows([]string{"id", "email", "login", "fullname",
		"password", "acc_verified"})
	for i := range users {
		expectedUser := users[i]
		rows.AddRow(
			expectedUser.Id,
			expectedUser.Email,
			expectedUser.Login,
			expectedUser.Fullname,
			expectedUser.Password,
			expectedUser.AccVerified)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM user_table WHERE id = $1")).
			WithArgs(expectedUser.Id).WillReturnRows(rows)
	}
	return &DB{db}, nil
}

func NewMockLoginDB(users []models.User, hash []string) (*DB, error) {
	db, mock, err := sqlmock.New()

	if err != nil {
		return &DB{}, err
	}
	rows := sqlmock.NewRows([]string{"id", "email", "login", "fullname",
		"password", "acc_verified"})
	for i := range users {
		expectedUser := users[i]
		rows.AddRow(
			expectedUser.Id,
			expectedUser.Email,
			expectedUser.Login,
			expectedUser.Fullname,
			hash[i],
			expectedUser.AccVerified)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM user_table WHERE login = $1")).
			WithArgs(expectedUser.Login).WillReturnRows(rows)
	}
	return &DB{db}, nil
}

func NewMockDeleteDB(users []models.User, ids []int) (*DB, error) {
	db, mock, err := sqlmock.New()

	if err != nil {
		return &DB{}, err
	}
	rows := sqlmock.NewRows([]string{"id", "email", "login", "fullname",
		"password", "acc_verified"})
	for i := range users {
		expectedUser := users[i]
		rows.AddRow(
			expectedUser.Id,
			expectedUser.Email,
			expectedUser.Login,
			expectedUser.Fullname,
			expectedUser.Password,
			expectedUser.AccVerified)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM user_table WHERE id = $1")).
			WithArgs(expectedUser.Id).WillReturnRows(rows)
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM user_table WHERE id = $1")).
			WithArgs(ids[i]).
			WillReturnResult(sqlmock.NewResult(1, 1))
	}
	return &DB{db}, nil
}
