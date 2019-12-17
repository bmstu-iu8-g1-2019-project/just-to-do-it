package services

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"os"
	"regexp"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bmstu-iu8-g1-2019-project/just-to-do-it/src/models"
)

// structure for functions that access the database
type DB struct {
	*sql.DB
}

func ReadConfig() (config string) {
	DbDriver := os.Getenv("DB_driver")
	DbUsername := os.Getenv("DB_username")
	DbPassword := os.Getenv("DB_password")
	DbHost := os.Getenv("DB_host")
	DbPort := os.Getenv("DB_port")
	DbName := os.Getenv("DB_name")
	DbSslmode := os.Getenv("DB_sslmode")

	config = fmt.Sprintf("%s://%s:%s@%s:%s/%s?sslmode=%s",
		DbDriver, DbUsername, DbPassword, DbHost, DbPort, DbName, DbSslmode)
	return
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

func Setup(filename string, db *DB) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Setupfile opening error: ", err)
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		fmt.Println("Error after opening setupfile: ", err)
		return
	}

	bs := make([]byte, stat.Size())
	_, err = file.Read(bs)
	if err != nil {
		fmt.Println("Error after opening setupfile: ", err)
		panic(err)
	}

	command := string(bs)
	_, err = db.Exec(command)
	if err != nil {
		fmt.Println("Command error")
	}
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
